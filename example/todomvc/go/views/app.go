package views

import (
	"fmt"
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
	"github.com/gophergala/humble/model"
	"github.com/gophergala/humble/view"
	"honnef.co/go/js/dom"
)

type App struct {
	humble.Identifier
	Model    []*models.Todo
	Children []*Todo
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
	elements.todoList, err = view.QuerySelector(v, todoListSelector)
	if err != nil {
		return err
	}
	elements.newTodo, err = view.QuerySelector(v, newTodoSelector)
	if err != nil {
		return err
	}
	elements.toggleBtn, err = view.QuerySelector(v, toggleBtnSelector)
	if err != nil {
		return err
	}

	if len(v.Model) > 0 {
		showTodosContainer()
	}

	//Create individual todo views
	for _, todo := range v.Model {
		todoView := &Todo{
			Model:  todo,
			Parent: v,
		}
		v.addChild(todoView)
		if err := view.AppendToParentHTML(todoView, todoListSelector); err != nil {
			return err
		}
	}

	if err := view.AddListener(v, newTodoSelector, "keyup", v.newTodoKeyUp); err != nil {
		return err
	}
	if err := view.AddListener(v, toggleBtnSelector, "click", v.toggleBtnClicked); err != nil {
		return err
	}

	return nil
}

func (v *App) removeChild(todoView *Todo) {
	for i, child := range v.Children {
		if child.Id == todoView.Id {
			v.Children = append(v.Children[:i], v.Children[i+1:]...)
		}
	}
}

func (v *App) addChild(todoView *Todo) {
	v.Children = append(v.Children, todoView)
}

// addTodoListener responds to DOM element input#new-todo being submitted by user to add a new todo to list and model
func (v *App) newTodoKeyUp(event dom.Event) {
	//If not Enter key, ignore event
	if event.(*dom.KeyboardEvent).KeyCode != EnterKey {
		return
	}
	//If newTodo input is empty, ignore event
	title := elements.newTodo.Underlying().Get("value").String()
	if title == "" {
		return
	}
	//This ensures the todo list container is visible. Does nothing if already visible, but costs no more than a check.
	showTodosContainer()
	//Create a model, send to server and append view
	m := &models.Todo{
		Title:       title,
		IsCompleted: false,
	}
	if err := model.Create(m); err != nil {
		panic(err)
	}
	todoView := &Todo{
		Model:  m,
		Parent: v,
	}
	if err := view.AppendToParentHTML(todoView, todoListSelector); err != nil {
		panic(err)
	}
	v.addChild(todoView)
	//Clear newTodo text input
	elements.newTodo.Underlying().Set("value", "")
}

// toggleBtnListener responds to DOM element input#toggle-all being clicked to toggle all todo
// items between the completed and active states.
func (v *App) toggleBtnClicked(event dom.Event) {
	isChecked := event.Target().(*dom.HTMLInputElement).Checked
	fmt.Println(isChecked)
	for _, todo := range v.Children {
		todo.setComplete(isChecked)
	}
}

// showTodosContainer sets the outer container of todos to visible when our first todo is added
func showTodosContainer() {
	doc.QuerySelector("#main").SetAttribute("style", "display: block;")
}
