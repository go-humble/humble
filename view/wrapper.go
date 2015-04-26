package view

import (
	"honnef.co/go/js/dom"
	"strings"
)

var (
	document = dom.GetWindow().Document()
)

type Wrapper struct {
	El dom.Element
}

func NewWrapper(tagName string, attributes map[string]string) *Wrapper {
	el := document.CreateElement(tagName)
	for name, value := range attributes {
		el.SetAttribute(name, value)
	}
	return &Wrapper{
		El: el,
	}
}

func (w *Wrapper) AppendTo(el dom.Element) {
	el.AppendChild(w.El)
}

func (w *Wrapper) Replace(el dom.Element) {
	w.El.ParentElement().ReplaceChild(w.El, el)
}

func (w *Wrapper) Remove() {
	w.El.ParentElement().RemoveChild(w.El)
}

func (w *Wrapper) Hide() {
	oldStyles := w.El.GetAttribute("style")
	newStyles := ""
	switch {
	case oldStyles == "":
		// There was no style attribute. We can safely set
		// the style attribute directly.
		newStyles = "display:none"
	case strings.Contains(oldStyles, "display:none"):
		// The element is already hidden. We should do
		// nothing.
		return
	case oldStyles[len(oldStyles)] == ';':
		// There was a style attribute and it ended in a semicolon,
		// We can safely append the new styles to the old.
		newStyles += "display:none;"
	default:
		// There was a style attribute and it didn't end in a semicolon,
		// in this case we should add our own semicolon.
		newStyles += ";display:none;"
	}
	w.El.SetAttribute("style", newStyles)
}

func (w *Wrapper) Show() {
	oldStyles := w.El.GetAttribute("style")
	// Try removing the with a semicolon version first.
	// If there is not a semicolon, this will have no effect.
	newStyles := strings.Replace(oldStyles, "display:none;", "", 1)
	// Then try removing the without a semicolon version.
	// If there was a semicolon, this will have no effect.
	newStyles = strings.Replace(newStyles, "display:none", "", 1)
	w.El.SetAttribute("style", newStyles)
}
