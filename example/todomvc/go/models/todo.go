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

func (t *Todo) UrlRoot() string {
	return "http://localhost:3000/todos"
}

func (t *Todo) CheckedStr() string {
	if t.IsCompleted {
		return "checked"
	}
	return ""
}
