package main

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed slides/*.md
var slideFS embed.FS

//go:embed templates/* templates/blocks/*
var tmplFS embed.FS

type Query struct {
	SlideName string `uri:"slideName" binding:"required"`
}

type SlideModel struct {
	Name, Title string
}

func main() {
	router := gin.Default()

	tmpl := template.Must(template.New("").ParseFS(tmplFS, "templates/*.htm", "templates/blocks/*.htm"))
	router.SetHTMLTemplate(tmpl)

	router.StaticFS("public", http.FS(slideFS))

	router.GET("/", func(c *gin.Context) {
		var slides []SlideModel

		entries, _ := os.ReadDir("slides")
		for _, v := range entries {
			if v.IsDir() {
				continue
			}
			if slideName, found := strings.CutSuffix(v.Name(), ".md"); found {
				slides = append(slides, SlideModel{
					Name:  slideName,
					Title: strings.ReplaceAll(slideName, "-", " "),
				})
			}
		}

		c.HTML(http.StatusOK, "index.htm", gin.H{"Slides": slides})
	})

	router.GET("slides/:slideName", func(c *gin.Context) {
		var query Query
		if err := c.BindUri(&query); err != nil {
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.HTML(http.StatusOK, "slide.htm", gin.H{
			"SlideName": query.SlideName,
		})
	})

	router.SetTrustedProxies(nil)
	router.Run("127.0.0.1:8000")
}
