package main

import (
	"choose-your-adventure/adventure"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := flag.Int("port", 3000, "the port to start the web server")
	fileName := flag.String("file", "gopher.json", "JSON file containing your adventure story")
	flag.Parse()

	fmt.Printf("Using story in %s\n", *fileName)

	fp, err := os.Open(*fileName)

	if err != nil {
		panic(err)
	}

	story, err := adventure.JSONStory(fp)

	if err != nil {
		panic(err)
	}

	h := adventure.NewHandler(story)
	fmt.Printf("Starting server at :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
