package views

import (
	"fmt"
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
	"honnef.co/go/js/dom"
)

type Todo struct {
	humble.Identifier
	Model *models.Todo
}

func (t *Todo) RenderHTML() string {
	return fmt.Sprintf(`<li class="todo-list-item">
		<input class="toggle" type="checkbox" %s>
		<label class="todo-label">%s</label>
		<button class="destroy"></button>
		<input class="edit" onfocus="this.value = this.value;" value="%s">
		</li>`,
		t.Model.CheckedStr(), t.Model.Title, t.Model.Title)
}

func (t *Todo) OnLoad() error {
	var err error
	err = humble.Views.AddListener(t, "button.destroy", "click", t.deleteButtonClicked)
	if err != nil {
		panic(err)
	}
	err = humble.Views.AddListener(t, "label.todo-label", "dblclick", t.todoDoubleClick)
	if err != nil {
		panic(err)
	}
	err = humble.Views.AddListener(t, "input.edit", "keyup", t.todoEditKeyUp)
	if err != nil {
		panic(err)
	}
	err = humble.Views.AddListener(t, "input.toggle", "click", t.checkboxClicked)
	if err != nil {
		panic(err)
	}
	err = humble.Views.AddListener(t, "input.edit", "blur", t.todoEditBlurred)
	if err != nil {
		panic(err)
	}
	return nil
}

func (t *Todo) OuterTag() string {
	return "div"
}

func (t *Todo) todoEditKeyUp(event dom.Event) {
	// If key is not Enter or Escape, we keep label and input.edit in sync but otherwise just return
	key := event.(*dom.KeyboardEvent).KeyCode
	// Grab contents of input.edit
	inputEdit, err := humble.Views.QuerySelector(t, "input.edit")
	if err != nil {
		panic(err)
	}
	title := inputEdit.Underlying().Get("value").String()
	if key != EnterKey && key != EscapeKey {
		// Change everything in label to match input.edit
		label, err := humble.Views.QuerySelector(t, "label.todo-label")
		if err != nil {
			panic(err)
		}
		label.SetInnerHTML(title)
		return
	}
	// If Escape or Enter key is entered, we want to get out of input.edit field
	if key == EscapeKey || key == EnterKey {
		t.removeEditTodo()
	}
	// If Enter key is pressed, we want to save to model
	if key == EnterKey {
		t.Model.Title = title
		if err := humble.Models.Update(t.Model); err != nil {
			panic(err)
		}
	}
}

func (t *Todo) todoEditBlurred(dom.Event) {
	t.removeEditTodo()
}

func (t *Todo) removeEditTodo() {
	label, err := humble.Views.QuerySelector(t, "label.todo-label")
	if err != nil {
		panic(err)
	}
	todoItem, err := humble.Views.QuerySelector(t, "li.todo-list-item")
	if err != nil {
		panic(err)
	}
	//Remove 'editing' class to todoItem to make it disappear
	todoItem.Class().Remove("editing")
	// Show our label while input.edit is open
	label.(*dom.HTMLLabelElement).Style().SetProperty("display", "block", "important")
}

func (t *Todo) todoDoubleClick(dom.Event) {
	// Get elements
	label, err := humble.Views.QuerySelector(t, "label.todo-label")
	if err != nil {
		panic(err)
	}
	todoItem, err := humble.Views.QuerySelector(t, "li.todo-list-item")
	if err != nil {
		panic(err)
	}
	inputEdit, err := humble.Views.QuerySelector(t, "input.edit")
	if err != nil {
		panic(err)
	}
	// Hide our label while input.edit is open
	label.(*dom.HTMLLabelElement).Style().SetProperty("display", "none", "important")
	// Append 'editing' class to todoItem to make it an editable input field
	todoItem.Class().Add("editing")
	// Set focus to input.edit field
	inputEdit.(dom.HTMLElement).Focus()
}

func (t *Todo) deleteButtonClicked(dom.Event) {
	if err := humble.Views.Remove(t); err != nil {
		panic(err)
	}
	if err := humble.Models.Delete(t.Model); err != nil {
		panic(err)
	}
}

func (t *Todo) checkboxClicked(event dom.Event) {
	isChecked := event.Target().(*dom.HTMLInputElement).Checked
	t.Model.IsCompleted = isChecked
	if err := humble.Models.Update(t.Model); err != nil {
		panic(err)
	}

}
