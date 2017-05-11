package io

import "fmt"
import "time"
import "sync"

type DataElem map[string]interface{}

type FileElem struct {
	FileName string
	Data     []byte
}

var dataElemPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return DataElem{}
	},
}

var fileElemPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		item := new(DataElem)
		item.Data = make([]byte, 0)
		return item
	},
}

func GetDataElem() *DataElem {
	elem := dataElemPool.Get().(*DataElem)
	return elem
}

func GetFileElem(name string, bytes []byte) *FileElem {
	elem := fileElemPool.Get().(*FileElem)
	elem.FileName = name
	elem.Data = bytes
	return elem
}

func PutDataElem(elem *DataElem) {
	for k, v := range *elem {
		delete(*elem, k)
	}
	dataElemPool.Put(elem)
}

func PutFileElem(elem *FileElem) {
	elem.FileName = ""
	elem.Data = nil
	fileElemPool.Put(elem)
}
