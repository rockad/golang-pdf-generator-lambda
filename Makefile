build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/pdf pdf-generator/main.go