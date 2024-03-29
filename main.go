package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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

func main() {
	sourcePtr := flag.String("source", "", "path of file to be parsed")
	flag.Parse()

	filePath := *sourcePtr

	file, err := os.Open(filePath)

	fileType := filepath.Ext(filePath)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	var data Routes
	if byteSlice, err := ReadAll(file); err == io.EOF && len(byteSlice) != 0 {
		if fileType == ".yaml" {
			if err := yaml.Unmarshal(byteSlice, &data); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}
		if fileType == ".json" {
			if err := json.Unmarshal(byteSlice, &data); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

		}
	}
	fmt.Println(data)

}
