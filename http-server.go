package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func RunHttpServer(addr string) error {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.GET("/images", GetImagesHandler)
	v1.POST("/images/upload", UploadImageHandler)

	v1.GET("/images/:id/resized", GetResizedImagesHandler)

	return r.Run(addr)
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
	images, err := GetAllImages()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	var result = []gin.H{}

	for _, image := range images {
		result = append(result, gin.H{
			"id": image.Id,
			"name": image.FileName,
			"url": image.Url(),
		})
	}

	c.JSON(200, gin.H{"images": result})
}

func GetResizedImagesHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "incorrect id"})
		return
	}

	images, err := GetResizedImages(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	var result = []gin.H{}

	for _, image := range images {
		result = append(result, gin.H{
			"id": image.Id,
			"name": image.FileName,
			"width": image.Width,
			"height": image.Height,
			"url": image.Url(),
		})
	}

	c.JSON(200, gin.H{"images": result})
}
