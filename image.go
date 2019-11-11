package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

//type ImageOriginal struct {
//	Id int
//	//UserId int
//	Name string
//	Resolution string
//	//Path string
//}
//
//type ImageResized struct {
//	Id int
//	OriginalId int
//	Resolution string
//	//Path string
//}


func SaveImage(imageRaw io.Reader, fileName string) (loc string, err error) {

	// TODO: save info to db

	// save image to aws
	return saveFile(imageRaw, fileName)
}

func saveFile(file io.Reader, key string) (location string, err error) {

	uploader := s3manager.NewUploader(s3sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AS3_BUCKET_NAME),
		Key:    aws.String(key),
		Body:   file,
		ACL: aws.String("public-read"),
	})
	if err != nil {
		return
	}
	return result.Location, err
}