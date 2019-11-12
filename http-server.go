package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
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

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "upload: " + err.Error()})
		return
	}
	width, errw := strconv.Atoi(c.PostForm("width"))
	height, errh := strconv.Atoi(c.PostForm("height"))

	if errw != nil || errh != nil {
		c.JSON(400, gin.H{"error": "incorrect input: width or height"})
		return
	}

	// open file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "open: " + err.Error()})
		return
	}

	img, err := NewImageOptions(fileHeader.Filename)
	if err != nil {
		c.JSON(400, gin.H{"error": "open: " + err.Error()})
		return
	}

	loc, err := img.Save(file)
	if err != nil {
		c.JSON(422, gin.H{"error": "save: " + err.Error()})
		return
	}

	imgResize, imgResizeRaw, err := img.Resize(file, width, height)
	if err != nil {
		c.JSON(422, gin.H{"error": "open: " + err.Error()})
		return
	}

	locResized, err := imgResize.Save(imgResizeRaw)
	if err != nil {
		c.JSON(422, gin.H{"error": "save: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"imageUrl": loc, "imageResizedUrl": locResized})

}

func GetImagesHandler(c *gin.Context) {
	// TODO: return slice of image links
}

func GetResizedImagesHandler(c *gin.Context) {
	//c.Query("id")
	// TODO: return slice of image links
}
