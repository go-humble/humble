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
	appHasLoaded = false
)

func init() {
	elements.body = doc.QuerySelector(bodySelector)
}

func main() {
	console.Log("Starting...")

	// Get existing todos
	todos := []*models.Todo{}
	if err := model.ReadAll(&todos); err != nil {
		panic(err)
	}
	//Start main app view, appView
	appView := &views.App{}
	if err := view.ReplaceParentHTML(appView, bodySelector); err != nil {
		panic(err)
	}

	r := router.New()
	r.HandleFunc("/", func(params map[string]string) {
		appView.InitChildren(todos)
		if err := view.Update(appView); err != nil {
			panic(err)
		}
		if err := view.Update(appView.Footer); err != nil {
			panic(err)
		}
		appView.ApplyFilter(views.FilterAll)
	})
	r.HandleFunc("/active", func(params map[string]string) {
		appView.InitChildren(todos)
		if err := view.Update(appView); err != nil {
			panic(err)
		}
		if err := view.Update(appView.Footer); err != nil {
			panic(err)
		}
		appView.ApplyFilter(views.FilterActive)
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		appView.InitChildren(todos)
		if err := view.Update(appView); err != nil {
			panic(err)
		}
		if err := view.Update(appView.Footer); err != nil {
			panic(err)
		}
		appView.ApplyFilter(views.FilterCompleted)
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		console.Log("At Completed")
	})
	r.Start()

}

func initChildrenOnce(appView *views.App, todos []*models.Todo) {
	if !appHasLoaded {
		appHasLoaded = true
		appView.InitChildren(todos)
	}

}
