package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
		if pathVal, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, pathVal, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)

	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func BuildHandlerFromFile(filename string, fallback http.Handler) (http.Handler, error) {
	// Read file
	readfileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Determine file type by extension
	ext := filepath.Ext(filename)
	switch ext {
	case ".yaml", ".yml":
		return YAMLHandler(readfileData, fallback)
	case ".json":
		return JSONHandler(readfileData, fallback)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parseedJSON, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	// Convert parsed JSON to map[string]string format
	// to make similar to MapHandler input
	// so that we can reuse MapHandler
	pathsToUrls := make(map[string]string)
	for _, pu := range parseedJSON {
		pathsToUrls[pu.Path] = pu.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	parsedYAML, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	// Convert parsed YAML to map[string]string format
	// to make similar to MapHandler input
	// so that we can reuse MapHandler
	pathsToUrls := make(map[string]string)
	for _, pu := range parsedYAML {
		pathsToUrls[pu.Path] = pu.URL
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func parseYAML(data []byte) ([]pathURL, error) {
	var parsed []pathURL
	err := yaml.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func parseJSON(data []byte) ([]pathURL, error) {
	var parsed []pathURL
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
