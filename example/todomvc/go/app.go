package main

import (
	"github.com/soroushjp/humble/example/todomvc/go/models"
	"github.com/soroushjp/humble/example/todomvc/go/views"
	"github.com/soroushjp/humble/rest"
	"github.com/soroushjp/humble/temple"
	"github.com/soroushjp/humble/view"
	"honnef.co/go/js/dom"
	"log"
)

var (
	document  = dom.GetWindow().Document()
	body      = document.QuerySelector("body")
	todos     = []*models.Todo{}
	todoViews = []*views.Todo{}
)

func init() {
	log.SetFlags(log.Lmicroseconds)
	log.Print("Loading todos...")
	if err := rest.ReadAll(&todos); err != nil {
		log.Fatal(err)
	}
	log.Print("Loading templates...")
	if err := temple.Load(); err != nil {
		log.Fatal(err)
	}
	log.Print("Done loading.")
}

func main() {
	if err := renderApp(); err != nil {
		log.Fatal(err)
	}
	if err := renderTodos(); err != nil {
		log.Fatal(err)
	}
	if err := renderFooter(); err != nil {
		log.Fatal(err)
	}
}

func renderApp() error {
	appView, err := views.NewApp()
	if err != nil {
		return err
	}
	if err := appView.Render(); err != nil {
		return err
	}
	return nil
}

func renderTodos() error {
	todoList := document.QuerySelector("#todo-list")
	for _, todo := range todos {
		todoView, err := views.NewTodo(todo)
		if err != nil {
			return err
		}
		view.AppendToEl(todoList, todoView)
		if err := todoView.Render(); err != nil {
			return err
		}
		todoViews = append(todoViews, todoView)
	}
	return nil
}

func renderFooter() error {
	footerView, err := views.NewFooter(todos)
	if err != nil {
		return err
	}
	if err := footerView.Render(); err != nil {
		return err
	}
	return nil
}
