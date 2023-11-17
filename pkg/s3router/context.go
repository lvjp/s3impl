//go:generate go run golang.org/x/tools/cmd/stringer -type=RequestStyle -output context_string.go -type RequestStyle -trimprefix RequestStyle

package s3router

import (
	"net/http"
	"net/url"
	"strings"
)

type RequestStyle int

const (
	RequestStyleUnknown RequestStyle = iota
	RequestStyleVHost
	RequestStyleVPath
	RequestStyleCName
)

type Context struct {
	Hostname string
	BaseHost string
	Style    RequestStyle
	Bucket   string
	Key      string
}

type ContextFactory struct {
	Hosts []string
}

func (cf ContextFactory) Build(r *http.Request) *Context {
	url := url.URL{
		Host: r.Host,
	}

	q := &Context{
		Hostname: url.Hostname(),
	}

	if !cf.detectVPath(q, r) && !cf.detectVHost(q, r) {
		cf.fillCName(q, r)
	}

	return q
}

func (cf ContextFactory) detectVPath(q *Context, r *http.Request) bool {
	for _, host := range cf.Hosts {
		if q.Hostname != host {
			continue
		}

		q.Style = RequestStyleVPath
		q.BaseHost = host

		path := strings.TrimPrefix(r.URL.Path, "/")
		q.Bucket, q.Key, _ = strings.Cut(path, "/")

		return true
	}

	return false
}

func (cf ContextFactory) detectVHost(q *Context, r *http.Request) bool {
	for _, host := range cf.Hosts {
		bucket, cut := strings.CutSuffix(q.Hostname, "."+host)
		if !cut {
			continue
		}

		q.Style = RequestStyleVHost
		q.BaseHost = host
		q.Bucket = bucket
		q.Key = strings.TrimPrefix(r.URL.Path, "/")

		return true
	}

	return false
}

func (cf ContextFactory) fillCName(q *Context, r *http.Request) {
	q.Style = RequestStyleCName
	q.BaseHost = q.Hostname
	q.Bucket = q.Hostname
	q.Key = strings.TrimPrefix(r.URL.Path, "/")
}
