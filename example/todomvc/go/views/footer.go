package views

import (
	"errors"
	"github.com/soroushjp/humble/example/todomvc/go/models"
	"github.com/soroushjp/humble/temple"
	"github.com/soroushjp/humble/view"
)

type Footer struct {
	Model models.Todos
	tmpl  *temple.Template
	view.DefaultView
}

func NewFooter(todos []*models.Todo) (*Footer, error) {
	tmpl, found := temple.Templates["footer-template"]
	if !found {
		return nil, errors.New("Could not find template named footer")
	}
	footerView := &Footer{
		Model: todos,
		tmpl:  tmpl,
	}
	footerView.SetElement(document.QuerySelector("footer#footer"))
	return footerView, nil
}

func (f *Footer) Render() error {
	return f.tmpl.Render(f.Element(), f.Model)
}
