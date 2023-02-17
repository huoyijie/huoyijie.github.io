package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/huoyijie/huoyijie.github.io/dadada/server"
)

const (
	SRV_ADDR  string = "127.0.0.1:8000"
	IDX_URL          = "http://" + SRV_ADDR
	SLIDE_DIR string = "dadada/server/slides"
)

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func download(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func generate(slides []server.SlideModel) {
	var (
		body []byte
		err  error
	)

	fmt.Println("download index.html")
	for body, err = download(IDX_URL); err != nil; {
		printErr(err)
		time.Sleep(3 * time.Second)
	}

	printErr(os.Remove("index.html"))
	printErr(os.WriteFile("index.html", body, fs.ModePerm))

	printErr(os.RemoveAll("slides"))
	printErr(os.Mkdir("slides", fs.ModePerm))

	fmt.Println("download slides")
	for _, v := range slides {
		slideUrl, _ := url.JoinPath(IDX_URL, "slides", v.Name)

		fmt.Println("download slide", slideUrl)
		if body, err = download(slideUrl); err == nil {
			slidePath := filepath.Join("slides", v.Name+".html")
			fmt.Println("save slide", slidePath)
			os.WriteFile(slidePath, body, fs.ModePerm)
		} else {
			printErr(err)
		}
	}

	fmt.Println("done")
	os.Exit(0)
}

func main() {
	var g bool
	flag.BoolVar(&g, "g", false, "generate htmls")
	flag.Parse()

	slides := server.InitSlides(SLIDE_DIR)

	if g {
		go generate(slides)
	}

	server.Run(SLIDE_DIR, slides)
}
