package view

import (
	"fmt"
	"github.com/gophergala/humble"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
	"regexp"
)

var document dom.Document

func init() {
	// If we are running this code in a test runner, document is undefined.
	// We only want to initialize document if we are running in the browser.
	if js.Global.Get("document") != js.Undefined {
		document = dom.GetWindow().Document()
	}
}

// View is the interface that must be implemented by all views.
// RenderHTML() returns the HTML to be inserted into the DOM.
// GetId() sets the unique ID of the View object.
//// To be given a random unique id, simply include humble.Identifer as an anonymous field ie.
//// type ExampleView struct {
//// 	humble.Identifier
//// }
// OuterTag() sets the tag name for the outer container that will contain HTML returned from getHTML().
//// This is required, but can be simply "div" or "span" for a semantically neutral HTML element.
type View interface {
	RenderHTML() string
	GetId() string
	OuterTag() string
}

// If a view implements OnLoader, humble will call the OnLoad method
// whenver the view's element is added to (or updated in) the DOM.
type OnLoader interface {
	OnLoad() error
}

// Listener is a callback function that will be triggered in response
// to some javascript event.
type Listener func(dom.Event)

// AppendToParentHTML appends a view to a parent DOM element. It takes a View interface and
// a parent DOM selector. parentSelector works identically to JavaScript's document.querySelector(selector)
// call. After the view's element is added to the DOM, AppendToParentHTML calls view.OnLoad if it is defined.
func AppendToParentHTML(view View, parentSelector string) error {
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
func ReplaceParentHTML(view View, parentSelector string) error {
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
func Update(view View) error {
	html := view.RenderHTML()
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
func Remove(view View) error {
	viewEl, err := getElementByViewId(view.GetId())
	if err != nil {
		return err
	}
	viewEl.ParentElement().RemoveChild(viewEl)
	return nil
}

func getElementByViewId(viewId string) (dom.Element, error) {
	// Use a query selector to find the element in the DOM
	selector := fmt.Sprintf("[data-humble-view-id='%s']", viewId)
	el := document.QuerySelector(selector)
	if el == nil {
		return nil, humble.NewViewElementNotFoundError(viewId)
	}
	return el, nil
}

// createViewElement creates a DOM element from HTML and a outer container tag.
// Takes innerHTML and outerTag and crafts a valid dom.Element. Returns the resultant
// dom.Element or an error.
func createViewElement(view View) (dom.Element, error) {
	// Check our outer container tag is valid
	if err := checkOuterTag(view.OuterTag()); err != nil {
		return nil, err
	}
	// Get our view HTML
	viewHTML := view.RenderHTML()
	// Check if view element exists in global map, otherwise create it
	var el dom.Element
	if existingEl, err := getElementByViewId(view.GetId()); err != nil {
		if _, notFound := err.(humble.ViewElementNotFoundError); notFound {
			// The view was not found in the DOM. We need to create it
			el = document.CreateElement(view.OuterTag())
		} else {
			// For any other type of error, return it.
			return nil, err
		}
	} else {
		el = existingEl
	}
	el.SetInnerHTML(viewHTML)
	// We set attribute data-humble-view-id on outer container so we can get it from the
	// DOM later on with a QuerySelector
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
		err := onLoader.OnLoad() //gopherjs:blocking
		if err != nil {
			return fmt.Errorf("Error in %T.OnLoad: %s", v, err.Error())
		}
	}
	return nil
}

// QuerySelector takes a selector string and returns the first matching element within the given view as a dom.Element.
// Will return an error if no matching element is found.
func QuerySelector(view View, selector string) (dom.Element, error) {
	fullSelector := fmt.Sprintf("[data-humble-view-id='%s'] %s", view.GetId(), selector)
	targetEls := document.QuerySelector(fullSelector)
	if targetEls == nil {
		return nil, fmt.Errorf("Could not find element with selector: `%s` inside of element for %T. Full selector was: `%s`", selector, view, fullSelector)
	}
	return targetEls, nil
}

// QuerySelectorAll takes a selector string and returns all matching elements within the given view as a []dom.Element.
// Will return an error if no matching elements are found.
func QuerySelectorAll(view View, selector string) ([]dom.Element, error) {
	fullSelector := fmt.Sprintf("[data-humble-view-id='%s'] %s", view.GetId(), selector)
	targetEls := document.QuerySelectorAll(fullSelector)
	if len(targetEls) == 0 {
		return nil, fmt.Errorf("Could not find any elements with selector: `%s` inside of element for %T. Full selector was: `%s`", selector, view, fullSelector)
	}
	return targetEls, nil
}

// AddListener adds an event listener to the element inside the view's element identified by childSelector.
// It expects a Listener as an argument, which will be called when the event is triggered. When this
// calls el.AddEventListener it sets the useCapture option to false. AddListener cannot be used to
// listen for events on elements outside of the view's element.
// Example:
// humble.AddListener(todoView, "button.destroy", "click", func(dom.Event) {
//		if err := views.Remove(todoView); err != nil {
// 		...
// 	}
// })
func AddListener(view View, childSelector string, eventName string, listener Listener) error {
	targetEls, err := QuerySelectorAll(view, childSelector)
	if err != nil {
		return err
	}
	for _, el := range targetEls {
		el.AddEventListener(eventName, false, nonBlockingListener(listener))
	}
	return nil
}

// Show will show the view's element by adding the style "display: block;"
// If the view's element is already shown, it will do nothing.
func Show(view View) error {
	el, err := getElementByViewId(view.GetId())
	if err != nil {
		return err
	}
	el.SetAttribute("style", "display: block;")
	return nil
}

// Hide will hide the view's element by adding the style "display: none;"
// If the view's element is already hidden, it will do nothing.
func Hide(view View) error {
	el, err := getElementByViewId(view.GetId())
	if err != nil {
		return err
	}
	el.SetAttribute("style", "display: none;")
	return nil
}

// nonBlockingListener takes care of wrapping our event listener functions with a goroutine to make these usually
// blocking calls non-blocking, as required by GopherJS
func nonBlockingListener(listener Listener) Listener {
	return func(ev dom.Event) {
		go func() {
			listener(ev) //gopherjs:blocking
		}()
	}
}
