package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aws/aws-lambda-go/lambda"
)

// Event Request struct
type Event struct {
	Filename string `json:"filename,omitempty" form:"filename"`
	Body     string `json:"html,omitempty" form:"input"`
}

// Response Respose struct
type Response struct {
	Filename string `json:"filename,omitempty"`
}

var (
	pdfPath   = "/tmp"
	extRegexp = regexp.MustCompile(`(?i).pdf$`)
)

func handler(event Event) (*Response, error) {
	res := new(Response)

	if event.Filename == "" {
		return res, errors.New("no filename specified")
	}

	generator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	generator.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(event.Body)))

	if err := generator.Create(); err != nil {
		return res, err
	}

	filename := extRegexp.ReplaceAllString(event.Filename, "") + ".pdf"

	res.Filename = filename

	putS3PDF("wkhtmltopdf-output-aleksandr-korolev", filename, generator.Bytes())

	return res, nil
}

func init() {
	os.Setenv("WKHTMLTOPDF_PATH", os.Getenv("LAMBDA_TASK_ROOT"))
}

func main() {
	lambda.Start(handler)
}
