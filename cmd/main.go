package main

import (
	"lambda-s3-ses-go/internal/s3"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler() {
	s3.GetFile()
}

func main() {
	lambda.Start(handler)
}
