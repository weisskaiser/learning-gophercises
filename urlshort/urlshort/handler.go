package urlshort

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		url, success := pathsToUrls[path]

		if success {
			http.Redirect(w, r, url, http.StatusFound)
		}

		fallback.ServeHTTP(w, r)

	}
}

type routeRedirect struct {
	Path string
	URL  string
}

func toURIByPathMap(r []routeRedirect) map[string]string {

	routeMap := map[string]string{}

	for _, routeRedirect := range r {
		routeMap[routeRedirect.Path] = routeRedirect.URL
	}

	return routeMap
}

func parseYAML(yml []byte) ([]routeRedirect, error) {

	var routesRedirects []routeRedirect

	err := yaml.Unmarshal(yml, &routesRedirects)

	return routesRedirects, err
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlContent []byte, fallback http.Handler) (http.HandlerFunc, error) {

	routesRedirect, err := parseYAML(yamlContent)

	if err != nil {
		return nil, err
	}

	return MapHandler(toURIByPathMap(routesRedirect), fallback), nil
}

func parseJSON(jsonContent []byte) ([]routeRedirect, error) {

	var routesRedirects []routeRedirect

	err := json.Unmarshal(jsonContent, &routesRedirects)

	return routesRedirects, err
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//     [{
//      "path": "/some-path"
//		"url": "https://www.some-url.com/demo"
//		}]
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {

	routesRedirect, err := parseJSON(json)

	if err != nil {
		return nil, err
	}

	return MapHandler(toURIByPathMap(routesRedirect), fallback), nil
}

func defaultValues(r io.Reader) (map[string]string, error) {

	content, readErr := ioutil.ReadAll(r)

	if readErr != nil {
		return nil, readErr
	}

	routesRedirect, err := parseJSON(content)

	if err != nil {
		return nil, err
	}

	return toURIByPathMap(routesRedirect), nil
}

func createDb(bucket string, tx *bolt.Tx, defaultFileValues io.Reader) (*bolt.Bucket, error) {

	log.Printf("Bucket \"%s\" was not found, creating it with default values", bucket)

	b, createErr := tx.CreateBucket([]byte(bucket))

	if createErr != nil {
		return nil, createErr
	}

	values, generateValuesErr := defaultValues(defaultFileValues)
	if generateValuesErr != nil {
		return nil, generateValuesErr
	}

	for k, v := range values {
		insertErr := b.Put([]byte(k), []byte(v))
		if insertErr != nil {
			return nil, insertErr
		}
	}

	log.Println("db was created successfuly")

	return b, nil
}

// BoltDbHandler will read from the db file and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the db, then the
// fallback http.Handler will be called instead.
func BoltDbHandler(file string, fallback http.Handler, defaultFileValues io.Reader) (http.HandlerFunc, error) {

	bucket := "redirectRoutes"

	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, err
	}
	defer db.Close()

	redirectRoutesMap := map[string]string{}

	viewErr := db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		if b == nil {
			createdBucket, creationErr := createDb(bucket, tx, defaultFileValues)
			if creationErr != nil {
				return creationErr
			}
			b = createdBucket
		}

		return b.ForEach(func(k, v []byte) error {
			redirectRoutesMap[string(k)] = string(v)
			return nil
		})
	})

	if viewErr != nil {
		return nil, viewErr
	}

	return MapHandler(redirectRoutesMap, fallback), nil
}
