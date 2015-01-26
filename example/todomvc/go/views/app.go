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
	Children      []*Todo
	Footer        *Footer
	CurrentFilter TodoFilter
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

type TodoFilter int

const (
	FilterAll TodoFilter = iota
	FilterActive
	FilterCompleted
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
		<footer id="footer">
		</footer>
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

func (v *App) InitChildren(todos []*models.Todo) {
	//Create individual todo views
	v.Children = []*Todo{}
	if v.Footer != nil {
		v.Footer.TodoViews = &v.Children
	}
	for _, todo := range todos {
		todoView := &Todo{
			Model:  todo,
			Parent: v,
		}
		v.addChild(todoView)
	}
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

	// Add each child view to the DOM
	for _, childView := range v.Children {
		view.AppendToParentHTML(childView, todoListSelector)
	}

	// Set up the footer view
	v.Footer = &Footer{}
	v.Footer.TodoViews = &v.Children
	view.ReplaceParentHTML(v.Footer, "#footer")

	if len(v.Children) > 0 {
		showTodosContainer()
		showTodosFooter()
	}

	if err := view.AddListener(v, newTodoSelector, "keyup", v.newTodoKeyUp); err != nil {
		return err
	}
	if err := view.AddListener(v, toggleBtnSelector, "click", v.toggleBtnClicked); err != nil {
		return err
	}

	return nil
}

func (v *App) ApplyFilter(filter TodoFilter) {
	v.CurrentFilter = filter
	for _, todoView := range v.Children {
		if todoView.GetId() == "" {
			continue
		}
		switch filter {
		case FilterAll:
			// For FilterAll we want to show all todos, regardless of whether they are complete
			if err := view.Show(todoView); err != nil {
				panic(err)
			}
		case FilterActive:
			// For the FilterActive, we want to hide views that are completed
			switch todoView.Model.IsCompleted {
			case true:
				if err := view.Hide(todoView); err != nil {
					panic(err)
				}
			case false:
				if err := view.Show(todoView); err != nil {
					panic(err)
				}
			}
		case FilterCompleted:
			// For the FilterCompleted, we want to hide views that are completed
			switch todoView.Model.IsCompleted {
			case true:
				if err := view.Show(todoView); err != nil {
					panic(err)
				}
			case false:
				if err := view.Hide(todoView); err != nil {
					panic(err)
				}
			}
		}
	}
}

func (v *App) removeChild(todoView *Todo) {
	for i, child := range v.Children {
		if child.Id == todoView.Id {
			v.Children = append(v.Children[:i], v.Children[i+1:]...)
		}
	}
	// Update the footer text
	if err := view.Update(v.Footer); err != nil {
		panic(err)
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
	showTodosFooter()
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

	// Update the footer text
	if err := view.Update(v.Footer); err != nil {
		panic(err)
	}
}

// toggleBtnListener responds to DOM element input#toggle-all being clicked to toggle all todo
// items between the completed and active states.
func (v *App) toggleBtnClicked(event dom.Event) {
	isChecked := event.Target().(*dom.HTMLInputElement).Checked
	for _, todo := range v.Children {
		todo.setComplete(isChecked)
	}

	if v.CurrentFilter == FilterActive {
		switch isChecked {
		case true:
			// If we are only showing active todos and we just completed all of them,
			// hide all the views
			for _, todoView := range v.Children {
				if err := view.Hide(todoView); err != nil {
					panic(err)
				}
			}
		case false:
			// If we are only showing active todos and we just uncompleted all of them,
			// show all the views
			for _, todoView := range v.Children {
				if err := view.Show(todoView); err != nil {
					panic(err)
				}
			}
		}
	} else if v.CurrentFilter == FilterCompleted {
		switch isChecked {
		case true:
			// If we are only showing completed todos and we just completed all of them,
			// show all the views
			for _, todoView := range v.Children {
				if err := view.Show(todoView); err != nil {
					panic(err)
				}
			}
		case false:
			// If we are only showing completed todos and we just uncompleted all of them,
			// hide all the views
			for _, todoView := range v.Children {
				if err := view.Hide(todoView); err != nil {
					panic(err)
				}
			}
		}
	}

	// Update the footer text
	if err := view.Update(v.Footer); err != nil {
		panic(err)
	}
}

// showTodosContainer sets the outer container of todos to visible when our first todo is added
func showTodosContainer() {
	doc.QuerySelector("#main").SetAttribute("style", "display: block;")
}

// showTodosContainer sets the outer todos footer (which contains links and the number of items left)
// to visible when our first todo is added
func showTodosFooter() {
	doc.QuerySelector("#footer").SetAttribute("style", "display: block;")
}
