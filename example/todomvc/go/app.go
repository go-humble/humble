package main

import (
	"github.com/gophergala/humble/example/todomvc/go/models"
	"github.com/gophergala/humble/example/todomvc/go/views"
	"github.com/gophergala/humble/model"
	"github.com/gophergala/humble/router"
	"github.com/gophergala/humble/view"
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

	r := router.New()
	r.HandleFunc("/", func(params map[string]string) {
		// Get existing todos
		todos := []*models.Todo{}
		if err := model.ReadAll(&todos); err != nil {
			panic(err)
		}
		//Start main app view, appView
		appView := &views.App{}
		appView.InitChildren(todos)
		if err := view.ReplaceParentHTML(appView, bodySelector); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/active", func(params map[string]string) {
		// Get existing todos
		todos := []*models.Todo{}
		if err := model.ReadAll(&todos); err != nil {
			panic(err)
		}

		// Get only those todos which are active
		activeTodos := []*models.Todo{}
		for _, todo := range todos {
			if !todo.IsCompleted {
				activeTodos = append(activeTodos, todo)
			}
		}

		//Start main app view, appView
		appView := &views.App{}
		appView.InitChildren(activeTodos)
		if err := view.ReplaceParentHTML(appView, bodySelector); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		// Get existing todos
		todos := []*models.Todo{}
		if err := model.ReadAll(&todos); err != nil {
			panic(err)
		}

		// Get only those todos which are completed
		completedTodos := []*models.Todo{}
		for _, todo := range todos {
			if todo.IsCompleted {
				completedTodos = append(completedTodos, todo)
			}
		}

		//Start main app view, appView
		appView := &views.App{}
		appView.InitChildren(completedTodos)
		if err := view.ReplaceParentHTML(appView, bodySelector); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		console.Log("At Completed")
	})
	r.Start()

}
