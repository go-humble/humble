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
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
	"regexp"
)

type View interface {
	GetHTML() string
	Id() string
	OuterTag() string
}

type viewsType struct{}

var viewsIndex = map[string]*dom.Element{}
var Views = viewsType{}
var document dom.Document

func init() {
	// If we are running this code in a test runner, document is undefined.
	// We only want to initialize document if we are running in the browser.
	if js.Global.Get("document") != js.Undefined {
		document = dom.GetWindow().Document()
	}
}

// AppendChild appends a view as a child to a parent DOM element. It takes a View interface and
// a parent DOM selector. parentSelector works identically to JavaScript's document.querySelector(selector)
// call.
func (*viewsType) AppendChild(view View, parentSelector string) error {
	//Grab DOM element matching parentSelector
	parent := document.QuerySelector(parentSelector)
	if parent == nil {
		return fmt.Errorf("Could not find element for parentSelector: %s", parentSelector)
	}
	//Get our view HTML
	viewHTML := view.GetHTML()
	if viewHTML == "" {
		return nil
	}
	//Create our child DOM element
	viewEl, err := createViewElement(viewHTML, view.OuterTag(), view.Id())
	if err != nil {
		return err
	}
	//Append as child to selected parent DOM element
	parent.AppendChild(viewEl)

	return nil
}

// SetInnerHTML replaces the current contents of the parent DOM element.
// It takes a View interface and a parent DOM selector. parentSelector works identically to
// JavaScript's document.querySelector(selector) call.
func (*viewsType) SetInnerHTML(view View, parentSelector string) error {
	//Grab DOM element matching parentSelector
	parent := document.QuerySelector(parentSelector)
	if parent == nil {
		return fmt.Errorf("Could not find element for parentSelector: %s", parentSelector)
	}
	//Get our view HTML
	viewHTML := view.GetHTML()
	if viewHTML == "" {
		return nil
	}

	//Create our view DOM element
	viewEl, err := createViewElement(viewHTML, view.OuterTag(), view.Id())
	if err != nil {
		return err
	}
	//Append as child to selected parent DOM element
	parent.SetInnerHTML("")
	parent.AppendChild(viewEl)
	return nil
}

// Remove removes a view element from the DOM, returning true if successful, false otherwise
func (*viewsType) Remove(view View) bool {
	viewElRef, found := viewsIndex[view.Id()]
	if found {
		(*viewElRef).ParentElement().RemoveChild(*viewElRef)
		delete(viewsIndex, view.Id())
		return true
	}
	return false
}

// createChildElement creates a DOM element from HTML and a outer container tag.
// Takes innerHTML and outerTag, crafts a valid *dom.Element and adds it to the global map viewsIndex
// for easy referencing. Returns the resultant *dom.Element or an error.
func createViewElement(innerHTML string, outerTag string, viewId string) (dom.Element, error) {
	//Check our outer container tag is valid
	if err := checkOuterTag(outerTag); err != nil {
		return nil, err
	}
	//Create our element to append, with outer tag
	var el dom.Element
	//Create unique element ID for the view element and add it to global map of existent view elements viewsIndex
	if indexedEl, found := viewsIndex[viewId]; found {
		el = (*indexedEl)
	} else {
		el = document.CreateElement(outerTag)
	}
	viewsIndex[viewId] = &el
	el.SetInnerHTML(innerHTML)
	//We set attribute data-humble-view-id on outer container to simplify debugging and as a secondary means of
	//selecting our View element from the DOM
	el.SetAttribute("data-humble-view-id", viewId)

	return el, nil
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
