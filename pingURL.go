package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	url  = flag.String("url", "", "URL to ping.")
	file = flag.String("file", "", "file to ping URLs in it.")
	dir  = flag.String("dir", "", "directory to ping URLs in it.")
)

func pingURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}

func handleFile(path string) {
}

func handleDir(path string) {
}

func main() {
	flag.Parse()
	fmt.Printf("url: \"%s\", file: \"%s\", dir: \"%s\"\n",
		*url, *file, *dir)
	if *url != "" && !pingURL(*url) {
		fmt.Printf("URL %s looks not alive.\n", *url)
		os.Exit(1)
	}
	if *file != "" {
		handleFile(*file)
	}
	if *dir != "" {
		handleDir(*dir)
	}
	os.Exit(0)
}
