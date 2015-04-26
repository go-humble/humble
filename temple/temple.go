// Package temple is a lightweight wrapper around standard html/template
// package. It includes utility functions for parsing inline templates
// and rendering templates in the DOM. For information about standard
// go templating see http://golang.org/pkg/text/template/ and
// http://golang.org/pkg/html/template/
package temple

import (
	"bytes"
	"honnef.co/go/js/dom"
	"html/template"
)

var (
	document = dom.GetWindow().Document()
	// Templates is a map of template names to temple.Template objects
	Templates = map[string]*Template{}
)

// Template is a lightweight wrapper around standard go templates from
// the html/template package.
type Template struct {
	GoTemplate template.Template
}

// Load loads inline templates from the current document. Templates must have
// the type "text/template", and the id will be used as the template name. So
// for example, the inline template identified by the tag
//   <script type="text/template" id="todo">
// will be loaded and parsed as a template with the name "todo".
// Load will return an error if there is a syntax error in any of the templates.
func Load() error {
	inlineTemplates := document.QuerySelectorAll(`script[type="text/template"]`)
	for _, inlineTemplate := range inlineTemplates {
		src := inlineTemplate.InnerHTML()
		name := inlineTemplate.ID()
		tmpl, err := template.New(name).Parse(src)
		if err != nil {
			return err
		}
		Templates[name] = &Template{
			GoTemplate: *tmpl,
		}
	}
	return nil
}

// Render renders the Template with the given data and sets el.innerHTML.
// It will return an error if there was a problem rendering the template
// with the given data.
func (tmpl *Template) Render(el dom.Element, data interface{}) error {
	buf := bytes.NewBuffer([]byte{})
	if err := tmpl.GoTemplate.Execute(buf, data); err != nil {
		return err
	}
	el.SetInnerHTML(buf.String())
	return nil
}
