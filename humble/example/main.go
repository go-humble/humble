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

func (t *TodoView) OuterTag() string {
	return "div"
}

func main() {
	console.Log("Starting...")

	t := &TodoView{}

	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		console.Log("At home page")
	})
	r.HandleFunc("/append", func(params map[string]string) {
		console.Log(t.Id())
		humble.Views.AppendChild(t, nil, "#current-page")
	})
	r.HandleFunc("/about/{person_id}", func(params map[string]string) {
		console.Log("At person with ID: ", params["person_id"])
	})
	r.HandleFunc("/replace", func(params map[string]string) {
		console.Log(t.Id())
		humble.Views.SetOnlyChild(t, nil, "#current-page")
	})
	r.HandleFunc("/removeLast", func(params map[string]string) {
		console.Log(t.Id())
		humble.Views.Remove(t)
	})
	r.HandleFunc("/buy/purchase/{item_id}/image/{image_size}/panoramic", func(params map[string]string) {
		console.Log("Item ID:", params["item_id"], " Image_size", params["image_size"])
	})

	r.Start()

}
