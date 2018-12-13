package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"title": strings.Title,
	})

	r.LoadHTMLGlob("templates/*")
	r.GET("/:preset/:this", handleTerm)
	r.Run()
}

func handleTerm(c *gin.Context) {
	preset := c.Param("preset")
	this := c.Param("this")

	if _, exists := presets[preset]; exists {
		c.HTML(http.StatusOK, "term.html", gin.H{
			"preset": preset,
			"this":   this,
			// TODO: Parallelize
			"paradigmatic": queryParadigmatic(preset, this),
			"syntagmatic":  querySyntagmatic(preset, this),
			"containing":   queryContaining(preset, this),
		})
	}
}
