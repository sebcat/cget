package cget

import (
	"net/http"
	"testing"
)

func TestSingleGet(t *testing.T) {
	c := CachingGetter{}
	rchan := make(chan *http.Response)

	c.Get("http://www.google.com/", rchan)
	resp := <-rchan
	if resp == nil {
		t.Fail()
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		t.Fail()
	}
}

func TestCachedGet(t *testing.T) {
	c := CachingGetter{}
	rchan := make(chan *http.Response)

	c.Get("http://www.google.com/", rchan)
	resp := <-rchan
	if resp == nil {
		t.Fail()
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		t.Fail()
	}

	c.Get("http://www.google.com/", rchan)
	resp = <-rchan
	if resp == nil {
		t.Fail()
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		t.Fail()
	}
}

func BenchmarkGet(b *testing.B) {
	c := CachingGetter{}
	rchan := make(chan *http.Response)

	for i := 0; i < b.N; i++ {
		c.Get("http://www.google.com/", rchan)
		resp := <-rchan
		if resp == nil {
			b.Fail()
		}

		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 399 {
			b.Fail()
		}
	}
}
