package models

import (
	"fmt"
	"strconv"
)

type Todo struct {
	Id          int
	Title       string
	IsCompleted bool
}

func (t *Todo) GetId() string {
	return strconv.Itoa(t.Id)
}

func (t *Todo) UrlRoot() string {
	return "http://localhost:3000/todos"
}

func (t *Todo) checkedStr() string {
	if t.IsCompleted {
		return "checked"
	}
	return ""
}

func (t *Todo) innerHtml() string {
	return fmt.Sprintf(`
		<input class="toggle" type="checkbox" %s>
		<label>%s</label>
		<button class="destroy"></button>
		<input class="edit" value="%s">
		`,
		t.checkedStr(), t.Title, t.Title)
}
