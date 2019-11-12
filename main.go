package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jmoiron/sqlx"
)

const (
	AS3_BUCKET_NAME = "mediaservice-test-zadanie-yalantis"
	AS3_REGION = "eu-central-1"
)

var (
	s3sess *session.Session
	mysqlSess *sqlx.DB
)

func main() {

	// init aws session
	s3sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(AS3_REGION)},
	))

	// init mysql session
	mysqlSess = sqlx.MustConnect("mysql", "root:@tcp(localhost:3306)/test_zadanie")

	// run http server
	err := RunHttpServer(":80")
	if err != nil {
		panic(err)
	}
}