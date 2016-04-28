package main

import (
	"flag"
	"fmt"
)

var (
	url  = flag.String("url", "", "URL to ping.")
	file = flag.String("file", "", "file to ping URLs in it.")
	dir  = flag.String("dir", "", "directory to ping URLs in it.")
)

func main() {
	flag.Parse()
	fmt.Printf("url: \"%s\", file: \"%s\", dir: \"%s\"\n",
		*url, *file, *dir)
}
