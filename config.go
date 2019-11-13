package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var config struct {
	HttpHost string `yaml:"http_host"`
	HttpPort string `yaml:"http_port"`

	DbHost string `yaml:"db_host"`
	DbPort string `yaml:"db_port"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
	DbDatabase string `yaml:"db_database"`

	S3bucketName string `yaml:"s3_bucket_name"`
	S3region     string `yaml:"s3_region"`
}

func InitConfig(path string) error {

	fileConfig, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(fileConfig)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &config)
	return err
}
