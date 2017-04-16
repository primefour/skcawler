/*******************************************************************************
 **      Filename: testGo/skclawer/downloader/request.go
 **        Author: crazyhorse
 **   Description: ---
 **        Create: 2017-04-14 14:48:38
 ** Last Modified: 2017-04-14 14:48:38
 ******************************************************************************/
package downloader

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"fmt"
	"github.com/primefour/skclawer/downloader/agent"
	"io"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	HTTP_HEAD  = "HEAD"
	HTTP_GET   = "GET"
	HTTP_POST  = "POST"
	HTTP_POSTM = "POST-M"
)

const (
	DefaultMethod      = HTTP_GET        // default method
	DefaultDialTimeout = 2 * time.Minute // default dial timeout
	DefaultRDWRTimeout = 2 * time.Minute // default connect timeout
	DefaultTryTimes    = 3               // the max try time
	DefaultRetryPause  = 2 * time.Second // retry pause interval
)

//this is input parameter
type HttpRequest struct {
	Url           string
	Method        string
	Header        http.Header
	EnableCookie  bool
	PostData      string
	DialTimeout   time.Duration
	RDWRTimeout   time.Duration
	TryTimes      int
	RetryPause    time.Duration
	RedirectTimes int
	Proxy         string
	StopClawer    bool
}

func (self *HttpRequest) Reset() {
	self.Method = DefaultMethod
	self.Header = make(http.Header)
	self.DialTimeout = DefaultDialTimeout
	self.RDWRTimeout = DefaultRDWRTimeout
	self.RetryPause = DefaultRetryPause
}

func NewHttpRequest() *HttpRequest {
	var obj = new(HttpRequest)
	obj.Reset()
	return obj
}

type RequestTask struct {
	method        string
	url           *url.URL
	proxy         *url.URL
	body          io.Reader
	header        http.Header
	enableCookie  bool
	dialTimeout   time.Duration
	rdwrTimeout   time.Duration
	tryTimes      int
	retryPause    time.Duration
	redirectTimes int
	client        *http.Client
}

func NewHttpTask(req HttpRequest) (task *RequestTask, err error) {
	task = new(RequestTask)

	task.url, err = UrlEncode(req.Url)

	if err != nil {
		return nil, err
	}

	if req.Proxy != "" {
		if task.proxy, err = url.Parse(req.Proxy); err != nil {
			return nil, err
		}
	}

	task.header = req.Header

	if task.header == nil {
		task.header = make(http.Header)
	}

	switch method := strings.ToUpper(req.Method); method {
	case HTTP_GET, HTTP_HEAD:
		task.method = method
	case HTTP_POST:
		task.method = method
		task.header.Add("Content-Type", "application/x-www-form-urlencoded")
		task.body = strings.NewReader(req.PostData)
	case HTTP_POSTM:
		task.method = "POST"
		body := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(body)
		values, _ := url.ParseQuery(req.PostData)
		for k, vs := range values {
			for _, v := range vs {
				writer.WriteField(k, v)
			}
		}
		err := writer.Close()
		if err != nil {
			return nil, err
		}
		task.header.Add("Content-Type", writer.FormDataContentType())
		task.body = body

	default:
		task.method = HTTP_GET
	}

	task.enableCookie = req.EnableCookie

	if len(task.header.Get("User-Agent")) == 0 {
		if task.enableCookie {
			fmt.Printf("%s", agent.UserAgents["common"][0])
			task.header.Add("User-Agent", agent.UserAgents["common"][0])
		} else {
			l := len(agent.UserAgents["common"])
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			task.header.Add("User-Agent", agent.UserAgents["common"][r.Intn(l)])
		}
	}

	task.dialTimeout = req.DialTimeout
	if task.dialTimeout < 0 {
		task.dialTimeout = 0
	}

	task.rdwrTimeout = req.RDWRTimeout
	task.tryTimes = req.TryTimes
	task.retryPause = req.RetryPause
	task.redirectTimes = req.RedirectTimes
	return
}

func (self *RequestTask) checkRedirect(req *http.Request, via []*http.Request) error {
	if self.redirectTimes == 0 {
		return nil
	}

	if len(via) >= self.redirectTimes {
		if self.redirectTimes < 0 {
			return fmt.Errorf("not allow redirects.")
		}
		return fmt.Errorf("stopped after %v redirects.", self.redirectTimes)
	}
	return nil
}

func (self *RequestTask) newCookieJar() *cookiejar.Jar {
	cookieJar, _ := cookiejar.New(nil)
	return cookieJar
}

// buildClient creates, configures, and returns a *http.Client type.
func (self *RequestTask) buildClient() error {
	client := &http.Client{
		CheckRedirect: self.checkRedirect,
	}

	if self.enableCookie {
		client.Jar = self.newCookieJar()
	}

	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, self.dialTimeout)
			if err != nil {
				return nil, err
			}
			if self.rdwrTimeout > 0 {
				c.SetDeadline(time.Now().Add(self.rdwrTimeout))
			}
			return c, nil
		},
	}

	if self.proxy != nil {
		transport.Proxy = http.ProxyURL(self.proxy)
	}

	if strings.ToLower(self.url.Scheme) == "https" {
		transport.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
		transport.DisableCompression = true
	}
	client.Transport = transport
	self.client = client
	return nil
}

// send uses the given *http.Request to make an HTTP request.
func (self *RequestTask) httpRequest() (resp *http.Response, err error) {
	req, err := http.NewRequest(self.method, self.url.String(), self.body)

	if err != nil {
		return nil, err
	}

	req.Header = self.header

	err = self.buildClient()

	if err != nil {
		return nil, err
	}

	if self.tryTimes <= 0 {
		for {
			resp, err = self.client.Do(req)
			if err != nil {
				if !self.enableCookie {
					l := len(agent.UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", agent.UserAgents["common"][r.Intn(l)])
				}
				time.Sleep(self.retryPause)
				continue
			}
			break
		}
	} else {
		for i := 0; i < self.tryTimes; i++ {
			resp, err = self.client.Do(req)
			if err != nil {
				if !self.enableCookie {
					l := len(agent.UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", agent.UserAgents["common"][r.Intn(l)])
				}
				time.Sleep(self.retryPause)
				continue
			}
			break
		}
	}

	return resp, err
}

func (self *RequestTask) DoTask() (resp *http.Response, err error) {
	resp, err = self.httpRequest()
	if err == nil {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			var gzipReader *gzip.Reader
			gzipReader, err = gzip.NewReader(resp.Body)
			if err == nil {
				resp.Body = gzipReader
			}

		case "deflate":
			resp.Body = flate.NewReader(resp.Body)

		case "zlib":
			var readCloser io.ReadCloser
			readCloser, err = zlib.NewReader(resp.Body)
			if err == nil {
				resp.Body = readCloser
			}
		}
	} else {
		fmt.Printf("%v ", err)
	}

	fmt.Printf("%s \n%s \n %s \n", resp.Request.Method, resp.Request.Header, resp.Request.Host)
	resp.Request.Method = self.method
	resp.Request.Header = self.header
	resp.Request.Host = self.url.Host
	return resp, nil
}
