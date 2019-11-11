package main

import "github.com/gin-gonic/gin"

func RunHttpServer(addr string) error {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.GET("/images", NotImplemented)
	v1.GET("/images/:id", NotImplemented)
	v1.POST("/images/upload", NotImplemented)

	v1.GET("/images/:id/resized", NotImplemented)

	return r.Run(addr)
}

func NotImplemented(c *gin.Context) {
	c.String(200, "not implemented")
}