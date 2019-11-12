package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"io"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type ImageOriginal struct {
	Id       int64
	UserId   int64
	FileName string
	Format   string
}

type ImageResized struct {
	Id         int64
	UserId     int64
	OriginalId int64
	FileName   string
	Format     string
	Width      int
	Height     int
}

func (i *ImageOriginal) Path() string {
	return fmt.Sprintf("%d/%s_%d.%s", i.UserId, i.FileName, i.Id, i.Format)
}

func (i *ImageOriginal) Url() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", AS3_BUCKET_NAME, AS3_REGION, i.Path())
}

func (i *ImageResized) Path() string {
	return fmt.Sprintf("%d/%s_%d_%dx%d.%s", i.UserId, i.FileName, i.Id, i.Width, i.Height, i.Format)
}

func (i *ImageResized) Url() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", AS3_BUCKET_NAME, AS3_REGION, i.Path())
}

//

func SaveImage(imageRaw io.Reader, fileFullName string) (loc string, err error) {

	fname, format, err := splitFileName(fileFullName)
	if err != nil {
		return
	}

	img := ImageOriginal{
		UserId:   1,
		FileName: fname,
		Format:   format,
	}

	err = img.Insert()
	if err != nil {
		logrus.Error("DB Insert error:", err)
		err = fmt.Errorf("Upload error")
		return
	}

	_, err = uploadFileS3(imageRaw, img.Path())
	if err != nil {
		logrus.Error("S3 upload file error:", err)
		err = fmt.Errorf("Upload error")
		return
	}

	return
}

func (i *ImageOriginal) Insert() (err error) {

	res, err := mysqlSess.Exec("INSERT INTO image_original (user_id, file_name, format) VALUES(?, ?, ?)", i.UserId, i.FileName, i.Format)
	if err != nil {
		return
	}

	i.Id, err = res.LastInsertId()
	return
}

//

func ResizeImage() {

}

//

func uploadFileS3(file io.Reader, key string) (location string, err error) {

	uploader := s3manager.NewUploader(s3sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AS3_BUCKET_NAME),
		Key:    aws.String(key),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return
	}
	return result.Location, err
}

func splitFileName(fileName string) (name, format string, err error) {
	slise := strings.Split(fileName, ".")
	if len(slise) < 2 {
		err = fmt.Errorf("File format unknown")
	}
	name = strings.Join(slise[:len(slise)-1], ".")
	format = slise[len(slise)-1]
	return
}