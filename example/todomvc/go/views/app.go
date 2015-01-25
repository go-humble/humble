package views

import (
	"fmt"
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
	"honnef.co/go/js/dom"
)

type App struct {
	humble.Identifier
	Model []*models.Todo
}

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
		body      dom.Element
		todoList  dom.Element
		newTodo   dom.Element
		toggleBtn dom.Element
	}{}
)

func (a *App) RenderHTML() string {
	return fmt.Sprintf(`
	<section id="todoapp">
		<header id="header">
			<h1>todos</h1>
			<input id="new-todo" placeholder="What needs to be done?" autofocus>
		</header>
		<section id="main">
			<input id="toggle-all" type="checkbox">
			<label for="toggle-all">Mark all as complete</label>
			<ul id="todo-list"></ul>
		</section>
		<footer id="footer"></footer>
	</section>
	<footer id="info">
		<p>Double-click to edit a todo</p>
		<p>Part of <a href="http://todomvc.com">TodoMVC</a>
		</p>
	</footer>
	<script src="js/app.js"></script>`)
}

func (v *App) OuterTag() string {
	return "div"
}

func (v *App) OnLoad() error {
	var err error
	elements.todoList, err = humble.Views.ChildQuerySelector(v, todoListSelector)
	if err != nil {
		return err
	}
	elements.newTodo, err = humble.Views.ChildQuerySelector(v, newTodoSelector)
	if err != nil {
		return err
	}
	elements.toggleBtn, err = humble.Views.ChildQuerySelector(v, toggleBtnSelector)
	if err != nil {
		return err
	}

	if len(v.Model) > 0 {
		showTodosContainer()
	}

	//Create individual todo views
	for _, todo := range v.Model {
		view := &Todo{
			Model: todo,
		}
		if err := humble.Views.AppendToParentHTML(view, todoListSelector); err != nil {
			panic(err)
		}
	}

	err = humble.Views.AddListener(v, newTodoSelector, "keyup", v.addTodoListener)
	if err != nil {
		return err
	}
	err = humble.Views.AddListener(v, toggleBtnSelector, "click", v.toggleBtnListener)
	if err != nil {
		return err
	}

	// //Attach listener to newTodo input onkeyup event
	// elements.newTodo.AddEventListener("keyup", false, nonBlockingListener(addTodoListener))
	// //Attach listener to toggle list button onclick event
	// elements.toggleBtn.AddEventListener("click", false, nonBlockingListener(toggleBtnListener))

	return nil
}

// addTodoListener responds to DOM element input#new-todo being submitted by user to add a new todo to list and model
func (v *App) addTodoListener(event dom.Event) {
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
	view := &Todo{
		Model: m,
	}
	if err := humble.Views.AppendToParentHTML(view, todoListSelector); err != nil {
		panic(err)
	}
}

// toggleBtnListener responds to DOM element input#toggle-all being clicked to trigger hide/show todo list
func (v *App) toggleBtnListener(event dom.Event) {

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
