package main

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

const (
	MB                       int64 = 1024 * 1024
	DEFAULT_FILE_NUM               = 1
	DEFAULT_FILE_NAME_LENGTH       = 10
	DEFAULT_FILE_FOLDER            = "testfiles"
)

var (
	filePath string
	fileNum  int
	fileName []byte // generate random file name with fixed length
)

func init() {
	var err error
	filePath, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
		panic("Error happened when get the current directory path.")
	}

	filePath = path.Join(filePath, DEFAULT_FILE_FOLDER)
	if err := os.MkdirAll(filePath, 0777); err != nil {
		log.Fatal(err)
		panic("Error happened when create folder 'testfiles'")
	}

	fileNum = DEFAULT_FILE_NUM
	fileName = make([]byte, DEFAULT_FILE_NAME_LENGTH)
}

func main() {
	app := cli.NewApp()
	app.Author = "csxuejin@gmail.com"
	app.Usage = "Create specific size file"
	app.Version = "0.0.1"

	app.Action = func(c *cli.Context) error {
		var (
			err          error
			fileSizeArgs string
		)

		switch {
		case c.NArg() == 1:
			fileSizeArgs = strings.ToUpper(c.Args()[0])

		case c.NArg() == 2:
			fileSizeArgs = strings.ToUpper(c.Args()[0])
			fileNum, err = strconv.Atoi(c.Args()[1])
			if err != nil {
				log.Fatal(err)
				return nil
			}

		default:
			log.Fatal("Please input the file size, eg: 100M or 1G.")
			return nil
		}

		unit := string(fileSizeArgs[len(fileSizeArgs)-1])
		sizeStr := fileSizeArgs[:len(fileSizeArgs)-1]
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			log.Fatalf("The size of file is not correct: %v\n", sizeStr)
			return nil
		}

		switch unit {
		case "M", "G":
			log.Printf("Now create a file of %v size in %v.\n", fileSizeArgs, filePath)

		default:
			log.Fatal("Please input the file size, eg: 100M or 1G.")
			return nil
		}

		for i := 0; i < fileNum; i++ {
			Helper.RandStringBytesMaskImpr(fileName)
			Writer.WriteFile(size, unit, path.Join(filePath, string(fileName)))
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	return
}
