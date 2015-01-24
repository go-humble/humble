package humble

import (
	"fmt"
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

// NewRouter creates and returns a new router
func NewRouter() *Router {
	return &Router{
		routes: []*route{},
	}
}

type route struct {
	regex      *regexp.Regexp
	paramNames []string
	handler    Handler
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
	removeEmptyStrings(strs)
	pattern := `^`
	for _, str := range strs {
		if !(str[0] == '{' && str[len(str)-1] == '}') {
			pattern += `/`
			pattern += `([\w-+]*)`
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
	// TODO detect support for onhashchange
	// and use legacy version if not supported
	r.watchHash()
}

func (r *Router) setInitialHash() {
	if hash := getHash(); hash == "" {
		setHash("/")
	} else {
		r.hashChanged(hash)
	}
}

func (r *Router) hashChanged(hash string) {
	fmt.Printf("hashChanged(%s)\n", hash)

	// path is everything after the '#'
	path := strings.SplitN(hash, "#", 2)[1]
	strs := strings.Split(path, "/")
	removeEmptyStrings(strs)
	fmt.Println("Parsed hash correctly!")
	// canditateRoutes := []route{}
	// copy(canditateRoutes, r.routes)

	// // Eliminate all the routes which want a different
	// // number of tokens
	// for j, route := range canditateRoutes {
	// 	if len(route.tokens) != len(strs) {
	// 		canditateRoutes = append(canditateRoutes[:j], canditateRoutes[j+1:]...)
	// 	}
	// }

	// // Eliminate other routes one by one
	// score := []int{}
	// for i, str := range strs {
	// 	for j, route := range canditateRoutes {
	// 		if route.tokens[i].kind == literalToken {
	// 			if str != route.tokens[i].name {
	// 				// If a route expected a literal token and the str
	// 				// we got didn't match, eliminate that route from
	// 				// the list of candidates
	// 				canditateRoutes = append(canditateRoutes[:j], canditateRoutes[j+1:]...)
	// 			}
	// 		}
	// 	}
	// }
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

	// If there is more than one candidate, throw an error
	// otherwise match the route with the one we left

	// todos/create
	// todos/{id}

	// todos/create

	// if _, found := r.routes[path]; found {
	// 	console.Log("Found handler for " + path)
	// 	// TODO: parse url params and call handler func
	// }
}

// removeEmptyStrings removes any empty strings from a
func removeEmptyStrings(a []string) {
	for i, s := range a {
		if s == "" {
			a = append(a[:i], a[i+1:]...)
		}
	}
}

func (r *Router) watchHash() {
	js.Global.Set("onhashchange", func() {
		r.hashChanged(getHash())
	})
}

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

func getHash() string {
	return js.Global.Get("location").Get("hash").String()
}

func setHash(string) {
	js.Global.Get("location").Set("hash", "/")
}
