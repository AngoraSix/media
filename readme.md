# Media Service

## Requisites

- Since we have the `go.mod` in place we can `export GO111MODULE=on`

- Init go mod (not always): `go mod init angorasix.com/media`

- Install go storage libraries: `go get cloud.google.com/go/storage`??

- Just once, generate code based on the design: `goa gen angorasix.com/media/design`

- Build the project with `go build -o media-svc cmd/main/main.go cmd/main/http.go` and run it `./media-svc`, possibly with env variables and arguments `GOOGLE_APPLICATION_CREDENTIALS=./config/gcpStorageCredentials.json ./media-svc -strategy google`

## Build and execute (with Docker)

- `docker build -t angorasix/hoc-media .`
- `docker run -d -p 7070:80 angorasix/hoc-media`

## Check the API with swagger

- Build the project with `go build -o openapi-svc cmd/openapi/main.go cmd/openapi/http.go` and run it `./opaenapi-svc`
- Browse `http://localhost/openapi.json`

- Run swagger UI `docker run -p 8080:8080 -e SWAGGER_JSON=localhost:80/openapi.json swaggerapi/swagger-ui`

## Upload image

- Request with `form-data`, Key => `file`, Key type => `File`, value => upload image

## Retrieve image

- Request e.g. `/static/uploads/2021122322421542_Screen Shot 2021-12-19 at 14.57.01.png`

