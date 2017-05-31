package io

import "time"
import "path"
import "path/filepath"

type PageData struct {
	SpiderName string //to fetch Spider
	Url        string
	ParentUrl  string
	TableData  []*DataElem
	FileList   []*FileElem
	ItemFields []string
	SubUrls    []string
	Time       time.Time //fetch time
}

const (
	CVS_TYPE = "cvs"
)

var DataWriterRouter = make(map[string](func(data PageData) error))

const (
	TABLE_COUNT = 10
	FILE_COUNT  = 10
)

var pageDataPool *sync.Pool = &sync.Pool{
	New: New,
}

func New() interface{} {
	ins = new(PageData)
	ins.Time = time.Now()
	ins.TableData = make([]*DataElem, TABLE_COUNT)
	ins.FileList = make([]*FileElem, FILE_COUNT)
	return ins
}

func GetPageData() *PageData {
	ins = pageDataPool.Get().(*PageData)
	return ins
}

func PutPageData(page *PageData) {
	for _, data := range page.TableData {
		PutDataElem(data)
	}
	page.TableData = nil

	for _, file := range page.FileList {
		PutFileElem(file)
	}
	page.FileList = nil
	page.ItemFields = nil
}

func (self *PageData) GetFileBaseName() string {
	spider, ok = GetSpiderByName(self.SpiderName)
	if ok {
		fileName := fmt.Sprintf("%s-%s-%v", pageData.Prefix, pageData.Time.Format(RFC3339))
	} else {
		log.E("get file base name failed ")
	}
}

func (self *PageData) GetSpider() (*Spider, bool) {
	spider, ok = GetSpiderByName(self.SpiderName)
	return spider, ok
}

func (self *PageData) GetRootPath() string {
	spider, ok = GetSpiderByName(self.SpiderName)
	if ok {
		path := spider.GetRootPath()
		if path == "" {
			//set current directory as root path
			path = util.GetPwd()
		}
	} else {
		log.E("get spider %s fail ", self.SpiderName)
		path = util.GetPwd()
	}

	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	return path
}

func (self *PageData) GetItemFields() []string {
	return ItemFields
}

func (self *PageData) setItemFields(fields []string) {
	self.ItemFields = fields
}

func (self *PageData) InsertDataItem(item *DataElem) {
	append(self.TableData, item)
}

func (self *PageData) InsertFileItem(item *FileElem) {
	append(self.FileList, item)
}
