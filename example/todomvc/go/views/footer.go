package views

import (
	"fmt"
	"github.com/gophergala/humble"
)

type Footer struct {
	humble.Identifier
	Remaining int
	TodoViews *[]*Todo
}

func (f *Footer) RenderHTML() string {
	return fmt.Sprintf(`
			<span id="todo-count">
				<strong>%d</strong> items left
			</span>
			<ul id="filters">
				<li>
					<a class="selected" href="#/">All</a>
				</li>
				<li>
					<a href="#/active">Active</a>
				</li>
				<li>
					<a href="#/completed">Completed</a>
				</li>
			</ul>`, f.countRemaining())
}

func (f *Footer) OuterTag() string {
	return "div"
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
