package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/rusco/qunit"
	"github.com/soroushjp/humble/router"
	"strings"
	"time"
)

// browserSupportsPushState will be true if the current browser
// supports history.pushState and the onpopstate event.
var browserSupportsPushState = (js.Global.Get("onpopstate") != js.Undefined) &&
	(js.Global.Get("history") != js.Undefined) &&
	(js.Global.Get("history").Get("pushState") != js.Undefined)

type route struct {
	path   string
	params map[string]string
}

func main() {
	qunit.Test("Navigate", func(assert qunit.QUnitAssert) {
		qunit.Expect(3)
		routeChan := make(chan route)
		r := router.New()
		r.HandleFunc("/foo", func(params map[string]string) {
			routeChan <- route{
				path:   "/foo",
				params: params,
			}
		})
		r.Start()
		done := assert.Async()
		go r.Navigate("/foo")
		go func() {
			select {
			case gotRoute := <-routeChan:
				checkPath(assert, "/foo")
				assert.Equal(gotRoute.path, "/foo", "Triggered route had incorrect path.")
				assert.DeepEqual(gotRoute.params, map[string]string{}, "Triggered route had incorrect params.")
				done()
			case <-time.After(200 * time.Millisecond):
				// This is admittedly very akward. But AFIAK there is no equivalent of t.Fail or t.Error
				// in qunit.
				assert.Ok(false, "Route was not triggered within 200 milliseconds")
				assert.Ok(true, "")
				assert.Ok(true, "")
				done()
			}
		}()
	})

	qunit.Test("Route Params", func(assert qunit.QUnitAssert) {
		qunit.Expect(3)
		routeChan := make(chan route)
		r := router.New()
		r.HandleFunc("/foo/{param1}/{param2}", func(params map[string]string) {
			routeChan <- route{
				path:   "/foo/{param1}/{param2}",
				params: params,
			}
		})
		r.Start()
		done := assert.Async()
		go r.Navigate("/foo/bar/baz")
		expectedParams := map[string]string{
			"param1": "bar",
			"param2": "baz",
		}
		go func() {
			select {
			case gotRoute := <-routeChan:
				checkPath(assert, "/foo/bar/baz")
				assert.Equal(gotRoute.path, "/foo/{param1}/{param2}", "Triggered route had incorrect path.")
				assert.DeepEqual(gotRoute.params, expectedParams, "Triggered route had incorrect params.")
				done()
			case <-time.After(200 * time.Millisecond):
				// This is admittedly very akward. But AFIAK there is no equivalent of t.Fail or t.Error
				// in qunit.
				assert.Ok(false, "Route was not triggered within 200 milliseconds")
				assert.Ok(true, "")
				assert.Ok(true, "")
				done()
			}
		}()
	})

	qunit.Test("Back", func(assert qunit.QUnitAssert) {
		qunit.Expect(3)
		routeChan := make(chan route)
		r := router.New()
		r.HandleFunc("/foo", func(params map[string]string) {
			routeChan <- route{
				path:   "/foo",
				params: params,
			}
		})
		r.HandleFunc("/bar", func(params map[string]string) {
			routeChan <- route{
				path:   "/bar",
				params: params,
			}
		})
		r.Start()
		done := assert.Async()
		go func() {
			// Navigate to /foo
			r.Navigate("/foo")
			// Wait for the "/foo" handler to be triggered
			// once before continuing.
			<-routeChan
			// Navigate to /bar
			r.Navigate("/bar")
			// Wait for the "/bar" handler to be triggered
			// once before continuing.
			<-routeChan
			// Navigate back to /foo, which should trigger the onpopstate listener
			// or the onhashchange listener, depending on browser support.
			js.Global.Get("history").Call("back")
		}()
		go func() {
			select {
			case gotRoute := <-routeChan:
				checkPath(assert, "/foo")
				assert.Equal(gotRoute.path, "/foo", "Triggered route had incorrect path.")
				assert.DeepEqual(gotRoute.params, map[string]string{}, "Triggered route had incorrect params.")
				done()
			case <-time.After(200 * time.Millisecond):
				// This is admittedly very akward. But AFIAK there is no equivalent of t.Fail or t.Error
				// in qunit.
				assert.Ok(false, "Route was not triggered within 200 milliseconds")
				assert.Ok(true, "")
				assert.Ok(true, "")
				done()
			}
		}()
	})
}

func checkPath(assert qunit.QUnitAssert, expected string) {
	gotPath := ""
	if browserSupportsPushState {
		gotPath = js.Global.Get("location").Get("pathname").String()
	} else {
		hash := js.Global.Get("location").Get("hash").String()
		gotPath = strings.SplitN(hash, "#", 2)[1]
	}
	assert.Equal(gotPath, expected, "Path was not set correctly.")
}
