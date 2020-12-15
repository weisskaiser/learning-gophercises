package main

import (
	"flag"
	"html-link-parser/link"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	fileName := flag.String("file", "example.html", "Used to inform which file should be parsed")

	flag.Parse()

	content := readContent(*fileName)

	links, err := link.ParseContent(content)

	if err != nil {
		log.Fatal(err)
	}

	for _, l := range links {
		log.Printf("%+v\n", l)
	}
}

func readContent(fileName string) []byte {

	log.Printf("reading content from: %s\n", fileName)

	fp, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	content, readErr := ioutil.ReadAll(fp)

	if readErr != nil {
		log.Fatal(readErr)
	}
	return content
}
