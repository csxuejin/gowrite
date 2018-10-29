package main

import (
	"log"
	"os"
	"sync"
	"time"
)

const (
	DEFAULT_WORKER_NUM  = 2
	DEFAULT_TIME_FORMAT = "2006-01-02 15:04:05"
)

type _Writer struct{}

var (
	Writer *_Writer
	wg     sync.WaitGroup
)

func (_ *_Writer) WriteFile(fileSize int64, unit string, fileAbsPath string) {
	log.Printf("Time before writing file %v: %v\n", fileAbsPath, time.Now().Format(DEFAULT_TIME_FORMAT))

	var cnt int64
	switch unit {
	case "M":
		cnt = fileSize

	case "G":
		cnt = 1024 * fileSize

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

	wg.Add(DEFAULT_WORKER_NUM)
	jobs := make(chan int64, 100)
	for w := 1; w <= DEFAULT_WORKER_NUM; w++ {
		go worker(w, jobs, fileAbsPath)
	}

	for i := int64(0); i < cnt; i++ {
		offset := i * MB
		jobs <- offset
	}
	close(jobs)

	wg.Wait()

	log.Printf("Time after writing file %v: %v\n\n", fileAbsPath, time.Now().Format(DEFAULT_TIME_FORMAT))
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
