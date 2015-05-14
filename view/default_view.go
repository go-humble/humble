package view

import (
	"honnef.co/go/js/dom"
)

type DefaultView struct {
	el dom.Element
}

func (v *DefaultView) Element() dom.Element {
	if v.el == nil {
		v.el = document.CreateElement("div")
	}
	return v.el
}

func (v *DefaultView) SetElement(el dom.Element) {
	v.el = el
}
