package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)


func ReadAll(r io.Reader) ([]byte , error){

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

	sourcePtr := flag.String("source" , "" , "path of file to be parsed")
	flag.Parse()

	filePath := *sourcePtr


	file , err := os.Open(filePath)
	fileType := filepath.Ext(filePath)
	fmt.Println(fileType)


	if err != nil {
		os.Exit(1)
	}

	if byteSLice , err := ReadAll(file); len(byteSLice) != 0 && err == io.EOF{
		fmt.Println(string(byteSLice))
	}

	defer file.Close()


}

