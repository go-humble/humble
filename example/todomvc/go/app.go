package main

import (
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
	"honnef.co/go/js/console"
	"honnef.co/go/js/dom"
)

const (
	EnterKey  = 13
	EscapeKey = 27
)

var (
	doc      = dom.GetWindow().Document()
	elements = struct {
		todoList dom.Element
		newTodo  dom.Element
	}{
		todoList: doc.QuerySelector("#todo-list"),
		newTodo:  doc.QuerySelector("#new-todo"),
	}
)

func main() {
	console.Log("Starting...")
	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		// Get existing todos
		todos := []*models.Todo{}
		if err := humble.Models.GetAll(&todos); err != nil {
			panic(err)
		}
		for _, todo := range todos {
			console.Log(todo)
		}
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		console.Log("At Completed")
	})
	r.Start()
}
