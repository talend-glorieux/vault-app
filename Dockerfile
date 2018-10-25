FROM golang AS build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo . 

FROM scratch
COPY --from=build /go/src/app/app /
CMD ["/app"]
