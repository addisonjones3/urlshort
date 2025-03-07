package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			url := pathsToUrls[r.URL.Path]
			if url != "" {
				http.Redirect(w, r, url, http.StatusPermanentRedirect)
			} else {
				fallback.ServeHTTP(w, r)
			}
		})
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := ParseYAML(yml)
	if err != nil {
		return nil, err
	}

	return MapHandler(BuildMap(pathsToUrls), fallback), err
}

type URLShortYAMLItem struct {
	Path string
	URL  string
}

func ParseYAML(yml []byte) (pathsToURLs []URLShortYAMLItem, err error) {
	err = yaml.Unmarshal(yml, &pathsToURLs)
	return
}

func BuildMap(URLShortYAMLItems []URLShortYAMLItem) map[string]string {
	URLPaths := make(map[string]string, len(URLShortYAMLItems))
	for _, item := range URLShortYAMLItems {
		URLPaths[item.Path] = item.URL
	}

	return URLPaths
}
