package main

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli"
)

type FileUnit int64

const (
	MB                       FileUnit = 1024 * 1024
	GB                       FileUnit = MB * 1024
	DEFAULT_FILE_NUM                  = 1
	DEFAULT_FILE_NAME_LENGTH          = 10
	DEFAULT_FILE_FOLDER               = "testfiles"
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
		size, err := strconv.Atoi(sizeStr)
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
			WriteFile(size, unit, path.Join(filePath, string(fileName)))
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	return
}

var (
	wg        sync.WaitGroup
	workerNum = 2
)

func WriteFile(fileSize int, unit string, fileAbsPath string) {
	log.Printf("Time before writing file %v: %v\n", fileAbsPath, time.Now())

	var cnt int64
	switch unit {
	case "M":
		cnt = int64(fileSize)

	case "G":
		cnt = int64(1024 * fileSize)

	default:
		log.Fatalf("The unit is wrong: %v\n", unit)
		return
	}

	f, err := os.Create(fileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Truncate(cnt * int64(MB))

	wg.Add(workerNum)
	jobs := make(chan int64, 100)
	for w := 1; w <= workerNum; w++ {
		go worker(w, jobs, fileAbsPath)
	}

	for i := int64(0); i < cnt; i++ {
		offset := i * int64(MB)
		jobs <- offset
	}
	close(jobs)

	wg.Wait()

	log.Printf("Time after writing file %v: %v\n\n", fileAbsPath, time.Now())
	return
}

func worker(id int, jobs <-chan int64, fileAbsPath string) {
	f, err := os.OpenFile(fileAbsPath, os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer func() {
		f.Close()
		wg.Done()
	}()

	buf := make([]byte, int(MB))
	for offset := range jobs {
		Helper.RandStringBytesMaskImpr(buf)
		if _, err := f.WriteAt(buf, offset); err != nil {
			log.Fatalln("f.WriteAt(): %v \n", err)
			panic(err)
		}
	}
}
