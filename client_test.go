package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/golib/assert"
)

func Test_CreateMultiFiles(t *testing.T) {
	assertion := assert.New(t)
	currentPath, err := os.Getwd()
	assertion.Nil(err)

	client := New()

	input := &CreateMultiFilesInput{
		FilePathPre: path.Join(currentPath, "testfiles"),
		FileSize:    100,
		FileUnit:    "M",
		FileNum:     5,
	}

	client.CreateMultiFiles(input)
}

func Test_GetCurrentDirectoryPath(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
