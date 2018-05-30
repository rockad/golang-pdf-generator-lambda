#!/usr/bin/env bash
mkdir -p ~/lambda

cd ~/lambda

cp /media/sf_go/src/golang-pdf-generator-lambda/main ~/lambda
cp /media/sf_go/src/golang-pdf-generator-lambda/wkhtmltopdf ~/lambda
chmod 755 main wkhtmltopdf
zip lambda-gen.zip main wkhtmltopdf
cp ~/lambda/lambda-gen.zip /media/sf_go/src/golang-pdf-generator-lambda