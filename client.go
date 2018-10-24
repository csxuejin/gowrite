package main

import (
	"log"
	"os"
)

type GoWriteClient struct{}

func New() *GoWriteClient {
	return &GoWriteClient{}
}

func (_ *GoWriteClient) CreateFile() {

}

func (_ *GoWriteClient) CreateMultiFiles(input *CreateMultiFilesInput) {
	var err error
	if input.FilePathPre == "" {
		input.FilePathPre, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < input.FileNum; i++ {
		fileName := make([]byte, 16)
		RandStringBytesMaskImpr(fileName)
		log.Printf("filepath: %v, filename: %v\n", input.FilePathPre, string(fileName))
		WriteFile(input.FileSize, input.FileUnit, string(fileName), input.FilePathPre)
	}
}
