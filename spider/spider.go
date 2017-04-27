package spider

import "fmt"

type PageTask struct {
	Request  *request.Request
	Response *http.Response
	//body content
	text []byte
	//dom tree
	dom   *goquery.Document
	items []data.DataCell // 存放以文本形式输出的结果数据
	files []data.FileCell // 存放欲直接输出的文件("Name": string; "Body": io.ReadCloser)
	err   error           // 错误标记
	sync.Mutex
}
