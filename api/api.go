package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define your handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Handle all requests using net/http
	http.Handle("/", r)
}
