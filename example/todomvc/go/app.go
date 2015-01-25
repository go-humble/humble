package main

import (
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
	"github.com/gophergala/humble/example/todomvc/go/views"
	"honnef.co/go/js/console"
	"honnef.co/go/js/dom"
)

const (
	EnterKey  = 13
	EscapeKey = 27

	todoListSelector = "#todo-list"
	newTodoSelector  = "#new-todo"
)

var (
	doc      = dom.GetWindow().Document()
	elements = struct {
		todoList dom.Element
		newTodo  dom.Element
	}{}
)

func main() {
	console.Log("Starting...")
	r := humble.NewRouter()

	err := humble.Models.Create(&models.Todo{
		Title:       "blahsdas&title=asadas !!!",
		IsCompleted: true,
	})
	if err != nil {
		panic(err)
	}
	r.HandleFunc("/", func(params map[string]string) {
		// Get existing todos
		todos := []*models.Todo{}
		if err := humble.Models.GetAll(&todos); err != nil {
			panic(err)
		}
		if len(todos) > 0 {
			showTodos()
		}
		for _, todo := range todos {
			view := &views.Todo{
				Model: todo,
			}
			if err := humble.Views.AppendToParentHTML(view, todoListSelector); err != nil {
				panic(err)
			}
		}
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		console.Log("At Completed")
	})
	r.Start()
}

func showTodos() {
	doc.QuerySelector("#main").SetAttribute("style", "display: block;")
}
