package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ProgramArgs struct {
	source string
}

type RouteMap = map[string]string

type Route struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

type Routes = []Route

func ReadAll(r io.Reader) ([]byte, error) {
	buffer := make([]byte, 0)

	for {
		tempBuffer := make([]byte, 1024)
		n, err := r.Read(tempBuffer)
		if err != nil {
			if err == io.EOF {
				return buffer, io.EOF
			}
			return nil, err
		}
		buffer = append(buffer, tempBuffer[:n]...)
	}
}

func getArgs() *ProgramArgs {

	sourcePtr := flag.String("source", "", "path of file to be parsed")

	flag.Parse()

	source := *sourcePtr

	return &ProgramArgs{
		source: source,
	}

}

func setupRouteMap(routes Routes) RouteMap {
	routeMap := make(RouteMap)

	for _, route := range routes {
		routeMap[route.Path] = route.URL
	}

	return routeMap
}

func handlerFuncGenerator(routeMap RouteMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, exists := routeMap[path]; exists {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}else{
			data , _ := json.Marshal(struct{
				User_agent string `json:"user_agent"`
				Host string	`json:"host"`
				Method string `json:"method"`

			}{
				User_agent:  r.UserAgent(),
				Host: r.Host,
				Method: r.Method,

			})
			fmt.Fprintf(w,"hello %v" , html.UnescapeString(string(data)) )
		}
	}
}

func parseRoutes(file *os.File) Routes {
	fileType := filepath.Ext(file.Name())
	temp := Routes{}
	if byteSlice, err := ReadAll(file); err == io.EOF && len(byteSlice) != 0 {

		switch fileType {
			case ".yaml" :
				if err := yaml.Unmarshal(byteSlice, &temp); err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
			case ".json" :
				if err := json.Unmarshal(byteSlice, &temp); err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
		}
	}
	return temp
}

func main() {

	flags := getArgs()

	file, err := os.Open(flags.source)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	parsedRoutes := parseRoutes(file)
	routeMap := setupRouteMap(parsedRoutes)

	http.HandleFunc("/", handlerFuncGenerator(routeMap))
	log.Fatal(http.ListenAndServe(":8800", nil))
}
