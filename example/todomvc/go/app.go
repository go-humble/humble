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

	todoListSelector  = "#todo-list"
	newTodoSelector   = "input#new-todo"
	toggleBtnSelector = "input#toggle-all"
)

var (
	doc      = dom.GetWindow().Document()
	elements = struct {
		todoList  dom.Element
		newTodo   dom.Element
		toggleBtn dom.Element
	}{}
)

// Listener is a callback function that will be triggered in response
// to some javascript event.
type Listener func(dom.Event)

func init() {
	elements.todoList = doc.QuerySelector(todoListSelector)
	elements.newTodo = doc.QuerySelector(newTodoSelector)
	elements.toggleBtn = doc.QuerySelector(toggleBtnSelector)
}

func main() {
	console.Log("Starting...")

	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		// Get existing todos
		todos := []*models.Todo{}
		if err := humble.Models.GetAll(&todos); err != nil {
			panic(err)
		}
		if len(todos) > 0 {
			showTodosContainer()
		}
		for _, todo := range todos {
			view := &views.Todo{
				Model: todo,
			}
			if err := humble.Views.AppendToParentHTML(view, todoListSelector); err != nil {
				panic(err)
			}
		}

		//Attach listener to newTodo input onkeyup event
		elements.newTodo.AddEventListener("keyup", false, nonBlockingListener(addTodoListener))
		//Attach listener to toggle list button onclick event
		elements.toggleBtn.AddEventListener("click", false, nonBlockingListener(toggleBtnListener))
	})
	r.HandleFunc("/completed", func(params map[string]string) {
		console.Log("At Completed")
	})
	r.Start()

}

// addTodoListener responds to DOM element input#new-todo being submitted by user to add a new todo to list and model
func addTodoListener(event dom.Event) {
	//If not Enter key, ignore event
	if event.(*dom.KeyboardEvent).KeyCode != EnterKey {
		return
	}
	//If newTodo input is empty, ignore event
	title := elements.newTodo.Underlying().Get("value").String()
	if title == "" {
		return
	}
	showTodosContainer()
	m := &models.Todo{
		Title:       title,
		IsCompleted: false,
	}
	if err := humble.Models.Create(m); err != nil {
		panic(err)
	}
	view := &views.Todo{
		Model: m,
	}
	if err := humble.Views.AppendToParentHTML(view, todoListSelector); err != nil {
		panic(err)
	}
}

// toggleBtnListener responds to DOM element input#toggle-all being clicked to trigger hide/show todo list
func toggleBtnListener(event dom.Event) {
	if elements.todoList.GetAttribute("style") == "" || elements.todoList.GetAttribute("style") == "null" {
		elements.todoList.SetAttribute("style", "visibility: hidden; height: 0;")
	} else {
		elements.todoList.SetAttribute("style", "")
	}
}

// showTodosContainer sets the outer container of todos to visible when our first todo is added
func showTodosContainer() {
	doc.QuerySelector("#main").SetAttribute("style", "display: block;")
}

// nonBlockingListener takes care of wrapping our event listener functions with a goroutine to make these usually
// blocking calls non-blocking, as required by GopherJS
func nonBlockingListener(listener Listener) Listener {
	return func(ev dom.Event) {
		go func() {
			listener(ev) //gopherjs:blocking
		}()
	}
}
