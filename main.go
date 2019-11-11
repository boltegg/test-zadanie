package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	AS3_BUCKET_NAME = "mediaservice-test-zadanie-yalantis"
)

var (
	s3sess *session.Session
)

func main() {

	// init aws session
	s3sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	))


	err := RunHttpServer(":80")
	if err != nil {
		panic(err)
	}
}