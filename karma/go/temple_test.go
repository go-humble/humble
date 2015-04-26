package main

import (
	"fmt"
	"github.com/albrow/qunit"
	"github.com/soroushjp/humble/temple"
	"honnef.co/go/js/dom"
)

var (
	document = dom.GetWindow().Document()
	body     = document.QuerySelector("body")
)

func main() {
	qunit.Test("Render", func(assert qunit.QUnitAssert) {
		qunit.Expect(5)
		createInlineTemplate()
		assert.Ok(len(document.QuerySelectorAll("script#greet")) == 1, "Inline template was not inserted into the DOM correctly.")
		err := temple.Load()
		assert.Ok(err == nil, fmt.Sprintf("Error loading inline templates: %v", err))
		greetTmpl, found := temple.Templates["greet"]
		assert.Ok(found, "greet template was not found")
		greetEl := document.CreateElement("div")
		greetEl.SetID("greeting")
		body.AppendChild(greetEl)
		err = greetTmpl.Render(greetEl, map[string]string{"Name": "Foo"})
		assert.Ok(err == nil, fmt.Sprintf("Error in greetTmpl.Render: %v", err))
		greetInner := greetEl.InnerHTML()
		assert.Equal(greetInner, "<h1>Hello, Foo</h1>", "Rendered template had the wrong inner html.")
	})
}

// createInlineTemplate adds a script tag to the body which contains a
// template called "greet"
func createInlineTemplate() {
	inlineTmpl := document.CreateElement("script")
	inlineTmpl.SetAttribute("type", "text/template")
	inlineTmpl.SetID("greet")
	inlineTmpl.SetInnerHTML("<h1>Hello, {{ .Name }}</h1>")
	body.AppendChild(inlineTmpl)
}
