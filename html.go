package main

const homePage = `<!DOCTYPE html>
<html>
	<head>
		<title>Vault example application - Connect</title>
	</head>
	<body>
		<form action="/vault" method="post">
			<label for="address">Address</label>
			<input type="text" id="address" name="address" placeholder="{{ .Address }}">

			<label for="approle">AppRole</label>
			<input type="text" id="approle" name="approle" value="applications">

			<label for="role_id">RoleID</label>
			<input type="text" id="role_id" name="role_id">

			<label for="secret_id">SecretID</label>
			<input type="text" id="secret_id" name="secret_id">
			<input type="submit">
		</form>
	</body>
</html>
`

const vaultPage = `<!DOCTYPE html>
<html>
	<head>
		<title>Vault example application - Encrypt</title>
		<style>
			dt {
				display: inline-block;
			}
			dt::after {
				content: ": ";
			}
			dd {
				display: inline;
				margin: 0;
			}
			dd:after{
				display: block;
				content: '';
			}
		</style>
	</head>
	<body>
		<h1>Connected to {{ .Address }}</h1>
		<dl>
			<dt>Initialized</dt>
			<dd>{{ .Initialized }}</dd>

			<dt>Sealed</dt>
			<dd>{{ .Sealed }}</dd>

			<dt>Standby</dt>
			<dd>{{ .Standby }}</dd>

			<dt>ServerTimeUTC</dt>
			<dd>{{ .ServerTimeUTC }}</dd>

			<dt>Version</dt>
			<dd>{{ .Version }}</dd>

			<dt>Cluster</dt>
			<dd>{{ .ClusterName }} - {{ .ClusterID }}</dd>
		</dl>
		<form action="/encrypt" method="post">
			<label for="path">Secrets engine path</label>
			<input type="text" id="path" name="path" placeholder="tenant-keyrings">
			<label for="key">key</label>
			<input type="text" id="key" name="key">
			<textarea name="data">text to encrypt</textarea>
			<input type="submit">
		</form>
	</body>
</html>
`

const decryptPage = `<!DOCTYPE html>
<html>
	<head>
		<title>Vault example application - Decrypt</title>
	</head>
	<body>
		<form action="/decrypt" method="post">
			<label for="path">Secrets engine path</label>
			<input type="text" id="path" name="path" value="{{ .Path }}">
			<label for="key">key</label>
			<input type="text" id="key" name="key" value="{{ .Key }}">
			<label for="encrypted">Key</label>
			<textarea id="encrypted" name="encrypted">{{ .Encrypted }}</textarea>
			<input type="submit">
		</form>
	</body>
</html>
`
