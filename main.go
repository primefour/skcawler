package main

import "fmt"
import "github.com/primefour/skclawer/downloader"

func main() {
	_, err := downloader.Downloader("http://www.baidu.com")
	if err != nil {
		fmt.Printf("%v ", err)
	}
}
