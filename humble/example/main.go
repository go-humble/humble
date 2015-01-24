package main

import (
	"github.com/gophergala/go_by_the_fireplace/humble"
	"honnef.co/go/js/console"
)

func main() {
	console.Log("Starting...")

	r := humble.NewRouter()
	// r.HandleFunc("/", func(params map[string]string) {
	// 	console.Log("At home page")
	// })
	r.HandleFunc("/about", func(params map[string]string) {
		console.Log("At about page")
	})
	// r.HandleFunc("/faq", func(params map[string]string) {
	// 	console.Log("At FAQ page")
	// })

	r.Start()
}
