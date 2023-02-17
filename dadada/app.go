package main

import (
	"bufio"
	"embed"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xeonx/timeago"
)

//go:embed slides/*.md
var slideFS embed.FS

//go:embed templates/* templates/blocks/*
var tmplFS embed.FS

type Query struct {
	SlideName string `uri:"slideName" binding:"required"`
}

type SlideModel struct {
	Name, Title, Desc, Ago string
	Ctime                  time.Time
}

func initSlides() (slides []SlideModel) {
	entries, _ := os.ReadDir("slides")
	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		if slideName, found := strings.CutSuffix(v.Name(), ".md"); found {
			info, err := v.Info()
			if err != nil {
				continue
			}

			stat := info.Sys().(*syscall.Stat_t)
			ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))

			var desc string
			file, err := os.Open("slides/" + v.Name())
			if err != nil {
				continue
			}

			if scanner := bufio.NewScanner(file); scanner.Scan() {
				desc = strings.TrimPrefix(scanner.Text(), "[//]: # (")
				desc = strings.TrimSuffix(desc, ")")
			}

			slides = append(slides, SlideModel{
				Name:  slideName,
				Title: strings.ReplaceAll(slideName, "-", " "),
				Desc:  desc,
				Ctime: ctime,
				Ago:   timeago.Chinese.Format(ctime),
			})
		}
	}

	if len(slides) > 0 {
		sort.Slice(slides, func(i, j int) bool {
			return slides[i].Ctime.After(slides[j].Ctime)
		})
	}
	return
}

func main() {

	router := gin.Default()

	tmpl := template.Must(template.New("").ParseFS(tmplFS, "templates/*.htm", "templates/blocks/*.htm"))
	router.SetHTMLTemplate(tmpl)

	router.StaticFS("public", http.FS(slideFS))

	slides := initSlides()
	router.GET("/", func(c *gin.Context) {
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
