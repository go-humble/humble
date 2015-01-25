package main

import (
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
	"github.com/gophergala/humble/example/todomvc/go/views"
	"honnef.co/go/js/console"
	"honnef.co/go/js/dom"
)

const (
	bodySelector = "body"
)

var (
	doc      = dom.GetWindow().Document()
	elements = struct {
		body      dom.Element
		todoList  dom.Element
		newTodo   dom.Element
		toggleBtn dom.Element
	}{}
)

func init() {
	elements.body = doc.QuerySelector(bodySelector)
}

func main() {
	console.Log("Starting...")

	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		// Get existing todos
		todos := []*models.Todo{}
		if err := humble.Models.ReadAll(&todos); err != nil {
			panic(err)
		}
		//Start main app view, appView
		appView := &views.App{
			Model: todos,
		}
		humble.Views.AppendToParentHTML(appView, bodySelector)
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		console.Log("At Completed")
	})
	r.Start()

}
