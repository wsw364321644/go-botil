package botil

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func PrintDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false
	fmt.Println("")
	for {
		select {
		case <-done:
			stop = true
		default:

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}
			var percent float64
			if(total<=0){
				percent=0
			}else{
				percent = float64(size) / float64(total) * 100
			}


			fmt.Printf("\r%d/%d(%.2f%%)\n",size,total,percent)
		}

		if stop {
			break
		}

		time.Sleep(time.Second*2)
	}
}

func DownloadFile(url string, dest string)error {

	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)

	start := time.Now()

	out, err := os.Create(path.String())

	if err != nil {
		fmt.Println(path.String())
		return err
	}

	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		return err
	}

	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		size=0
	}

	resp, err := http.Get(url)
	if err != nil {

		return err
	}
	defer resp.Body.Close()

	done := make(chan int64,1)
	go PrintDownloadPercent(done, path.String(), int64(size))


	n, err := io.Copy(out, resp.Body)
	done <- n
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
	return nil
}
