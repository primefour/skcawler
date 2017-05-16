package io

import (
	"encoding/csv"
	"fmt"
	"github.com/primefour/skclawer/log"
	"github.com/primefour/skclawer/utils"
	"os"
)

//this is output data elem as csv file
func init() {
	DataWriterRouter[CVS_TYPE] = func(pageData *PageData) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			}
		}()

		path = pageData.GetRootPath()
		f, err := os.Stat(path)
		if err != nil || !f.IsDir() {
			if err := os.MkdirAll(path, 0777); err != nil {
				log.E("Error: %v\n", err)
			}
		}

		fileName := fmt.Sprintf("%s-%s", pageData.GetFileBaseName(), "table.csv")
		filePath := path + fileName
		file, err := os.Create(filePath)
		if err != nil {
			log.E("%v", err)
		}

		spider := pageData.GetSpider()
		writer, ok := spider.GetCSVWriter(filePath)
		th := pageData.GetItemFields()

		if !ok {
			file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
			writer := csv.NewWriter(file)
			writer.Write(th)
			spider.SetCSVWriter(filePath, writer)
		}

		row := []string{}

		for _, dataElemPtr := range pageData.TableData {
			for _, title := range th {
				vd := *dataElemPtr
				if v, ok := vd[title].(string); ok || vd[title] == nil {
					row = append(row, v)
				} else {
					row = append(row, util.JsonString(vd[title]))
				}
			}
			writer.Write(row)
		}

		writer.Flush()
		file.Close()
	}
}
