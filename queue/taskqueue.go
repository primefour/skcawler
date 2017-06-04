/*******************************************************************************
 **      Filename: skclawer/schd/taskqueue.go
 **        Author: crazyhorse
 **   Description: ---
 **        Create: 2017-05-31 14:55:37
 ** Last Modified: 2017-05-31 14:55:37
 ******************************************************************************/
package queue

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"sort"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type TaskQueue struct {
	SpiderName string //the own of TaskQueue
	priority   int
	RQueue     chan string
	RHistory   []int32 //sha256
	FQueue     map[int32]string
	hasher     hash.Hash
	lock       sync.Mutex
}

func NewTaskQueue(spiderName string, level int) *TaskQueue {
	var tmpQueue = new(TaskQueue{
		SpiderName: spiderName,
		priority:   level,
		RQueue:     make(chan string, 1024),
		RHistory:   make([]int32, 1024),
		FQueue:     make(map[int32]string, 1024),
		hash256:    sha256.New(),
	})
	return tmpQueue
}

func (self *TaskQueue) checkHistory(int32 value) bool {
	length := len(self.RHistory)
	var low, top int = 0, length - 1
	var found bool = false
	for low < top {
		m = (low + top) / 2
		if self.RHistory[m] == value {
			found = 1
			break
		} else if self.RHistory[m] < value {
			low = m + 1
		} else {
			top = m - 1
		}
	}
	return found
}

func (self *TaskQueue) linkHash(url string) int32 {
	self.hasher.Reset()
	self.hasher.Write([]byte(url))
	sum := self.hasher.Sum256()
	uptr := (*int32)(unsafe.Pointer(&sum[0]))
	return *uptr
}

func (self *TaskQueue) checkFail(value int32) bool {
	_, ok := self.FQueue[value]
	return ok
}

func (self *TaskQueue) AddRequest(url string) {
	value := linkHash(url)
	lock.Lock()
	defer lock.Unlock()
	if !(checkHistory(value) || checkFail(value)) {
		self.RQueue <- url
	}
}

func (self *TaskQueue) fetchRequest() string {
	lock.Lock()
	defer lock.Unlock()
	url <- self.RQueue
	return url
}

func (self *TaskQueue) FailRequest(url string) {
	value := linkHash(url)
	lock.Lock()
	defer lock.Unlock()
	self.FQueue[value] = url
}

func (self *TaskQueue) DoneRequest(url string) {
	lock.Lock()
	defer lock.Unlock()
	hashValue := linkHash(url)
	append(self.RHistory, hashValue)
	length = len(self.RHistroy)
	var i int = length - 1

	for i >= 1 {
		if self.RHistory[i] < self.RHistory[i-1] {
			self.RHistory[i], self.RHistory[i-1] = self.RHistory[i-1], self.RHistory[i]
		} else {
			break
		}
		i--
	}
}
