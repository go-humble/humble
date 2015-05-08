package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/rusco/qunit"
	"github.com/soroushjp/humble/router"
	"honnef.co/go/js/dom"
	"strings"
	"time"
)

// browserSupportsPushState will be true if the current browser
// supports history.pushState and the onpopstate event.
var (
	browserSupportsPushState = (js.Global.Get("onpopstate") != js.Undefined) &&
		(js.Global.Get("history") != js.Undefined) &&
		(js.Global.Get("history").Get("pushState") != js.Undefined)
	document = dom.GetWindow().Document()
)

// route is an internal representation of a route that has been triggered.
type route struct {
	path   string
	params map[string]string
}

func main() {
	qunit.Test("Navigate", func(assert qunit.QUnitAssert) {
		qunit.Expect(3)
		r := router.New()
		// For this test and other similar tests in this file, we're going to
		// create a channel called routeChan. Any router.Handlers that we set
		// up are just going to send a route object through the routeChan. That
		// way, we can recieve from the routeChan to detect when certain routes
		// were triggered. And by extension, we can check that the expected route
		// is triggered within a certain amount of time.
		routeChan := make(chan route)
		r.HandleFunc("/foo", newChanHandlerFunc("/foo", routeChan))
		r.Start()
		defer r.Stop()
		done := assert.Async()
		go r.Navigate("/foo")
		expectedRoute := route{
			path:   "/foo",
			params: map[string]string{},
		}
		go expectTriggeredRoute(assert, routeChan, "/foo", expectedRoute, done)
	})

	qunit.Test("Navigate with Params", func(assert qunit.QUnitAssert) {
		qunit.Expect(3)
		r := router.New()
		routeChan := make(chan route)
		r.HandleFunc("/foo/{param1}/{param2}", newChanHandlerFunc("/foo/{param1}/{param2}", routeChan))
		r.Start()
		defer r.Stop()
		done := assert.Async()
		go r.Navigate("/foo/bar/baz")
		expectedRoute := route{
			path: "/foo/{param1}/{param2}",
			params: map[string]string{
				"param1": "bar",
				"param2": "baz",
			},
		}
		go expectTriggeredRoute(assert, routeChan, "/foo/bar/baz", expectedRoute, done)
	})

	qunit.Test("Navigate Back", func(assert qunit.QUnitAssert) {
		qunit.Expect(3)
		r := router.New()
		routeChan := make(chan route)
		r.HandleFunc("/foo", newChanHandlerFunc("/foo", routeChan))
		r.HandleFunc("/bar", newChanHandlerFunc("/bar", routeChan))
		r.Start()
		defer r.Stop()
		done := assert.Async()
		go func() {
			// Navigate to "/foo"
			r.Navigate("/foo")
			// Wait for the "/foo" handler to be triggered
			// once before continuing.
			<-routeChan
			// Navigate to "/bar"
			r.Navigate("/bar")
			// Wait for the "/bar" handler to be triggered
			// once before continuing.
			<-routeChan
			// Navigate back to "/foo", which should trigger the onpopstate listener
			// or the onhashchange listener, depending on browser support, and in turn
			// trigger the corresponding router.Handler again.
			r.Back()
		}()
		expectedRoute := route{
			path:   "/foo",
			params: map[string]string{},
		}
		go expectTriggeredRoute(assert, routeChan, "/foo", expectedRoute, done)
	})

	qunit.Test("Intercept Links", func(assert qunit.QUnitAssert) {
		qunit.Expect(3)
		done := assert.Async()

		go func() {
			// Set up and start the router
			r := router.New()
			routeChan := make(chan route)
			r.HandleFunc("/biz", newChanHandlerFunc("/biz", routeChan))
			r.Start()
			defer r.Stop()

			// Create a link element, call InterceptLinks, and then simulate clicking it
			link := createLinkElement("/biz", "biz-link")
			r.InterceptLinks()
			link.Click()

			// Make sure that we navigated to the /biz route
			expectedRoute := route{
				path:   "/biz",
				params: map[string]string{},
			}
			expectTriggeredRoute(assert, routeChan, "/biz", expectedRoute, done)
		}()
	})
}

// newChanHandlerFunc will create and return a router.Handler which, when run,
// will send the route that was triggered (along with its params) through the
// given routeChan. This serves as an easy way for us to test which routes were
// triggered and what params were provided.
func newChanHandlerFunc(path string, routeChan chan route) router.Handler {
	return func(params map[string]string) {
		routeChan <- route{
			path:   path,
			params: params,
		}
	}
}

// checkPath checks that the browser is currently at the expected path. It
// knows whether or not to check the hash or the pathname based on browser
// support.
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

// expectTriggeredRoute will first wait to receive from routeChan. When it receives, it checks the route
// object to make sure it matches expectedRoute. It also checks window.location to make sure that the
// browser actually navigated to expectedPath. If the route is not triggered within 200 milliseconds, the
// check will fail. The function will call done (which should be from assert.Async()) when it is done.
func expectTriggeredRoute(assert qunit.QUnitAssert, routeChan chan route, expectedPath string, expectedRoute route, done func()) {
	select {
	case gotRoute := <-routeChan:
		checkPath(assert, expectedPath)
		assert.Equal(gotRoute.path, expectedRoute.path, "Triggered route had incorrect path.")
		assert.DeepEqual(gotRoute.params, expectedRoute.params, "Triggered route had incorrect params.")
		done()
	case <-time.After(200 * time.Millisecond):
		// This is admittedly very akward. But AFIAK there is no equivalent of t.Fail or t.Error
		// in qunit.
		assert.Ok(false, fmt.Sprintf("Route %s was not triggered within 200 milliseconds", expectedRoute.path))
		assert.Ok(true, "")
		assert.Ok(true, "")
		done()
	}
}

// createLinkElement will add an <a> element to the DOM with an href property equal to path
// and an id equal to id. It returns the created element.
func createLinkElement(path string, id string) *dom.HTMLAnchorElement {
	el := document.CreateElement("a")
	el.SetAttribute("href", path)
	el.SetID(id)
	document.QuerySelector("body").AppendChild(el)
	return el.(*dom.HTMLAnchorElement)
}
