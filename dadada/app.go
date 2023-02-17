package main

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed slides/*
var slideFS embed.FS

//go:embed templates/*
var tmplFS embed.FS

func main() {
	router := gin.Default()

	tmpl := template.Must(template.New("").ParseFS(tmplFS, "templates/*.htm"))
	router.SetHTMLTemplate(tmpl)

	router.StaticFS("public", http.FS(slideFS))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.htm", gin.H{})
	})

	router.SetTrustedProxies(nil)
	router.Run("127.0.0.1:8000")
}
