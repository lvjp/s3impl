package s3router

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/lvjp/s3impl/pkg/utils"
)

func DetermineRoute(r *http.Request, acceptedHosts []string) (*Route, error) {
	d := derminator{
		Request:       r,
		AcceptedHosts: acceptedHosts,
	}
	d.hostname()

	if !d.detectVPath(r) && !d.detectVHost(r) {
		d.fillCName(r)
	}

	if err := d.action(); err != nil {
		return nil, err
	}

	return &d.Route, nil
}

type derminator struct {
	AcceptedHosts []string
	Request       *http.Request
	Route         Route
}

func (d *derminator) hostname() {
	url := url.URL{
		Host: d.Request.Host,
	}

	d.Route.Hostname = url.Hostname()
}

func (d *derminator) detectVPath(r *http.Request) bool {
	for _, host := range d.AcceptedHosts {
		if d.Route.Hostname != host {
			continue
		}

		d.Route.Style = RequestStyleVPath
		d.Route.BaseHost = host

		path := strings.TrimPrefix(r.URL.Path, "/")
		d.Route.Bucket, d.Route.Key, _ = strings.Cut(path, "/")

		return true
	}

	return false
}

func (d *derminator) detectVHost(r *http.Request) bool {
	for _, host := range d.AcceptedHosts {
		bucket, cut := strings.CutSuffix(d.Route.Hostname, "."+host)
		if !cut {
			continue
		}

		d.Route.Style = RequestStyleVHost
		d.Route.BaseHost = host
		d.Route.Bucket = bucket
		d.Route.Key = strings.TrimPrefix(r.URL.Path, "/")

		return true
	}

	return false
}

func (d *derminator) fillCName(r *http.Request) {
	d.Route.Style = RequestStyleCName
	d.Route.BaseHost = d.Route.Hostname
	d.Route.Bucket = d.Route.Hostname
	d.Route.Key = strings.TrimPrefix(r.URL.Path, "/")
}

func (d *derminator) routesTree() routesTree {
	haveBucket := d.Route.Bucket != ""
	haveKey := d.Route.Key != ""

	switch {
	case !haveBucket && !haveKey: // 00
		return routesForRoot
	case haveBucket && !haveKey: // 10
		return routesForbucket
	case haveBucket && haveKey: // 11
		return routesForObject
	default: // 01
		panic("cannot have a key wuthout a bucket")
	}
}

func (d *derminator) action() error {
	routesTree := d.routesTree()
	queries := d.Request.URL.Query()
	subresources := utils.KeysIntersection(routesTree, queries)

	var routeSelectorMap map[string]routeSelector
	switch len(subresources) {
	case 0:
		// Use the default route
		routeSelectorMap = routesTree[""]
	case 1:
		routeSelectorMap = routesTree[subresources[0]]
	default:
		return errors.New(
			"Conflicting query string parameters: " +
				strings.Join(subresources, ", "),
		)
	}

	selector, exists := routeSelectorMap[d.Request.Method]
	if !exists {
		return errors.New("Method not allowed: " + d.Request.Method)
	}

	d.Route.Action = selector(d.Route, queries, d.Request.Header)
	return nil
}
