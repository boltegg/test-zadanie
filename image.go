package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"image"
	"io"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type ImageOptions struct {
	Id       int64  `db:"id"`
	UserId   int64  `db:"user_id"`
	FileName string `db:"file_name"`
	Format   string `db:"format"`
}

type ImageResizedOptions struct {
	Id         int64  `db:"id"`
	UserId     int64  `db:"user_id"`
	OriginalId int64  `db:"original_id"`
	FileName   string `db:"file_name"`
	Format     string `db:"format"`
	Width      int    `db:"width"`
	Height     int    `db:"height"`
}

func (i *ImageOptions) Path() string {
	return fmt.Sprintf("%d/%s_%d.%s", i.UserId, i.FileName, i.Id, i.Format)
}

func (i *ImageOptions) Url() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.S3bucketName, config.S3region, i.Path())
}

func (i *ImageResizedOptions) Path() string {
	return fmt.Sprintf("%d/%s_%d_%dx%d.%s", i.UserId, i.FileName, i.OriginalId, i.Width, i.Height, i.Format)
}

func (i *ImageResizedOptions) Url() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.S3bucketName, config.S3region, i.Path())
}

//

func NewImageOptions(fileNameFull string) (img ImageOptions, err error) {
	fname, format, err := parseFileName(fileNameFull)
	if err != nil {
		return
	}

	_, err = imaging.FormatFromExtension(format)
	if err != nil {
		err = fmt.Errorf("unsupported image format")
		return
	}

	img = ImageOptions{
		UserId:   1,
		FileName: fname,
		Format:   format,
	}

	return
}

func (i *ImageOptions) Save(imageRaw io.Reader) (location string, err error) {

	err = i.insert()
	if err != nil {
		return
	}

	location, err = uploadFileS3(imageRaw, i.Path())
	if err != nil {
		return
	}

	return
}

func (i *ImageOptions) insert() (err error) {

	res, err := mysqlSess.Exec("INSERT INTO image_original (user_id, file_name, format) VALUES(?, ?, ?)", i.UserId, i.FileName, i.Format)
	if err != nil {
		return
	}

	i.Id, err = res.LastInsertId()
	return
}

func (i *ImageOptions) Resize(imageRaw io.Reader, width int, height int) (imgResized ImageResizedOptions, dstRaw io.Reader, err error) {
	srcImage, _, err := image.Decode(imageRaw)
	if err != nil {
		return
	}

	dstImage := imaging.Resize(srcImage, width, height, imaging.Lanczos)
	format, _ := imaging.FormatFromExtension(i.Format)

	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)
	err = imaging.Encode(bufWriter, dstImage, format)
	if err != nil {
		return
	}
	dstRaw = bytes.NewReader(buf.Bytes())

	imgResized = ImageResizedOptions{
		UserId:     i.UserId,
		OriginalId: i.Id,
		FileName:   i.FileName,
		Format:     i.Format,
		Width:      width,
		Height:     height,
	}

	return
}

func (i *ImageOptions) Raw() (raw io.Reader, err error) {

	downloader := s3manager.NewDownloader(s3sess)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(config.S3bucketName),
			Key:    aws.String(i.Path()),
		})
	if err != nil {
		return
	}

	raw = bytes.NewReader(buf.Bytes())
	return
}

//

func (i *ImageResizedOptions) Save(imageRaw io.Reader) (location string, err error) {

	err = i.insert()
	if err != nil {
		return
	}

	location, err = uploadFileS3(imageRaw, i.Path())
	if err != nil {
		return
	}

	return
}

func (i *ImageResizedOptions) insert() (err error) {

	res, err := mysqlSess.Exec("INSERT INTO image_resized (original_id, width, height) VALUES(?, ?, ?)", i.OriginalId, i.Width, i.Height)
	if err != nil {
		return
	}

	i.Id, err = res.LastInsertId()
	return
}

//

func GetImage(id int) (image ImageOptions, err error) {

	err = mysqlSess.Get(&image, "SELECT * FROM image_original WHERE id=?", id)
	return
}

func GetAllImages() (images []ImageOptions, err error) {

	err = mysqlSess.Select(&images, "SELECT * FROM image_original")
	return
}

func GetResizedImages(id int) (images []ImageResizedOptions, err error) {

	err = mysqlSess.Select(&images, `SELECT image_resized.id, user_id, original_id, width, height, file_name,format FROM image_resized
LEFT JOIN image_original ON image_resized.original_id = image_original.id
WHERE image_resized.original_id=?`, id)
	return
}

//

func uploadFileS3(file io.Reader, key string) (location string, err error) {

	uploader := s3manager.NewUploader(s3sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.S3bucketName),
		Key:    aws.String(key),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return
	}
	return result.Location, err
}

func parseFileName(fileName string) (name, format string, err error) {
	slise := strings.Split(fileName, ".")
	if len(slise) < 2 {
		err = fmt.Errorf("format unknown")
	}
	name = strings.Join(slise[:len(slise)-1], ".")
	format = slise[len(slise)-1]
	return
}
