package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/rusco/qunit"
	"github.com/soroushjp/humble/router"
	"time"
)

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
				assert.Equal(js.Global.Get("location").Get("pathname").String(), "/foo", "Path was not set correctly.")
				assert.Equal(gotRoute.path, "/foo", "Triggered route had incorrect path.")
				assert.DeepEqual(gotRoute.params, map[string]string{}, "Triggered route had incorrect params.")
				done()
			case <-time.After(200 * time.Millisecond):
				// This is admiteddly very akward. But AFIAK there is no equivalent of t.Fail or t.Error
				// in qunit.
				assert.Ok(false, "Route was not triggered within 200 milliseconds")
				assert.Ok(true, "")
				assert.Ok(true, "")
				done()
			}
		}()
	})
}
