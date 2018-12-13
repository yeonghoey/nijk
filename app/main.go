package main

import (
	"fmt"
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
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"title": strings.Title,
		"sectionParam": func(t string, ts []string) sectionParam {
			return sectionParam{t, ts}
		},
	})

	router.LoadHTMLGlob("_templates/*")
	router.GET("/", handleIndex)
	router.GET("/:preset", handlePreset)
	router.GET("/:preset/:this", handleThis)
	router.Run()
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"presets": presets,
	})
}

func handlePreset(c *gin.Context) {
	preset := c.Param("preset")
	items, ok := presets[preset]
	if !ok {
		return
	}

	term := c.Query("term")
	term = strings.TrimSpace(term)
	if term == "" {
		c.HTML(http.StatusOK, "preset.html", gin.H{
			"preset": preset,
			"items":  items,
		})
	} else {
		location := fmt.Sprintf("/%s/%s", preset, term)
		c.Redirect(http.StatusMovedPermanently, location)
	}
}

func handleThis(c *gin.Context) {
	preset := c.Param("preset")
	this := c.Param("this")

	if _, exists := presets[preset]; !exists {
		return
	}

	c.HTML(http.StatusOK, "this.html", gin.H{
		"preset": preset,
		"this":   this,
		// TODO: Parallelize
		"paradigmatic": queryParadigmatic(preset, this),
		"syntagmatic":  querySyntagmatic(preset, this),
		"containing":   queryContaining(preset, this),
	})
}
