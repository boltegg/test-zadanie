package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jmoiron/sqlx"
)

var (
	s3sess *session.Session
	mysqlSess *sqlx.DB
)

func main() {

	err := InitConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	// init aws session
	s3sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.S3region)},
	))

	// init mysql session
	mysqlSess = sqlx.MustConnect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort,config.DbDatabase))

	// run http server
	err = RunHttpServer(fmt.Sprintf("%s:%s", config.HttpHost, config.HttpPort))
	if err != nil {
		panic(err)
	}
}