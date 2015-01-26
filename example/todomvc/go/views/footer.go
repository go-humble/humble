package views

import (
	"fmt"
	"github.com/gophergala/humble"
	"github.com/gophergala/humble/view"
	"honnef.co/go/js/dom"
	"strings"
)

type Footer struct {
	humble.Identifier
	TodoViews *[]*Todo
}

func (f *Footer) RenderHTML() string {
	return fmt.Sprintf(`
			<span id="todo-count">
				<strong>%d</strong> items left
			</span>
			<ul id="filters">
				<li>
					<a href="#/">All</a>
				</li>
				<li>
					<a href="#/active">Active</a>
				</li>
				<li>
					<a href="#/completed">Completed</a>
				</li>
			</ul>
			<button id="clear-completed">Clear completed (%d)</button>`, f.countRemaining(), f.countCompleted())
}

func (f *Footer) OuterTag() string {
	return "div"
}

func (f *Footer) OnLoad() error {
	// Set some things in the DOM which may have changed
	// after an update.
	if err := f.setSelected(); err != nil {
		return err
	}
	if f.countCompleted() == 0 {
		f.hideClearCompleted()
	} else {
		f.showClearCompleted()
	}

	// Add listeners
	if err := view.AddListener(f, "button#clear-completed", "click", f.clearCompleted); err != nil {
		return err
	}

	return nil
}

func (f *Footer) countRemaining() int {
	count := 0
	if f.TodoViews == nil {
		return 0
	}
	for _, todoView := range *f.TodoViews {
		if !todoView.Model.IsCompleted {
			count++
		}
	}
	return count
}

func (f *Footer) countCompleted() int {
	count := 0
	if f.TodoViews == nil {
		return 0
	}
	for _, todoView := range *f.TodoViews {
		if todoView.Model.IsCompleted {
			count++
		}
	}
	return count
}

func (f *Footer) setSelected() error {
	hash := dom.GetWindow().Location().Hash
	links, err := view.QuerySelectorAll(f, "#filters li a")
	if err != nil {
		return err
	}
	for _, linkEl := range links {
		if linkEl.GetAttribute("href") == hash {
			linkEl.SetAttribute("class", linkEl.GetAttribute("class")+" selected")
		} else {
			oldClass := linkEl.GetAttribute("class")
			newClass := strings.Replace(oldClass, "selected", "", 1)
			linkEl.SetAttribute("class", newClass)
		}
	}
	return nil
}

func (f *Footer) clearCompleted(dom.Event) {
	f.hideClearCompleted()
	if f.TodoViews == nil {
		return
	}
	todosToRemove := []*Todo{}
	for _, todoView := range *(f.TodoViews) {
		if todoView.Model.IsCompleted {
			todosToRemove = append(todosToRemove, todoView)
		}
	}
	for _, todoView := range todosToRemove {
		todoView.remove()
	}
}

func (f *Footer) showClearCompleted() {
	clrCompletedEl, err := view.QuerySelector(f, "button#clear-completed")
	if err != nil {
		panic(err)
	}
	clrCompletedEl.SetAttribute("style", "display: block;")
}

func (f *Footer) hideClearCompleted() {
	clrCompletedEl, err := view.QuerySelector(f, "button#clear-completed")
	if err != nil {
		panic(err)
	}
	clrCompletedEl.SetAttribute("style", "display: none;")
}
