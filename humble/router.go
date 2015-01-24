package humble

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
	"strings"
	"time"
)

// Router is responsible for handling routes
type Router struct {
	routes map[string]Handler
}

// Handler is a function which is run in response to a specific
// route. A Handler takes the url parameters as an argument.
type Handler func(params map[string]string)

// NewRouter creates and returns a new router
func NewRouter() *Router {
	return &Router{
		routes: map[string]Handler{},
	}
}

// HandleFunc will cause the router to call f whenever the
// hash of the url (everything after the '#' symbol) matches path.
func (r *Router) HandleFunc(path string, handler Handler) {
	r.routes[path] = handler
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
	// path is everything after the '#'
	path := strings.SplitN(hash, "#", 2)[1]
	if _, found := r.routes[path]; found {
		console.Log("Found handler for " + path)
		// TODO: parse url params and call handler func
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
