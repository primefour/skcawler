package io

type PageData struct {
	StartTime  time.Time
	SpiderName string
	Prefix     string
	Url        string
	ParentUrl  string
	TableData  []DataElem
	FileList   []FileElem
	Time       time.Time
	RootPath   string
}

const (
	CVS_TYPE = "cvs"
)

var DataWriterRouter = make(map[string](func(data PageData) error))
