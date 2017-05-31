/*******************************************************************************
 **      Filename: skclawer/schd/taskqueue.go
 **        Author: crazyhorse
 **   Description: ---
 **        Create: 2017-05-31 14:55:37
 ** Last Modified: 2017-05-31 14:55:37
 ******************************************************************************/
package schd

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
	RQueue     []string
	RHistory   []int32 //sha256
	hasher     hash.Hash
	lock       sync.Mutex
}

func NewTaskQueue(spiderName string, level int) *TaskQueue {
	var tmpQueue = new(TaskQueue{
		SpiderName: spiderName,
		priority:   level,
		RQueue:     make([]string, 1024),
		RHistory:   make([]int32, 1024),
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

func (self *TaskQueue) AddRequest(url string) {
	value := linkHash(url)
	if !checkHistory(value) {
		append(self.RQueue, url)
	}
}

func (self *TaskQueue) FetchRequest() string {

}

func (self *TaskQueue) DoneRequest(url string) {
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
