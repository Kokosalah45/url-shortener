package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func YAMLHandler(r io.Reader) ([]byte , error){

	buffer := make([]byte , 100)
	
	for {
		tempBuffer := make([]byte , 50)
		_ , err := r.Read(tempBuffer)
		if  err != nil {
			if err == io.EOF{
				return buffer , io.EOF
			}
			
			return  nil, err
		}
		buffer = append(buffer , tempBuffer...)

	}
	
}

func main (){

	filePathPtr := flag.String("path" , "" , "path of file to be parsed")
	flag.Parse()

	filePath := *filePathPtr

	yamlFile , err := os.Open(filePath)

	if err != nil {
		os.Exit(1)
	}
	if byteSLice , err := YAMLHandler(yamlFile); len(byteSLice) != 0 && err == io.EOF{
		fmt.Println(string(byteSLice))
	}

	defer yamlFile.Close()


}

