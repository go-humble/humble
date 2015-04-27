package views

import (
	"errors"
	"github.com/soroushjp/humble/example/todomvc/go/models"
	"github.com/soroushjp/humble/temple"
	"github.com/soroushjp/humble/view"
)

type Todo struct {
	Model *models.Todo
	tmpl  *temple.Template
	view.DefaultView
}

func NewTodo(model *models.Todo) (*Todo, error) {
	tmpl, found := temple.Templates["todo-template"]
	if !found {
		return nil, errors.New("Could not find template named todo")
	}
	todoView := &Todo{
		Model: model,
		tmpl:  tmpl,
	}
	return todoView, nil
}

func (t *Todo) Render() error {
	return t.tmpl.Render(t.Element(), t.Model)
}
