package humble

// View is the interface that must be implemented by all views.
// GetHTML(Model) returns the HTML to be inserted into the DOM, given an object implementing Model.
// Id() sets the unique ID of the View object.
//// To be given a random unique id, simply include humble.Identifer as an anonymous field ie.
//// type ExampleView struct {
//// 	humble.Identifier
//// }
// OuterTag() sets the tag name for the outer container that will contain HTML returned from getHTML().
//// This is required, but can be simply "div" or "span" for a semantically neutral HTML element.

import (
	"fmt"
	"honnef.co/go/js/dom"
	"regexp"
)

type View interface {
	GetHTML(Model) string
	Id() string
	OuterTag() string
}

type viewsType struct{}

var viewsIndex = map[string]*dom.Element{}
var Views = viewsType{}
var document = dom.GetWindow().Document()

// doSomething does something
func (v *View) doSomething() {

}

// AppendChild appends a view as a child to a parent DOM element. It takes a View interface, a Model provided to the view
// and a parent DOM selector. Parent selector works identically to JavaScript's document.querySelector(selector) call.
func (*viewsType) AppendChild(view View, model Model, parentSelector string) error {
	//Grab DOM element matching parentSelector
	parent := document.QuerySelector(parentSelector)
	if parent == nil {
		return fmt.Errorf("Could not find element for parentSelector: %s", parentSelector)
	}
	//Get our view HTML given the model
	html := view.GetHTML(model)
	if html == "" {
		return nil
	}
	//Check our outer container tag is valid
	if err := checkOuterTag(view.OuterTag()); err != nil {
		return err
	}
	//Create our element to append, with outer tag
	el := document.CreateElement(view.OuterTag())
	if _, found := viewsIndex[view.Id()]; found {
		return fmt.Errorf("Duplicate humble.View Id: %s", view.Id())
	}
	viewsIndex[view.Id()] = &el
	el.SetInnerHTML(html)
	//We set attribute data-humble-view-id on outer container to simplify debugging and as a secondary means of
	//selecting our View element from the DOM
	el.SetAttribute("data-humble-view-id", view.Id())
	//Append as child to selected parent DOM element
	parent.AppendChild(el)

	return nil
}

// checkOuterTag will check that the given HTML tag is composed of alphabetical characters
func checkOuterTag(tag string) error {
	match, err := regexp.Match("[a-zA-Z]", []byte(tag))
	if err != nil {
		fmt.Errorf("Invalid outer tag for humble.View: %s", err.Error())
	}
	if !match {
		return fmt.Errorf("Outer tag must be alphabetical characters")
	}
	return nil
}
