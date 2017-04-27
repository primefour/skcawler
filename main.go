package main

import "fmt"
import "strings"

//import "github.com/primefour/skclawer/downloader"

func BaseName(path string) string {
	x_idx := strings.LastIndex(path, "/")
	if x_idx < 0 {
		return ""
	}
	y_idx := strings.LastIndex(path, ".")

	if y_idx < 0 {
		return path[x_idx+1:]
	} else {
		return path[x_idx+1 : y_idx]
	}
}

func main() {
	/*_, err := downloader.Downloader("http://www.baidu.com")
	if err != nil {
		fmt.Printf("%v ", err)
	}
	*/

	var test_str = "Hello world"
	fmt.Printf("up string is %s \n", strings.ToUpper(test_str))
	fmt.Printf("low string is  %s \n", strings.ToLower(test_str))
	test_str = "hello/world.cap"
	fmt.Printf("base name of path is %s \n", BaseName(test_str))

	fmt.Printf("base name contain is %t \n", strings.Contains(test_str, "world"))

	test_str = "hewllo ssjfas dsdfalf a"

	xx := strings.Fields(test_str)

	fmt.Printf("%v \n", xx)
	for k, v := range xx {
		fmt.Printf(" k %d v %s \n", k, v)
	}
}
