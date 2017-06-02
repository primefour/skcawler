package schd

import (
	"sync"
	"time"
)

const (
	SCHEDULE_MAX_PRIORITY = 32
)

type PTasks struct {
	matrix []*TaskQueue
	clock  int32
	level  int32
	index  int32
	sparse int32
}

func newPTasks(int level) *PTasks {
	ptask := new(PTasks)
	ptask.clock = level
	ptask.level = level
	ptask.index = 0
	return ptask
}

func (self *PTasks) getTaskQueue() (*TaskQueue, bool) {
	self.sparse = 0
	length = len(matrix)
	if self.clock > 0 {
		self.clock--
		for self.index < length {
			item, ok = self.matrix[self.index]
			if ok {
				return item, ok
			} else {
				self.sparse++
			}
			self.index++
		}
		if self.index >= length {
			self.index = 0
		}
	} else {
		return nil, false
	}
}

func (self *PTasks) resetClock() {
	self.clock = self.level
}

func (self *PTasks) addTaskQueue(task *TaskQueue) {
	append(matrix, task)
}

func (self *PTasks) removeTaskQueue(spiderName string) {
	var found bool = false
	length = len(self.matrix)
	for i, item := range self.matrix {
		if item.SpiderName == spiderName {
			found = true
			if i == length {
				self.matrix[i] = nil
			} else {
				self.matrix[i] = self.matrix[i+1]
			}
		}
		if found {
			if i == length {
				self.matrix[i] = nil
			} else {
				self.matrix[i] = self.matrix[i+1]
			}
		}
	}
}

type scheduler struct {
	pqueue []PTasks
	pindex int32
	sync.RWMutex
}

var ss = &scheduler{}

func init() {
	ss.pqueue = make(PTasks, SCHEDULE_MAX_PRIORITY)
	for i, item := range ss.pqueue {
		item = newPTasks(32 - i)
	}
}

func AddTaskQueue(spiderName string, level int) *TaskQueue {
	tqueue := NewTaskQueue(spiderName, level)
	ss.RLock()
	defer ss.RUnlock()
	ss.pqueue[level] = append(ss.pqueue[level], tqueue)
	return tqueue
}

func FetchTaskQueue() *TaskQueue {
	i := pindex
	for i < SCHEDULE_MAX_PRIORITY {
	}
}

func Stop() {
	ss.Lock()
	defer ss.Unlock()
	defer func() {
		recover()
	}()
	sdl.matrices = nil
}
