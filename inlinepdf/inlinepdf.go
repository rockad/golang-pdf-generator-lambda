package main

import (
	"bytes"
	"context"
	"errors"
	"golang-pdf-generator-lambda/helpers"
	"regexp"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jung-kurt/gofpdf"
)

// Event Request struct
type Event struct {
	Filename string `json:"filename,omitempty" form:"filename"`
	Body     string `json:"html,omitempty" form:"input"`
}

// Response Respose struct
type Response struct {
	Filename string `json:"filename,omitempty"`
	Message  string `json:"message,omitempty"`
}

var (
	pdfPath   = "/tmp"
	extRegexp = regexp.MustCompile(`(?i).pdf$`)
)

func handler(ctx context.Context, event Event) (*Response, error) {
	ctx, seg := xray.BeginSegment(ctx, "generate-inline-pdf")

	res := new(Response)

	if event.Filename == "" {
		return res, errors.New("no filename specified")
	}

	ctx, subseg := xray.BeginSubsegment(ctx, "filename regex")
	filename := extRegexp.ReplaceAllString(event.Filename, "") + ".pdf"
	subseg.Close(nil)

	ctx, subseg = xray.BeginSubsegment(ctx, "create pdf")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 20)
	_, lineHt := pdf.GetFontSize()
	subseg.Close(nil)

	ctx, subseg = xray.BeginSubsegment(ctx, "write html")
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, event.Body)
	subseg.Close(nil)

	res.Filename = filename

	ctx, subseg = xray.BeginSubsegment(ctx, "write bytes")
	buf := new(bytes.Buffer)

	seg.AddMetadata("filename", filename)

	if err := pdf.Output(buf); err != nil {
		return res, err
	}
	subseg.Close(nil)

	ctx, subseg = xray.BeginSubsegment(ctx, "store in S3")
	helpers.PutS3PDF("wkhtmltopdf-output-aleksandr-korolev", filename, buf.Bytes())
	subseg.Close(nil)

	seg.Close(nil)

	return res, nil
}

func init() {
	xray.Configure(xray.Config{
		LogLevel:       "info", // default
		ServiceVersion: "1.2.3",
	})
}

func main() {
	lambda.Start(handler)
}
