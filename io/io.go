package io

type PageData struct {
	TableData []DataElem
	FileList  []FileElem
}

const (
	CVS_TYPE = "cvs"
)

var DataWriterRouter = make(map[string](func(data PageData) error))
