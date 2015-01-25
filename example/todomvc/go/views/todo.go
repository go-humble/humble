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

func (t *Todo) GetHTML() string {
	return fmt.Sprintf(`<li>
		<input class="toggle" type="checkbox" %s>
		<label>%s</label>
		<button class="destroy"></button>
		<input class="edit" value="%s">
		</li>`,
		t.Model.CheckedStr(), t.Model.Title, t.Model.Title)
}

func (t *Todo) OnLoad() error {
	err := humble.Views.AddListener(t, "button.destroy", "click", t.deleteButtonClicked)
	if err != nil {
		panic(err)
	}
	return nil
}

func (t *Todo) OuterTag() string {
	return "div"
}

func (t *Todo) deleteButtonClicked(dom.Event) {
	if err := humble.Views.Remove(t); err != nil {
		panic(err)
	}
	if err := humble.Models.Delete(t.Model); err != nil {
		panic(err)
	}
}
