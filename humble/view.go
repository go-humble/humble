package humble

import (
	"fmt"
	"honnef.co/go/js/dom"
)

type View interface {
	GetHTML(Model) string
	Id() string
}

type viewsType struct{}

var Views = viewsType{}
var document = dom.GetWindow().Document()

func (*viewsType) Append(view View, model Model, domSelector string) error {
	//Grab DOM element matching domSelector
	parent := document.QuerySelector(domSelector)
	if parent == nil {
		return fmt.Errorf("Could not find element for domSelector: %s", domSelector)
	}
	html := view.GetHTML(model)
	if html == "" {
		return nil
	}
	parent.SetInnerHTML(parent.InnerHTML() + html)
	return nil
}
