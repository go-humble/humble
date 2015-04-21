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
	qunit.Test("ReadAll", func(assert qunit.QUnitAssert) {
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
		done := assert.Call("async")
		go func() {
			gotTodos := []*Todo{}
			err := model.ReadAll(&gotTodos)
			assert.Ok(err == nil, fmt.Sprintf("model.ReadAll returned an error: %v", err))
			assert.Ok(reflect.DeepEqual(gotTodos, expectedTodos), fmt.Sprintf("Expected: %v, Got: %v", expectedTodos, gotTodos))
			done.Invoke()
		}()
	})

	qunit.Test("Read", func(assert qunit.QUnitAssert) {
		qunit.Expect(2)
		expectedTodo := &Todo{
			Id:          0,
			Title:       "Write a frontend framework in Go",
			IsCompleted: false,
		}
		done := assert.Call("async")
		go func() {
			gotTodo := &Todo{}
			err := model.Read("0", gotTodo)
			assert.Ok(err == nil, fmt.Sprintf("model.Read returned an error: %v", err))
			assert.Ok(reflect.DeepEqual(gotTodo, expectedTodo), fmt.Sprintf("Expected: %v, Got: %v", expectedTodo, gotTodo))
			done.Invoke()
		}()
	})
}
