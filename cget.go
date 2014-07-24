package cget

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type CachedResponse struct {
	Response *http.Response
	Body     []byte
}

type CachingGetter struct {
	sync.RWMutex
	respCache map[string]*CachedResponse
	cli       http.Client
}

type noopCloser struct {
	io.Reader
}

func (noopCloser) Close() error {
	return nil
}

func (c *CachedResponse) ToHTTP() *http.Response {
	return &http.Response{
		Status:           c.Response.Status,
		StatusCode:       c.Response.StatusCode,
		Proto:            c.Response.Proto,
		ProtoMajor:       c.Response.ProtoMajor,
		ProtoMinor:       c.Response.ProtoMinor,
		Header:           c.Response.Header,
		Body:             noopCloser{bytes.NewReader(c.Body)},
		ContentLength:    c.Response.ContentLength,
		TransferEncoding: c.Response.TransferEncoding,
		// And those are the fields we care about, I guess
	}
}

func (c *CachingGetter) get(url string, respChan chan<- *http.Response) (err error) {
	resp, err := c.cli.Get(url)
	if err != nil {
		close(respChan)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		close(respChan)
	} else {
		resp := &CachedResponse{Response: resp, Body: body}
		c.Lock()
		c.respCache[url] = resp
		c.Unlock()
		respChan <- resp.ToHTTP()
	}

	return err
}

// XXX: Except for Body, the response fields should not be modified. There's no
//      deep copy of Header e.g., so that's shared among the cached responses
func (c *CachingGetter) Get(url string, respChan chan<- *http.Response) {

	if c.respCache == nil {
		c.respCache = make(map[string]*CachedResponse)
	}

	c.RLock()
	val, existsInCache := c.respCache[url]
	c.RUnlock()

	if existsInCache {
		go func() {
			respChan <- val.ToHTTP()
		}()
	} else {
		go c.get(url, respChan)
	}
}
