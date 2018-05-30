export GOOS=linux
export GOARCH=amd64

dep ensure

GOOS=linux GOARCH=amd64 go build -v -o ./bin/pdf-gen ./pdf-generator
cp -f wkhtmltopdf bin/wkhtmltopdf

serverless deploy -v