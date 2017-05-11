package io

import (
	"encoding/csv"
	"fmt"
	"github.com/primefour/skclawer/log"
	"github.com/primefour/skclawer/utils"
	"os"
)

/*
type PageData struct {
	SpiderName string //as table name
	Url        string //for store file elem
	ParentUrl  string
	TableData  []DataElem //for schema
	FileList   []FileElem //will create a seperate file
	Time       time.Time //record the time of scraw
}
*/

//this is output data elem as csv file
func init() {
	DataWriterRouter[CVS_TYPE] = func(pageData *PageData) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			}
		}()

		dirName := utils.FileNameReplace(pageData.SpiderName)
		dirName = dirName + pageData.StartTime.Format(RFC3339)

		path := pageData.RootPath + dirName

		fileName := fmt.Sprintf("%s-%s-%v", pageData.Prefix, pageData.Time.Format(RFC3339), "_table.csv")

		filePath := path + fileName

		f, err := os.Stat(path)

		if err != nil || !f.IsDir() {
			if err := os.MkdirAll(path, 0777); err != nil {
				log.E("Error: %v\n", err)
			}
		}

		file, err := os.Create(filePath)

		if err != nil {
			log.E("%v", err)
		}

		file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
		writer := csv.NewWriter(file)

		th := self.MustGetRule(datacell["RuleName"].(string)).ItemFields

		if self.Spider.OutDefaultField() {
			th = append(th, "当前链接", "上级链接", "下载时间")
		}
		sheets[subNamespace].Write(th)

		defer func(file *os.File) {
			// 发送缓存数据流
			sheets[subNamespace].Flush()
			// 关闭文件
			file.Close()
		}(file)

		for _, datacell := range self.dataDocker {
			var subNamespace = util.FileNameReplace(self.subNamespace(datacell))
			if _, ok := sheets[subNamespace]; !ok {
				folder := config.TEXT_DIR + "/" + cache.StartTime.Format("2006-01-02 150405") + "/" + joinNamespaces(namespace, subNamespace)
				filename := fmt.Sprintf("%v/%v-%v.csv", folder, self.sum[0], self.sum[1])

				// 创建/打开目录
				f, err := os.Stat(folder)
				if err != nil || !f.IsDir() {
					if err := os.MkdirAll(folder, 0777); err != nil {
						logs.Log.Error("Error: %v\n", err)
					}
				}

				// 按数据分类创建文件
				file, err := os.Create(filename)

				if err != nil {
					logs.Log.Error("%v", err)
					continue
				}

				file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

				sheets[subNamespace] = csv.NewWriter(file)
				th := self.MustGetRule(datacell["RuleName"].(string)).ItemFields
				if self.Spider.OutDefaultField() {
					th = append(th, "当前链接", "上级链接", "下载时间")
				}
				sheets[subNamespace].Write(th)

				defer func(file *os.File) {
					// 发送缓存数据流
					sheets[subNamespace].Flush()
					// 关闭文件
					file.Close()
				}(file)
			}

			row := []string{}
			for _, title := range self.MustGetRule(datacell["RuleName"].(string)).ItemFields {
				vd := datacell["Data"].(map[string]interface{})
				if v, ok := vd[title].(string); ok || vd[title] == nil {
					row = append(row, v)
				} else {
					row = append(row, util.JsonString(vd[title]))
				}
			}
			if self.Spider.OutDefaultField() {
				row = append(row, datacell["Url"].(string))
				row = append(row, datacell["ParentUrl"].(string))
				row = append(row, datacell["DownloadTime"].(string))
			}
			sheets[subNamespace].Write(row)
		}
		return
	}
}
