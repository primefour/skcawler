package schd

import (
	"sync"
	"time"
)

type scheduler struct {
	status   int
	count    chan bool
	useProxy bool
	proxy    *proxy.Proxy
	matrices []*TaskQueue
	sync.RWMutex
}
