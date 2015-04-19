package main

import (
	"fmt"
	"github.com/rusco/qunit"
	"github.com/soroushjp/humble/model"
	"reflect"
	"strconv"
)

type Todo struct {
	Id          int
	Title       string
	IsCompleted bool
}

func (t Todo) GetId() string {
	return strconv.Itoa(t.Id)
}

func (t Todo) RootURL() string {
	return "http://localhost:3000/todos"
}

func main() {
	qunit.AsyncTest("ReadAll", func() interface{} {
		qunit.Expect(2)
		expectedTodos := []*Todo{
			{
				Id:          0,
				Title:       "Write a frontend framework in Go",
				IsCompleted: false,
			},
			{
				Id:          1,
				Title:       "???",
				IsCompleted: false,
			},
			{
				Id:          2,
				Title:       "Profit!",
				IsCompleted: false,
			},
		}
		gotTodos := []*Todo{}
		err := model.ReadAll(&gotTodos)
		qunit.Ok(err == nil, fmt.Sprintf("model.ReadAll returned an error: %v", err))
		qunit.Ok(reflect.DeepEqual(gotTodos, expectedTodos), fmt.Sprintf("Expected: %v, Got: %v", expectedTodos, gotTodos))
		qunit.Start()
		return nil
	})
}
