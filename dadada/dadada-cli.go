package main

import (
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

func main() {
	slides := server.InitSlides(SLIDE_DIR)

	go func() {
		var (
			body []byte
			err  error
		)

		for body, err = download(IDX_URL); err != nil; {
			time.Sleep(3 * time.Second)
		}

		os.WriteFile("index.html", body, fs.ModePerm)

		_ = os.Mkdir("slides", fs.ModePerm)

		for _, v := range slides {
			slideUrl, _ := url.JoinPath(IDX_URL, "slides", v.Name)
			fmt.Println(slideUrl)
			body, err = download(slideUrl)
			if err == nil {
				os.WriteFile(filepath.Join("slides", v.Name), body, fs.ModePerm)
			} else {
				fmt.Println(err)
			}
		}

		os.Exit(0)
	}()

	server.Run(SLIDE_DIR, slides)
}
