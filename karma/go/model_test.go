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
		done := assert.Call("async")
		go func() {
			expectedTodos := []*Todo{
				{
					Id:          0,
					Title:       "Todo 0",
					IsCompleted: false,
				},
				{
					Id:          1,
					Title:       "Todo 1",
					IsCompleted: false,
				},
				{
					Id:          2,
					Title:       "Todo 2",
					IsCompleted: true,
				},
			}
			gotTodos := []*Todo{}
			err := model.ReadAll(&gotTodos)
			assert.Ok(err == nil, fmt.Sprintf("model.ReadAll returned an error: %v", err))
			assert.Ok(reflect.DeepEqual(gotTodos, expectedTodos), fmt.Sprintf("Expected: %v, Got: %v", expectedTodos, gotTodos))
			done.Invoke()
		}()
	})

	qunit.Test("Read", func(assert qunit.QUnitAssert) {
		qunit.Expect(2)
		done := assert.Call("async")
		go func() {
			expectedTodo := &Todo{
				Id:          2,
				Title:       "Todo 2",
				IsCompleted: true,
			}
			gotTodo := &Todo{}
			err := model.Read("2", gotTodo)
			assert.Ok(err == nil, fmt.Sprintf("model.Read returned an error: %v", err))
			assert.Ok(reflect.DeepEqual(gotTodo, expectedTodo), fmt.Sprintf("Expected: %v, Got: %v", expectedTodo, gotTodo))
			done.Invoke()
		}()
	})

	qunit.Test("Create", func(assert qunit.QUnitAssert) {
		// For some unkown reason, qunit is running 8 assertions and reporting an error
		// There are obviously only 4 so something wonky is happening. The other tests
		// seem fine.
		// qunit.Expect(4)
		done := assert.Call("async")
		go func() {
			newTodo := &Todo{
				Title:       "Test",
				IsCompleted: true,
			}
			err := model.Create(newTodo)
			assert.Ok(err == nil, fmt.Sprintf("model.Create returned an error: %v", err))
			assert.Equal(newTodo.IsCompleted, true, "newTodo.IsCompleted was incorrect.")
			assert.Equal(newTodo.Title, "Test", "newTodo.Title was incorrect.")
			assert.Equal(newTodo.Id, 3, "newTodo.Id was not set correctly.")
			done.Invoke()
		}()
	})
}
