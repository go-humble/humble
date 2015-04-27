package models

import (
	"strconv"
)

type Todo struct {
	Id          int
	Title       string
	IsCompleted bool
}

func (t *Todo) GetId() string {
	return strconv.Itoa(t.Id)
}

func (t *Todo) RootURL() string {
	return "http://localhost:3000/todos"
}

func (t Todo) CheckedStr() string {
	if t.IsCompleted {
		return "checked"
	}
	return ""
}

func (t Todo) CompletedStr() string {
	if t.IsCompleted {
		return "completed"
	}
	return ""
}

type Todos []*Todo

func (todos Todos) NumCompleted() int {
	count := 0
	for _, todo := range todos {
		if todo.IsCompleted {
			count += 1
		}
	}
	return count
}

func (todos Todos) NumRemaining() int {
	count := 0
	for _, todo := range todos {
		if !todo.IsCompleted {
			count += 1
		}
	}
	return count
}
