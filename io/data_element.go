package io

import "fmt"
import "time"
import "sync"

type DataValue struct {
	Name  string
	Value interface{}
}

type DataElem struct {
	SpiderName string
	Data       DataValue
	Url        string
	ParentUrl  string
	Time       time.Time
}

type FileElem struct {
	SpiderName string
	FileName   string
	Data       []byte
}

var dataElemPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return DataElem{}
	},
}

var fileElemPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return FileElem{}
	},
}

func GetDataElem(spiderName, url, parentUrl, string, data DateValue, time time.Time) DataElem {
	elem := dataElemPool.Get().(DataElem)
	elem.SpiderName = spiderName
	elem.Data = data
	elem.Url = url
	elem.ParentUrl = parentUrl
	elem.Time = time
	return elem
}

func GetFileElem(spiderName, name string, bytes []byte) FileElem {
	elem := fileElemPool.Get().(FileElem)
	elem.SpiderName = spiderName
	elem.FileName = name
	elem.Data = bytes
	return elem
}

func PutDataElem(elem DataElem) {
	elem.SpiderName = ""
	elem.Data = nil
	elem.Url = ""
	elem.ParentUrl = ""
	dataElemPool.Put(elem)
}

func PutFileElem(elem FileElem) {
	elem.SpiderName = ""
	elem.FileName = ""
	elem.Data = nil
	fileElemPool.Put(elem)
}
