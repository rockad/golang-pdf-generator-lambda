SET GOOS=linux
SET GOARCH=amd64

go build -o ./bin/pdfgenerator ./pdfgenerator
go build -o ./bin/inlinepdf ./inlinepdf