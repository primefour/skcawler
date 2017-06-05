package spider

//this for downloader deamon
//spider ==> task queue ==> task ==> downloader ==> parser ==> task queue

func main_deamon() {

	for key, value := range SpiderMap {
		queue.AddTaskQueue(key, value.SpiderLevel)
	}

}
