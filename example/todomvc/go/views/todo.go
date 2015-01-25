package views

import (
	"fmt"
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/example/todomvc/go/models"
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

func (t *Todo) OuterTag() string {
	return "div"
}
