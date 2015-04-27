package view

import (
	"honnef.co/go/js/dom"
	"strings"
)

var (
	document = dom.GetWindow().Document()
)

type View interface {
	Render() error
	Element() dom.Element
	SetElement(dom.Element)
}

func Append(parent View, child View) {
	parent.Element().AppendChild(child.Element())
}

func AppendToEl(parent dom.Element, child View) {
	parent.AppendChild(child.Element())
}

func Replace(old View, new View) {
	old.Element().ParentElement().ReplaceChild(old.Element(), new.Element())
}

func ReplaceEl(old dom.Element, new View) {
	old.ParentElement().ReplaceChild(old, new.Element())
}

func Remove(v View) {
	v.Element().ParentElement().RemoveChild(v.Element())
}

func Hide(v View) {
	oldStyles := v.Element().GetAttribute("style")
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
	v.Element().SetAttribute("style", newStyles)
}

func Show(v View) {
	oldStyles := v.Element().GetAttribute("style")
	// Try removing the with a semicolon version first.
	// If there is not a semicolon, this will have no effect.
	newStyles := strings.Replace(oldStyles, "display:none;", "", 1)
	// Then try removing the without a semicolon version.
	// If there was a semicolon, this will have no effect.
	newStyles = strings.Replace(newStyles, "display:none", "", 1)
	v.Element().SetAttribute("style", newStyles)
}
