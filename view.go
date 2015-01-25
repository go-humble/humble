package humble

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
	"regexp"
)

// View is the interface that must be implemented by all views.
// GetHTML() returns the HTML to be inserted into the DOM.
// GetId() sets the unique ID of the View object.
//// To be given a random unique id, simply include humble.Identifer as an anonymous field ie.
//// type ExampleView struct {
//// 	humble.Identifier
//// }
// OuterTag() sets the tag name for the outer container that will contain HTML returned from getHTML().
//// This is required, but can be simply "div" or "span" for a semantically neutral HTML element.
type View interface {
	GetHTML() string
	GetId() string
	OuterTag() string
}

// If a view implements OnLoader, humble will call the OnLoad method
// whenver the view's element is added to (or updated in) the DOM.
type OnLoader interface {
	OnLoad() error
}

type viewsType struct{}

var viewsIndex = map[string]dom.Element{}
var Views = viewsType{}
var document dom.Document

func init() {
	// If we are running this code in a test runner, document is undefined.
	// We only want to initialize document if we are running in the browser.
	if js.Global.Get("document") != js.Undefined {
		document = dom.GetWindow().Document()
	}
}

// AppendToParentHTML appends a view to a parent DOM element. It takes a View interface and
// a parent DOM selector. parentSelector works identically to JavaScript's document.querySelector(selector)
// call. After the view's element is added to the DOM, AppendToParentHTML calls view.OnLoad if it is defined.
func (*viewsType) AppendToParentHTML(view View, parentSelector string) error {
	// Grab DOM element matching parentSelector
	parent := document.QuerySelector(parentSelector)
	if parent == nil {
		return fmt.Errorf("Could not find element for parentSelector: %s", parentSelector)
	}
	// Create our child DOM element
	viewEl, err := createViewElement(view)
	if err != nil {
		return err
	}
	// Append as child to selected parent DOM element
	parent.AppendChild(viewEl)

	// Call view.OnLoad if it is defined
	if err := viewOnLoad(view); err != nil {
		return err
	}
	return nil
}

// ReplaceParentHTML replaces the current inner HTML of the parent DOM element with the view.
// It takes a View interface and a parent DOM selector. parentSelector works identically to
// JavaScript's document.querySelector(selector) call. After the view's element is added to the
// DOM, ReplaceParentHTML calls view.OnLoad if it is defined.
func (*viewsType) ReplaceParentHTML(view View, parentSelector string) error {
	// Grab DOM element matching parentSelector
	parent := document.QuerySelector(parentSelector)
	if parent == nil {
		return fmt.Errorf("Could not find element for parentSelector: %s", parentSelector)
	}
	// Create our view DOM element
	viewEl, err := createViewElement(view)
	if err != nil {
		return err
	}
	// Append as child to selected parent DOM element
	parent.SetInnerHTML("")
	parent.AppendChild(viewEl)

	// Call view.OnLoad if it is defined
	if err := viewOnLoad(view); err != nil {
		return err
	}
	return nil
}

// Update updates a view in place by calling SetInnerHTML on the view's element.
// Returns an error if the dom element for this view does not exist. After the view's
// element is added to the DOM, Update calls view.OnLoad if it is defined.
func (*viewsType) Update(view View) error {
	html := view.GetHTML()
	el, err := getElementByViewId(view.GetId())
	if err != nil {
		return err
	}
	el.SetInnerHTML(html)

	// Call view.OnLoad if it is defined
	if err := viewOnLoad(view); err != nil {
		return err
	}
	return nil
}

// Remove removes a view element from the DOM, returning true if successful, false otherwise
func (*viewsType) Remove(view View) error {
	viewEl, err := getElementByViewId(view.GetId())
	if err != nil {
		return err
	}
	viewEl.ParentElement().RemoveChild(viewEl)
	return nil
}

func getElementByViewId(viewId string) (dom.Element, error) {
	if indexedEl, found := viewsIndex[viewId]; found {
		return indexedEl, nil
	} else {
		// The element wasn't in our index. Try finding in the DOM as a last
		// resort. (Maybe our index got out of date because the DOM was changed
		// outside of humble).
		selector := fmt.Sprintf("[data-humble-view-id='%s']", viewId)
		el := document.QuerySelector(selector)
		if el == nil {
			return nil, ViewElementNotFoundError{viewId: viewId}
		}
		viewsIndex[viewId] = el //Add our element to index since it exists in DOM but was not found in index
		return el, nil
	}
}

// createViewElement creates a DOM element from HTML and a outer container tag.
// Takes innerHTML and outerTag, crafts a valid *dom.Element and adds it to the global map viewsIndex
// for easy referencing. Returns the resultant *dom.Element or an error.
func createViewElement(view View) (dom.Element, error) {
	//Check our outer container tag is valid
	if err := checkOuterTag(view.OuterTag()); err != nil {
		return nil, err
	}
	//Get our view HTML
	viewHTML := view.GetHTML()
	//Check if view element exists in global map, otherwise create it
	var el dom.Element
	if indexedEl, err := getElementByViewId(view.GetId()); err != nil {
		if _, notFound := err.(ViewElementNotFoundError); notFound {
			// The view was not found in the DOM. We need to create it
			el = document.CreateElement(view.OuterTag())
			viewsIndex[view.GetId()] = el
		} else {
			// For any other type of error, return it.
			return nil, err
		}
	} else {
		el = indexedEl
	}
	el.SetInnerHTML(viewHTML)
	//We set attribute data-humble-view-id on outer container to simplify debugging and as a secondary means of
	//selecting our View element from the DOM
	el.SetAttribute("data-humble-view-id", view.GetId())

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

func viewOnLoad(v View) error {
	if onLoader, hasOnLoad := v.(OnLoader); hasOnLoad {
		if err := onLoader.OnLoad(); err != nil {
			return fmt.Errorf("Error in %T.OnLoad: %s", v, err.Error())
		}
	}
	return nil
}
