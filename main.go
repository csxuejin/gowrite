package main

import (
	"fmt"
	"log"
	"math/rand"
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
	MB               FileUnit = 1024 * 1024
	GB               FileUnit = MB * 1024
	FILE_NAME_LENGTH          = 10
)

func main() {
	app := cli.NewApp()
	app.Author = "csxuejin@gmail.com"
	app.Usage = "Create specific size file"
	app.Version = "0.0.1"

	app.Action = func(c *cli.Context) error {
		fileNum := 1
		fileSizeArgs := ""
		filePath, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
			return nil
		}

		filePath = path.Join(filePath, "testfiles")
		if err := os.MkdirAll(filePath, 0777); err != nil {
			log.Fatal(err)
			return nil
		}

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
			fmt.Printf("Please input the file size, eg: 100M or 1G.\nYou also can specify the file path, eg: 100M /tmp. The default path is current directory.\n")
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
		case "M":
			fmt.Printf("Now create a file of %v size in %v.\n", fileSizeArgs, filePath)

		case "G":
			fmt.Printf("Now create a file of %v size in %v.\n", fileSizeArgs, filePath)

		default:
			fmt.Println("Please input the correct file size, eg: 100M or 1G")
			return nil
		}

		fileName := make([]byte, FILE_NAME_LENGTH)
		for i := 0; i < fileNum; i++ {
			RandStringBytesMaskImpr(fileName)
			WriteFile(size, unit, string(fileName), filePath)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	wg        sync.WaitGroup
	workerNum = 2
)

func WriteFile(fileSize int, unit string, fileName string, filePathPre string) {
	fmt.Printf("Time before writing file %v: %v\n", fileName, time.Now())

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

	filePath := path.Join(filePathPre, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Truncate(cnt * int64(MB))

	wg.Add(workerNum)
	jobs := make(chan int64, 100)
	for w := 1; w <= workerNum; w++ {
		go worker(w, jobs, filePath)
	}

	for i := int64(0); i < cnt; i++ {
		offset := i * int64(MB)
		jobs <- offset
	}
	close(jobs)

	wg.Wait()

	fmt.Printf("Time after writing file %v: %v\n\n", fileName, time.Now())
	return
}

func worker(id int, jobs <-chan int64, filePath string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
		wg.Done()
	}()

	buf := make([]byte, int(MB))
	for offset := range jobs {
		RandStringBytesMaskImpr(buf)
		if _, err := f.WriteAt(buf, offset); err != nil {
			log.Fatalln("err is : ", err)
			panic(err)
		}
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImpr(buff []byte) {
	n := len(buff)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			buff[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return
}
