package queue

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
	ptask.matrix = make(*TaskQueue, 0)
	return ptask
}

func (self *PTasks) getTaskQueue() (*TaskQueue, bool) {
	self.sparse = 0
	length = len(self.matrix)

	if length == 0 {
		return nil, false
	}

	if self.index >= length {
		self.index = 0
	}

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
	append(self.matrix, task)
}

func (self *PTasks) removeTaskQueue(spiderName string) {
	length = len(self.matrix)
	for i, item := range self.matrix {
		if item.SpiderName == spiderName {
			tmp := self.matrix[0:i]
			tmp = append(tmp, self.matrix[i+1:])
			self.matrix = tmp
		}
	}
}

type scheduler struct {
	pqueue []*PTasks
	pindex int32
	sync.RWMutex
}

var ss = &scheduler{}

func init() {
	ss.pqueue = make(*PTasks, SCHEDULE_MAX_PRIORITY)

	i := 0
	for i < SCHEDULE_MAX_PRIORITY {
		ss.pqueue = newPTasks(32 - i)
		i++
	}
}

func AddTaskQueue(spiderName string, level int) *TaskQueue {
	tqueue := NewTaskQueue(spiderName, level)
	ss.RLock()
	defer ss.RUnlock()
	ss.pqueue[level].addTaskQueue(tqueue)
	return tqueue
}

func RemoveTaskQueue(spiderName string, level int) {
	ss.RLock()
	defer ss.RUnlock()
	ss.pqueue[level].removeTaskQueue(tqueue)
}

func fetchTaskQueue() *TaskQueue {
	ss.RLock()
	defer ss.RUnlock()
	i := pindex
	for i < SCHEDULE_MAX_PRIORITY {
		q, ok = ss.pqueue[i].getTaskQueue()
		if ok {
			i++
			pindex = i
			return q
		} else {
			i++
		}
	}

	if i > SCHEDULE_MAX_PRIORITY {
		i = 0
		for i < SCHEDULE_MAX_PRIORITY {
			ss.pqueue[i].resetClock()
		}
		i = 0
	}
	pindex = 0
	return nil
}

func FetchTask() string {
	queue := fetchTaskQueue()
	return queue.FetchQueue()
}
