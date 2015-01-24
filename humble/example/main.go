package main

import (
	"fmt"
	"github.com/gophergala/go_by_the_fireplace/humble"
	"honnef.co/go/js/console"
)

type CurrentPageView struct {
	humble.Identifier
}

type CurrentPage struct {
	Name string
	humble.Identifier
}

func (t *CurrentPageView) GetHTML(m humble.Model) string {
	cp := m.(*CurrentPage)
	fmt.Println(cp)
	return fmt.Sprintf("<p>Current page is %s</p>", cp.Name)
}

func (t *CurrentPageView) OuterTag() string {
	return "div"
}

func main() {
	console.Log("Starting...")

	v := &CurrentPageView{}
	homePage := &CurrentPage{Name: "Home"}
	aboutPage := &CurrentPage{Name: "About"}
	faqPage := &CurrentPage{Name: "FAQ"}

	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		console.Log(v.Id())
		console.Log("At home page")
		if err := humble.Views.SetOnlyChild(v, homePage, "#current-page"); err != nil {
			console.Error(err)
		}
	})
	r.HandleFunc("/about", func(params map[string]string) {
		console.Log(v.Id())
		console.Log("At about page")
		if err := humble.Views.SetOnlyChild(v, aboutPage, "#current-page"); err != nil {
			console.Error(err)
		}
	})
	r.HandleFunc("/faq", func(params map[string]string) {
		console.Log(v.Id())
		console.Log("At faq page")
		if err := humble.Views.SetOnlyChild(v, faqPage, "#current-page"); err != nil {
			console.Error(err)
		}
	})
	r.HandleFunc("/remove", func(params map[string]string) {
		console.Log(v.Id())
		if ok := humble.Views.Remove(v); !ok {
			console.Error("Not Ok!")
		}
	})

	r.Start()

}
