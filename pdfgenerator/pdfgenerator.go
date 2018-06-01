package main

import (
	"context"
	"errors"
	"golang-pdf-generator-lambda/helpers"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
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

func handler(ctx context.Context, event Event) (*Response, error) {
	ctx, seg := xray.BeginSegment(ctx, "generate-wkhtmltopdf-pdf")

	res := new(Response)

	if event.Filename == "" {
		return res, errors.New("no filename specified")
	}

	ctx, subseg := xray.BeginSubsegment(ctx, "new pdf generator")
	generator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	subseg.Close(nil)

	ctx, subseg = xray.BeginSubsegment(ctx, "new reader")
	reader := strings.NewReader(event.Body)
	subseg.Close(nil)

	ctx, subseg = xray.BeginSubsegment(ctx, "new page")
	generator.AddPage(wkhtmltopdf.NewPageReader(reader))
	subseg.Close(nil)

	ctx, subseg = xray.BeginSubsegment(ctx, "create pdf")
	if err := generator.Create(); err != nil {
		return res, err
	}
	subseg.Close(nil)

	filename := extRegexp.ReplaceAllString(event.Filename, "") + ".pdf"

	res.Filename = filename

	helpers.PutS3PDF("wkhtmltopdf-output-aleksandr-korolev", filename, generator.Bytes())
	seg.Close(nil)
	return res, nil
}

func init() {
	os.Setenv("WKHTMLTOPDF_PATH", os.Getenv("LAMBDA_TASK_ROOT"))
	xray.Configure(xray.Config{
		LogLevel:       "info", // default
		ServiceVersion: "1.2.3",
	})
}

func main() {
	lambda.Start(handler)
}
