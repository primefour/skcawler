package spider

import "fmt"

type PageTask struct {
	Request  *request.Request
	Response *http.Response
	//body content
	text []byte
	//dom tree
	dom *goquery.Document
	sync.Mutex
}

type ParseInterface interface {
	PageParse(page *PageTask) *PageData
}

type PageParser struct {
	RegExpUrl string
	ParseInterface
}

type Spider struct {
	SaveRootPath string
	SpiderName   string
	SpiderStatus string
	SpiderLevel  int32
	SpiderTime   time.Time        //start time
	Request      *request.Request //the start request
	PageParses   []PageParser
}

var SpiderMap = make(map[string]*Spider, 20)

func GetSpiderByName(name string) (*Spider, bool) {
	spider, ok := SpiderMap[name]
	return spider, ok
}

func (self *Spider) GetRootPath() string {
	return SaveRootPath
}
