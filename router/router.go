package router

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
	"honnef.co/go/js/dom"
	"regexp"
	"strings"
)

var (
	// browserSupportsPushState will be true if the current browser
	// supports history.pushState and the onpopstate event.
	browserSupportsPushState = (js.Global.Get("onpopstate") != js.Undefined) &&
		(js.Global.Get("history") != js.Undefined) &&
		(js.Global.Get("history").Get("pushState") != js.Undefined)
	document dom.HTMLDocument
)

func init() {
	if js.Global != nil {
		var ok bool
		document, ok = dom.GetWindow().Document().(dom.HTMLDocument)
		if !ok {
			panic("Could not convert document to dom.HTMLDocument")
		}
	}
}

// Router is responsible for handling routes. If history.pushState is
// supported, it uses that to navigate from page to page and will listen
// to the "onpopstate" event. Otherwise, it sets the hash component of the
// url and listens to changes via the "onhashchange" event.
type Router struct {
	routes []*route
}

// Handler is a function which is run in response to a specific
// route. A Handler takes the url parameters as an argument.
type Handler func(params map[string]string)

// New creates and returns a new router
func New() *Router {
	return &Router{
		routes: []*route{},
	}
}

type route struct {
	regex      *regexp.Regexp // Regex pattern that matches route
	paramNames []string       // Ordered list of query parameters expected by route handler
	handler    Handler        // Handler called when route is matched
}

// HandleFunc will cause the router to call f whenever window.location.pathname
// (or window.location.hash, if history.pushState is not supported) matches path.
// path can contain any number of parameters which are denoted with curly brackets.
// So, for example, a path argument of "users/{id}" will be triggered when the user
// visits users/123 and will call the handler function with params["id"] = "123".
func (r *Router) HandleFunc(path string, handler Handler) {
	r.routes = append(r.routes, newRoute(path, handler))
}

// newRoute returns a route with the given arguments. paramNames and regex
// are calculated from the path
func newRoute(path string, handler Handler) *route {
	route := &route{
		handler: handler,
	}
	strs := strings.Split(path, "/")
	strs = removeEmptyStrings(strs)
	pattern := `^`
	for _, str := range strs {
		if str[0] == '{' && str[len(str)-1] == '}' {
			pattern += `/`
			pattern += `([\w+-]*)`
			route.paramNames = append(route.paramNames, str[1:(len(str)-1)])
		} else {
			pattern += `/`
			pattern += str
		}
	}
	pattern += `/?$`
	route.regex = regexp.MustCompile(pattern)
	return route
}

// Start causes the router to listen for changes to window.location and
// trigger the appropriate handler whenever there is a change.
func (r *Router) Start() {
	if browserSupportsPushState {
		r.watchHistory()
	} else {
		r.setInitialHash()
		r.watchHash()
	}
}

// Stop causes the router to stop listening for changes, and therefore
// the router will not trigger any more router.Handler functions.
func (r *Router) Stop() {
	if browserSupportsPushState {
		js.Global.Set("onpopstate", nil)
	} else {
		js.Global.Set("onhashchange", nil)
	}
}

// Navigate will trigger the handler associated with the given path
// and update window.location accordingly. If the browser supports
// history.pushState, that will be used. Otherwise, Navigate will
// set the hash component of window.location to the given path.
func (r *Router) Navigate(path string) {
	if browserSupportsPushState {
		pushState(path)
		r.pathChanged(path)
	} else {
		setHash(path)
	}
}

// Back will cause the browser to go back to the previous page.
// It has the same effect as the user pressing the back button,
// and is just a wrapper around history.back()
func (r *Router) Back() {
	js.Global.Get("history").Call("back")
}

// InterceptLinks intercepts click events on links of the form <a href="/foo"></a>
// and calls router.Navigate("/foo") instead, which triggers the appropriate Handler
// instead of requesting a new page from the server.
func (r *Router) InterceptLinks() {
	for _, link := range document.Links() {
		link.AddEventListener("click", true, func(event dom.Event) {
			href := link.GetAttribute("href")
			switch {
			case href == "":
				return
			case strings.HasPrefix(href, "http://"), strings.HasPrefix(href, "https://"), strings.HasPrefix(href, "//"):
				// These are external links and should behave normally.
				return
			case strings.HasPrefix(href, "#"):
				// These are anchor links and should behave normally.
				// Recall that even when we are using the hash trick, href
				// attributes should be relative paths without the "#" and
				// router will handle them appropriately.
				return
			case strings.HasPrefix(href, "/"):
				// These are relative links. The kind that we want to intercept.
				event.PreventDefault()
				go r.Navigate(href)
			}
		})
	}
}

// setInitialHash will set hash to / if there is currently no hash.
// Then it will trigger the appropriate
func (r *Router) setInitialHash() {
	if getHash() == "" {
		setHash("/")
	} else {
		r.pathChanged(getPathFromHash(getHash()))
	}
}

// pathChanged should be called whenever the path changes and will trigger
// the appropriate handler
func (r *Router) pathChanged(path string) {
	// path is everything after the '#'
	strs := strings.Split(path, "/")
	strs = removeEmptyStrings(strs)
	// Compare given path against regex patterns of routes. Preference given to routes with most literal (non-query) matches.
	// Route 1: /todos/work
	// Route 2: /todos/{category}
	// Path /todos/work will match Route #1
	leastParams := -1
	var bestRoute *route
	var bestMatches []string
	for _, route := range r.routes {
		matches := route.regex.FindStringSubmatch(path)
		if matches != nil {
			if (leastParams == -1) || (len(matches) < leastParams) {
				leastParams = len(matches)
				bestRoute = route
				bestMatches = matches[1:]
			}
		}
	}
	// If no routes match, we throw console error and no handlers are called
	if bestRoute == nil {
		console.Error("Could not find route to match: " + path)
		return
	}
	// Make the params map and pass it to the handler
	params := map[string]string{}
	for i, match := range bestMatches {
		params[bestRoute.paramNames[i]] = match
	}
	bestRoute.handler(params)
}

// removeEmptyStrings removes any empty strings from a
func removeEmptyStrings(a []string) []string {
	for i, s := range a {
		if s == "" {
			a = append(a[:i], a[i+1:]...)
		}
	}
	return a
}

// watchHash listens to the onhashchange event and calls r.pathChanged when
// it changes
func (r *Router) watchHash() {
	js.Global.Set("onhashchange", func() {
		go func() {
			path := getPathFromHash(getHash())
			r.pathChanged(path)
		}()
	})
}

// watchHistory listens to the onpopstate event and calls r.pathChanged when
// it changes
func (r *Router) watchHistory() {
	js.Global.Set("onpopstate", func() {
		go func() {
			r.pathChanged(getPath())
		}()
	})
}

// getPathFromHash returns everything after the "#" character in hash.
func getPathFromHash(hash string) string {
	return strings.SplitN(hash, "#", 2)[1]
}

// getHash is an alias for js.Global.Get("location").Get("hash").String()
func getHash() string {
	return js.Global.Get("location").Get("hash").String()
}

// setHash is an alias for js.Global.Get("location").Set("hash", hash)
func setHash(hash string) {
	js.Global.Get("location").Set("hash", hash)
}

// getPath is an alias for js.Global.Get("location").Get("pathname").String()
func getPath() string {
	return js.Global.Get("location").Get("pathname").String()
}

// pushState is an alias for js.Global.Get("history").Call("pushState", nil, "", path)
func pushState(path string) {
	js.Global.Get("history").Call("pushState", nil, "", path)
}
