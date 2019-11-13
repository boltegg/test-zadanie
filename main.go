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

	err := initApp()
	if err != nil {
		panic(err)
	}

	// run http server
	err = RunHttpServer(fmt.Sprintf("%s:%s", config.HttpHost, config.HttpPort))
	if err != nil {
		panic(err)
	}
}

func initApp() error {
	err := InitConfig("config.yaml")
	if err != nil {
		return err
	}

	// init aws session
	s3sess, err = session.NewSession(&aws.Config{
		Region: aws.String(config.S3region)},
	)
	if err != nil {
		return err
	}

	// init mysql session
	mysqlSess, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort,config.DbDatabase))

	return err
}