package main

import (
	"github.com/gin-gonic/gin"
)

func RunHttpServer(addr string) error {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.GET("/images", NotImplemented)
	v1.POST("/images/upload", UploadImageHandler)

	v1.GET("/images/:id/resized", NotImplemented)

	return r.Run(addr)
}

func NotImplemented(c *gin.Context) {
	c.String(200, "not implemented")
}

func UploadImageHandler(c *gin.Context) {

	// TODO: check image (is correct file)
	// TODO: save info to db + save image to aws
	// TODO: resize
	// TODO: save resized info to db + save resized image to aws
	// TODO: return original image + resized image

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error":"Error while upload file, " + err.Error()})
	}
	//c.PostForm("width")
	//c.PostForm("height")

	// open file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, gin.H{"error":"Error while open file, " + err.Error()})
	}

	//img, _, err := image.Decode(file)

	loc, err := SaveImage(file, fileHeader.Filename)
	if err != nil {
		c.JSON(400, gin.H{"error":"Error while open file, " + err.Error()})
	}

	c.JSON(200, gin.H{"imageUrl":loc})

}

func GetImagesHandler(c *gin.Context) {
	// TODO: return slice of image links
}

func GetResizedImagesHandler(c *gin.Context) {
	//c.Query("id")
	// TODO: return slice of image links
}