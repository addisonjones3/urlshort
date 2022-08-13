package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/addisonjones3/gophercises/urlshort"
)

func main() {
	mux := defaultMux()
	YAMLFilePath := flag.String("yamlfilepath", "", "Pass full path to YAML config file")

	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yaml := `
        - path: /urlshort
          url: https://github.com/gophercises/urlshort
        - path: /urlshort-final
          url: https://github.com/gophercises/urlshort/tree/solution
`

	var YAMLHandlerInput []byte
	if *YAMLFilePath != "" {
		fmt.Printf("YAMLFilePath: %s\n", *YAMLFilePath)
		YAMLHandlerInput, _ = os.ReadFile(*YAMLFilePath)
	} else {
		YAMLHandlerInput = []byte(yaml)
	}

	yamlHandler, err := urlshort.YAMLHandler(YAMLHandlerInput, mapHandler)

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
