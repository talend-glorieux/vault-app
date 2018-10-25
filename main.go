package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := flag.String("port", "4242", "Application's port")
	dev := flag.Bool("dev", false, "Enables unsecure mode")
	flag.Parse()
	log.Infof("Running at http://localhost:%s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, NewApplication(*dev)))
}

// Application is a vault web application
type Application struct {
	templates   map[string]*template.Template
	vaultClient *api.Client
	router      *http.ServeMux
}

// NewApplication return a new Application
func NewApplication(dev bool) *Application {
	templates := make(map[string]*template.Template)
	templates["home"] = template.Must(template.New("home").Parse(homePage))
	templates["vault"] = template.Must(template.New("vault").Parse(vaultPage))
	templates["decrypt"] = template.Must(template.New("decrypt").Parse(decryptPage))

	config := api.DefaultConfig()
	if dev {
		config = &api.Config{
			Address:    "http://127.0.0.1:8200",
			HttpClient: http.DefaultClient,
		}
	}

	vaultClient, err := api.NewClient(config)
	if err != nil {
		log.Errorf("New vault client: %s", err)
		return nil
	}

	app := &Application{templates, vaultClient, http.NewServeMux()}
	app.router.HandleFunc("/", app.home)
	app.router.HandleFunc("/vault", app.vault)
	app.router.HandleFunc("/encrypt", app.encrypt)
	app.router.HandleFunc("/decrypt", app.decrypt)
	return app
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.router.ServeHTTP(w, r)
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Address string
		Token   string
	}{
		app.vaultClient.Address(),
		app.vaultClient.Token(),
	}
	err := app.templates["home"].Execute(w, data)
	if err != nil {
		log.Errorf("Execute template: %s", err)
	}
}

func (app *Application) vault(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Parsing form: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	address := r.Form.Get("address")
	if address != "" {
		app.vaultClient.SetAddress(r.Form.Get("address"))
	}

	token := r.Form.Get("token")
	if token != "" {
		app.vaultClient.SetToken(r.Form.Get("token"))
	}

	response, err := app.vaultClient.Sys().Health()
	if err != nil {
		log.Errorf("Vault health: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("RESPONSE: %+v\n", response)

	data := struct {
		Address       string
		Initialized   bool
		Sealed        bool
		Standby       bool
		ServerTimeUTC int64
		Version       string
		ClusterName   string
		ClusterID     string
	}{
		app.vaultClient.Address(),
		response.Initialized,
		response.Sealed,
		response.Standby,
		response.ServerTimeUTC,
		response.Version,
		response.ClusterName,
		response.ClusterID,
	}

	err = app.templates["vault"].Execute(w, data)
	if err != nil {
		log.Errorf("Execute template: %s", err)
	}
}

func (app *Application) encrypt(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Parsing form: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := path.Join(
		r.Form.Get("path"),
		"encrypt",
		r.Form.Get("key"),
	)

	data := make(map[string]interface{})
	data["plaintext"] = base64.StdEncoding.EncodeToString([]byte(r.Form.Get("data")))

	secret, err := app.vaultClient.Logical().Write(path, data)
	if err != nil {
		log.Errorf("Parsing form: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Warn(secret.Data)

	responseData := struct {
		Path      string
		Key       string
		Encrypted string
	}{
		r.Form.Get("path"),
		r.Form.Get("key"),
		secret.Data["ciphertext"].(string),
	}
	err = app.templates["decrypt"].Execute(w, responseData)
	if err != nil {
		log.Errorf("Execute template: %s", err)
	}
}

func (app *Application) decrypt(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Parsing form: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := path.Join(
		r.Form.Get("path"),
		"decrypt",
		r.Form.Get("key"),
	)

	data := make(map[string]interface{})
	data["ciphertext"] = r.Form.Get("encrypted")

	decryptedData, err := app.vaultClient.Logical().Write(path, data)
	if err != nil {
		log.Errorf("Parsing form: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decryptedText, err := base64.StdEncoding.DecodeString(decryptedData.Data["plaintext"].(string))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Fprint(w, string(decryptedText))
}
