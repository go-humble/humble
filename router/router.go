package router

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
	"regexp"
	"strings"
	"time"
)

// Router is responsible for handling routes
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
	regex      *regexp.Regexp //Regex pattern that matches route
	paramNames []string       //Ordered list of query parameters expected by route handler
	handler    Handler        //Handler called when route is matched
}

// HandleFunc will cause the router to call f whenever the
// hash of the url (everything after the '#' symbol) matches path.
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

// Start will listen for changes in the hash of the url and
// trigger the appropriate handler function.
func (r *Router) Start() {
	r.setInitialHash()
	r.watchHash()
}

// setInitialHash will set hash to / when none is given
func (r *Router) setInitialHash() {
	if hash := getHash(); hash == "" {
		setHash("/")
	} else {
		r.hashChanged(hash)
	}
}

// hashChanged is called whenever DOM onhashchange event is fired
func (r *Router) hashChanged(hash string) {
	// path is everything after the '#'
	path := strings.SplitN(hash, "#", 2)[1]
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

	//If no routes match, we throw console error and no handlers are called
	if bestRoute == nil {
		console.Error("Could not find route to match: " + path)
		return
	}

	// Make the params map and pass it to the handler
	params := map[string]string{}
	for i, match := range bestMatches {
		params[bestRoute.paramNames[i]] = match
	}
	bestRoute.handler(params) //gopherjs:blocking
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

// watchHash watches DOM onhashchange and calls route.hashChanged
func (r *Router) watchHash() {
	if js.Global.Get("onhashchange") != js.Undefined {
		js.Global.Set("onhashchange", func() {
			go func() {
				r.hashChanged(getHash())
			}()
		})
	} else {
		console.Warn("onhashchange is not supported. Humble is falling back to a legacy version.")
		r.legacyWatchHash()
	}
}

// legacyWatchHash sets a ticker to check for hash changes every 50ms in browsers where onhashchange is not supported
func (r *Router) legacyWatchHash() {
	t := time.NewTicker(50 * time.Millisecond)
	go func() {
		hash := getHash()
		for {
			<-t.C
			newHash := getHash()
			if hash != newHash {
				hash = newHash
				r.hashChanged(hash)
			}
		}
	}()
}

// getHash gets DOM window.location.hash
func getHash() string {
	return js.Global.Get("location").Get("hash").String()
}

// setHash sets DOM window.location.hash to given hash
func setHash(hash string) {
	js.Global.Get("location").Set("hash", hash)
}
