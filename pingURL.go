package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	url  = flag.String("url", "", "URL to ping.")
	file = flag.String("file", "", "file or dir to ping URLs in it.")
)

func pingURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}

func asyncPingURL(url string, c chan string) {
	res := ""
	if !pingURL(url) {
		res = url
	}
	c <- res
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func urlsIn(text string) []string {
	var urls []string

	idx := 0
	for {
		idx = strings.Index(text, "http")
		if idx == -1 {
			break
		}
		fields_after := strings.Fields(text[idx:])
		if len(fields_after) > 0 {
			url := fields_after[0]
			if url[len(url)-1] == ')' {
				url = url[:len(url)-1]
			}
			urls = append(urls, url)
		}
		text = text[idx+4:]
	}

	return urls
}

func handleRegularFile(path string) {
	f, err := os.Open(path)
	handleError(err)

	s, err := f.Stat()
	handleError(err)
	f.Close()

	if s.IsDir() {
		log.Printf("-file argument is path to dir.\n")
		os.Exit(1)
	}

	dat, err := ioutil.ReadFile(path)
	urls := urlsIn(string(dat))
	nrAsyncPings := len(urls)
	fmt.Printf("found %d URLS in %s. Now ping...\n", len(urls), path)
	c := make(chan string)
	for _, url := range urls {
		go asyncPingURL(url, c)
	}

	for nrDonePings := 0; nrDonePings < nrAsyncPings; nrDonePings++ {
		pingRes := <-c
		if pingRes != "" {
			fmt.Printf("\t%s%s looks not alive.\n",
				"\x1B[33m", pingRes)
			fmt.Printf("%s", "\x1B[0m")
		}
	}
	close(c)
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		handleRegularFile(path)
	}

	return nil
}

func handleFile(path string) {
	f, err := os.Open(path)
	handleError(err)

	s, err := f.Stat()
	handleError(err)

	if !s.IsDir() {
		handleRegularFile(path)
	} else {
		err = filepath.Walk(path, visit)
		handleError(err)
	}
}

func main() {
	flag.Parse()
	if *url != "" && !pingURL(*url) {
		fmt.Printf("URL %s looks not alive.\n", *url)
		os.Exit(1)
	}
	if *file != "" {
		handleFile(*file)
	}
	os.Exit(0)
}
