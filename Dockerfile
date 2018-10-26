FROM golang:1.11 AS build
WORKDIR /
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo .

FROM scratch
COPY --from=build /vault-app /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/vault-app"]
