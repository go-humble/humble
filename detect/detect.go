package detect

import (
	"github.com/gopherjs/gopherjs/js"
)

// IsClient returns true iff the code that is currently running is
// compiled javascript code running in a browser with a global document
// property.
func IsClient() bool {
	return IsJavascript() && js.Global.Get("document") != js.Undefined
}

// IsJavascript return true iff the code that is currently running
// is compiled javascript code.
func IsJavascript() bool {
	return js.Global != nil
}

// IsServer returns true iff the code that is currently running is
// pure go code. That is, if the code that is currently running has
// not been compiled to javascript and is not running inside a browser.
func IsServer() bool {
	return !IsClient()
}
