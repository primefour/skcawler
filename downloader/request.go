/*******************************************************************************
 **      Filename: testGo/skclawer/downloader/request.go
 **        Author: crazyhorse
 **   Description: ---
 **        Create: 2017-04-14 14:48:38
 ** Last Modified: 2017-04-14 14:48:38
 ******************************************************************************/
package downloader

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	HTTP_GET   = "GET"
	HTTP_POST  = "POST"
	HTTP_POSTM = "POST-M"
)

const (
	DefaultMethod      = HTTP_GET        // default method
	DefaultDialTimeout = 2 * time.Minute // default dial timeout
	DefaultConnTimeout = 2 * time.Minute // default connect timeout
	DefaultTryTimes    = 3               // the max try time
	DefaultRetryPause  = 2 * time.Second // retry pause interval
)

type HttpRequest struct {
	Url           string
	Method        string
	Header        http.Header
	EnableCookie  bool
	PostData      string
	DialTimeout   time.Duration
	ConnTimeout   time.Duration
	TryTimes      int
	RetryPause    time.Duration
	RedirectTimes int
	Proxy         string
}

func (self *HttpRequest) Reset() {
	self.Method = DefaultMethod
	self.Header = make(http.Header)
	self.DialTimeout = DefaultDialTimeout
	self.ConnTimeout = DefaultConnTimeout
	self.RetryPause = DefaultRetryPause
}

func (self *DefaultRequest) check() {
	self.Method = strings.ToUpper(self.Method)
	if self.DialTimeout < 0 {
		self.DialTimeout = 0
	}

	if self.Method == "" {
		self.Method = DefaultMethod
	}

	if self.ConnTimeout < 0 {
		self.ConnTimeout = 0
	}
	if self.TryTimes < 0 {
		self.TryTimes = 1
	}
	if self.RetryPause <= 0 {
		self.RetryPause = DefaultRetryPause
	}
}

// url
func (self *DefaultRequest) GetUrl() string {
	self.check()
	return self.Url
}

// GET POST POST-M HEAD
func (self *DefaultRequest) GetMethod() string {
	self.check()
	return self.Method
}

// POST values
func (self *DefaultRequest) GetPostData() string {
	self.check()
	return self.PostData
}

// http header
func (self *DefaultRequest) GetHeader() http.Header {
	self.check()
	return self.Header
}

// enable http cookies
func (self *DefaultRequest) GetEnableCookie() bool {
	self.check()
	return self.EnableCookie
}

// dial tcp: i/o timeout
func (self *DefaultRequest) GetDialTimeout() time.Duration {
	self.check()
	return self.DialTimeout
}

// WSARecv tcp: i/o timeout
func (self *DefaultRequest) GetConnTimeout() time.Duration {
	self.check()
	return self.ConnTimeout
}

// the max times of download
func (self *DefaultRequest) GetTryTimes() int {
	self.check()
	return self.TryTimes
}

// the pause time of retry
func (self *DefaultRequest) GetRetryPause() time.Duration {
	self.check()
	return self.RetryPause
}

// the download ProxyHost
func (self *DefaultRequest) GetProxy() string {
	self.check()
	return self.Proxy
}

// max redirect times
func (self *DefaultRequest) GetRedirectTimes() int {
	self.check()
	return self.RedirectTimes
}
