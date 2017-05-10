package io

import "fmt"
import "time"
import "sync"

type DataElem struct {
	Name  string
	Value interface{}
}

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
		return FileElem{}
	},
}

func GetDataElem(name string, value interface{}) DataElem {
	elem := dataElemPool.Get().(DataElem)
	elem.Name = name
	elem.Value = value
	return elem
}

func GetFileElem(name string, bytes []byte) FileElem {
	elem := fileElemPool.Get().(FileElem)
	elem.FileName = name
	elem.Data = bytes
	return elem
}

func PutDataElem(elem DataElem) {
	elem.Name = ""
	elem.Value = nil
	dataElemPool.Put(elem)
}

func PutFileElem(elem FileElem) {
	elem.FileName = ""
	elem.Data = nil
	fileElemPool.Put(elem)
}
