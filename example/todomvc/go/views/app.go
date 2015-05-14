package views

import (
	"errors"
	"github.com/soroushjp/humble/temple"
	"github.com/soroushjp/humble/view"
)

type App struct {
	tmpl *temple.Template
	view.DefaultView
}

func NewApp() (*App, error) {
	tmpl, found := temple.Templates["app-template"]
	if !found {
		return nil, errors.New("Could not find template named app")
	}
	appView := &App{
		tmpl: tmpl,
	}
	appView.SetElement(document.QuerySelector("#todoapp"))
	return appView, nil
}

func (a *App) Render() error {
	return a.tmpl.Render(a.Element(), nil)
}
