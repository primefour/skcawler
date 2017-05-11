package io

import "time"

type PageData struct {
	SpiderName string //to fetch Spider
	Url        string
	ParentUrl  string
	TableData  []*DataElem
	FileList   []*FileElem
	ItemFields []string
	Time       time.Time
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
