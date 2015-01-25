package main

import (
	"fmt"
	"github.com/gophergala/humble"
	"honnef.co/go/js/console"
)

type CurrentPageView struct {
	Name string
	humble.Identifier
}

func (cp *CurrentPageView) GetHTML() string {
	return fmt.Sprintf("<p>Current page is %s</p>", cp.Name)
}

func (t *CurrentPageView) OuterTag() string {
	return "div"
}

func main() {
	console.Log("Starting...")

	v := &CurrentPageView{}

	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		v.Name = "Home"
		if err := humble.Views.SetInnerHTML(v, "#current-page"); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/about", func(params map[string]string) {
		v.Name = "About"
		if err := humble.Views.Update(v); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/faq", func(params map[string]string) {
		v.Name = "FAQ"
		if err := humble.Views.SetInnerHTML(v, "#current-page"); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/remove", func(params map[string]string) {
		if err := humble.Views.Remove(v); err != nil {
			panic(err)
		}
	})

	r.Start()

}
