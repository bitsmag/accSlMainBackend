In order to deploy this app to elastic beanstalk the application has to be build first

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/application application.go

then the entire content of this repository needs to be compressed to a .zip-file.
