package main

import (
	"fmt"
	"github.com/gophergala/humble"
	"honnef.co/go/js/console"
	"strconv"
)

type TodoModel struct {
	Id          int
	Title       string
	IsCompleted bool
}

func (tm *TodoModel) UrlRoot() string {
	return "http://localhost:3000/todos"
}

func (tm *TodoModel) GetId() string {
	return strconv.Itoa(tm.Id)
}

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

	//Model
	mCollection := []*TodoModel{}
	err := humble.Models.GetAll(&mCollection)
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range mCollection {
		fmt.Println("m: ", *m)
	}

	//View
	v := &CurrentPageView{}

	r := humble.NewRouter()
	r.HandleFunc("/", func(params map[string]string) {
		v.Name = "Home"
		fmt.Println(v.GetId())
		if err := humble.Views.ReplaceParentHTML(v, "#current-page"); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/about", func(params map[string]string) {
		v.Name = "About"
		fmt.Println(v.GetId())
		if err := humble.Views.Update(v); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/faq", func(params map[string]string) {
		v.Name = "FAQ"
		fmt.Println(v.GetId())
		if err := humble.Views.ReplaceParentHTML(v, "#current-page"); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/remove", func(params map[string]string) {
		fmt.Println(v.GetId())
		if err := humble.Views.Remove(v); err != nil {
			panic(err)
		}
	})

	r.Start()

}
