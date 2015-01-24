package main

import (
	"github.com/gophergala/go_by_the_fireplace/humble"
	"honnef.co/go/js/console"
)

type TodoView struct {
	humble.Identifier
}

func (t *TodoView) GetHTML(m humble.Model) string {
	return "<p>Bringing Go to the frontend :) please no more Javascript :(</p>"
}

func main() {
	console.Log("Starting...")

	t := &TodoView{}

	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		console.Log("At home page")
	})
	r.HandleFunc("/about", func(params map[string]string) {
		console.Log("At about page")
		humble.Views.Append(t, nil, "#current-page")
	})
	r.HandleFunc("/about/{person_id}", func(params map[string]string) {
		console.Log("At person with ID: ", params["person_id"])
	})
	r.HandleFunc("/faq", func(params map[string]string) {
		console.Log("At FAQ page")
	})
	r.HandleFunc("/buy/purchase/{item_id}/image/{image_size}/panoramic", func(params map[string]string) {
		console.Log("Item ID:", params["item_id"], " Image_size", params["image_size"])
	})

	r.Start()

}
