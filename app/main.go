package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type sectionParam struct {
	Title string
	Thats []string
}

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"title": strings.Title,
		"sectionParam": func(t string, ts []string) sectionParam {
			return sectionParam{t, ts}
		},
	})

	r.LoadHTMLGlob("_templates/*")
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
