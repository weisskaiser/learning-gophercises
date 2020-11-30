package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"learning/url-shortener/urlshort"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func redirectConfig(file string) ([]byte, error){

	log.Printf("Reading config from %s\n", file)

	content, readErr := ioutil.ReadFile(file)

	return content, readErr
}

func parseFilePath(file, fileType string) string{

	dir := path.Dir(file)
	fileName := path.Base(file)

	if !strings.Contains(fileName, "."){
		fileName = fmt.Sprintf("%s.%s", fileName, fileType)
	}

	return path.Join(dir, fileName)

}

func main() {

	file := flag.String("file", "redirects", "File that should contain a list of redirects")
	configType := flag.String("type", "db", "Type is used to indicate if the format will be a JSON, YAML or db file")
	seedsFile := flag.String("seedsFile", "redirects.json", "JSON file that should contain a list of redirects to populate the db")

	flag.Parse()

	configPath := parseFilePath(*file, *configType)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var handler http.HandlerFunc
	var err error

	switch *configType {
		case "json", "yaml":
			configContent, readErr := redirectConfig(configPath)
			if readErr != nil {
				log.Fatal(readErr)
			}
			if(*configType == "json") {
				handler, err = urlshort.JSONHandler(configContent, mapHandler)
			}else{
				handler, err = urlshort.YAMLHandler(configContent, mapHandler)
			}
		case "db":
			seedsFilePointer, openErr := os.Open(*seedsFile)
			defer seedsFilePointer.Close()
			if openErr != nil {
				log.Fatal(openErr)
			}

			handler, err = urlshort.BoltDbHandler(configPath, mapHandler, seedsFilePointer)
		default:
			log.Fatal("Invalid config %s", *configType)
	}
	
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
