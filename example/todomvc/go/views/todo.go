package views

import (
	"fmt"
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
	"github.com/gophergala/humble/model"
	"github.com/gophergala/humble/view"
	"honnef.co/go/js/dom"
)

type Todo struct {
	humble.Identifier
	Model  *models.Todo
	Parent *App
}

func (t *Todo) RenderHTML() string {
	return fmt.Sprintf(`<li class="todo-list-item %s">
		<input class="toggle" type="checkbox" %s>
		<label class="todo-label">%s</label>
		<button class="destroy"></button>
		<input class="edit" onfocus="this.value = this.value;" value="%s">
		</li>`,
		t.Model.CompletedStr(), t.Model.CheckedStr(), t.Model.Title, t.Model.Title)
}

func (t *Todo) OnLoad() error {
	if err := view.AddListener(t, "button.destroy", "click", t.deleteButtonClicked); err != nil {
		return err
	}
	if err := view.AddListener(t, "label.todo-label", "dblclick", t.todoDoubleClick); err != nil {
		return err
	}
	if err := view.AddListener(t, "input.edit", "keyup", t.todoEditKeyUp); err != nil {
		return err
	}
	if err := view.AddListener(t, "input.toggle", "click", t.checkboxClicked); err != nil {
		return err
	}
	if err := view.AddListener(t, "input.edit", "blur", t.todoEditBlurred); err != nil {
		return err
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
	inputEdit, err := view.QuerySelector(t, "input.edit")
	if err != nil {
		panic(err)
	}
	title := inputEdit.Underlying().Get("value").String()
	if key != EnterKey && key != EscapeKey {
		// Change everything in label to match input.edit
		label, err := view.QuerySelector(t, "label.todo-label")
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
		if err := model.Update(t.Model); err != nil {
			panic(err)
		}
	}
}

func (t *Todo) todoEditBlurred(dom.Event) {
	t.removeEditTodo()
}

func (t *Todo) removeEditTodo() {
	label, err := view.QuerySelector(t, "label.todo-label")
	if err != nil {
		panic(err)
	}
	todoItem, err := view.QuerySelector(t, "li.todo-list-item")
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
	label, err := view.QuerySelector(t, "label.todo-label")
	if err != nil {
		panic(err)
	}
	todoItem, err := view.QuerySelector(t, "li.todo-list-item")
	if err != nil {
		panic(err)
	}
	inputEdit, err := view.QuerySelector(t, "input.edit")
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
	t.remove()
}

func (t *Todo) remove() {
	if err := view.Remove(t); err != nil {
		panic(err)
	}
	if err := model.Delete(t.Model); err != nil {
		panic(err)
	}
	t.Parent.removeChild(t)
}

func (t *Todo) checkboxClicked(event dom.Event) {
	isChecked := event.Target().(*dom.HTMLInputElement).Checked
	t.setComplete(isChecked)
	if err := view.Update(t.Parent.Footer); err != nil {
		panic(err)
	}

	if t.Parent.CurrentFilter == FilterActive {
		switch isChecked {
		case true:
			// If we are only showing active todos and we just completed this one,
			// hide it
			if err := view.Hide(t); err != nil {
				panic(err)
			}
		case false:
			// If we are only showing active todos and we just uncompleted this on,
			// show it.
			if err := view.Show(t); err != nil {
				panic(err)
			}
		}
	} else if t.Parent.CurrentFilter == FilterCompleted {
		switch isChecked {
		case true:
			// If we are only showing completed todos and we just completed this on,
			// show it.
			if err := view.Show(t); err != nil {
				panic(err)
			}
		case false:
			// If we are only showing completed todos and we just uncompleted this on,
			// hide it.
			if err := view.Hide(t); err != nil {
				panic(err)
			}
		}
	}
}

func (t *Todo) setComplete(isCompleted bool) {
	t.Model.IsCompleted = isCompleted
	if err := model.Update(t.Model); err != nil {
		panic(err)
	}
	view.Update(t)
}
