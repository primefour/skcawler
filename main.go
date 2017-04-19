package main

import "fmt"
import "io"
import "github.com/primefour/skclawer/downloader"

func main() {
	http_req := downloader.NewHttpRequest()
	http_req.Url = "http://www.baidu.com"
	http_req.Method = downloader.HTTP_GET
	task, err := downloader.NewHttpTask(*http_req)
	if err != nil {
		return
	}
	resp, err := task.DoTask()
	if err != nil {
		fmt.Printf("get a error %v ", err)
	}
	reader := resp.Body
	defer reader.Close()
	var buffer [10]byte

	for true {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("get a eof")
			}
			break
		}
		fmt.Printf("buffer %s ", buffer)
	}
}
