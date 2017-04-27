package downloader

import "fmt"
import "io"
import "bufio"
import "bytes"
import "os"

var pwd string

func init() {
	pwd, _ = os.Getwd()
	fmt.Printf("pwd = %s \n", pwd)
}

func Downloader(url string) ([]byte, error) {
	http_req := NewHttpRequest()
	http_req.Url = "http://www.sina.com.cn"
	http_req.Method = HTTP_GET
	task, err := NewHttpTask(*http_req)
	if err != nil {
		return nil, err
	}
	resp, err := task.DoTask()
	if err != nil {
		fmt.Printf("get a error %v ", err)
		return nil, err
	}

	var buff bytes.Buffer

	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	var buffer = make([]byte, 1024)
	for true {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("get a eof")
			} else {
				return nil, err
			}
			break
		}
		if n > 0 {
			buff.Write(buffer)
		}
	}
	fmt.Printf("buffer %s \n", buff)
	return buff.Bytes(), nil
}
