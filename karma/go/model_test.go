package main

import (
	"fmt"
	"github.com/albrow/qunit"
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
		qunit.Expect(4)
		done := assert.Call("async")
		go func() {
			newTodo := &Todo{
				Title:       "Test",
				IsCompleted: true,
			}
			err := model.Create(newTodo)
			assert.Ok(err == nil, fmt.Sprintf("model.Create returned an error: %v", err))
			assert.Equal(newTodo.Id, 3, "newTodo.Id was not set correctly.")
			assert.Equal(newTodo.Title, "Test", "newTodo.Title was incorrect.")
			assert.Equal(newTodo.IsCompleted, true, "newTodo.IsCompleted was incorrect.")
			done.Invoke()
		}()
	})

	qunit.Test("Update", func(assert qunit.QUnitAssert) {
		qunit.Expect(4)
		done := assert.Call("async")
		go func() {
			updatedTodo := &Todo{
				Id:          1,
				Title:       "Updated Title",
				IsCompleted: true,
			}
			err := model.Update(updatedTodo)
			assert.Ok(err == nil, fmt.Sprintf("model.Update returned an error: %v", err))
			assert.Equal(updatedTodo.Id, 1, "updatedTodo.Id was incorrect.")
			assert.Equal(updatedTodo.Title, "Updated Title", "updatedTodo.Title was incorrect.")
			assert.Equal(updatedTodo.IsCompleted, true, "updatedTodo.IsCompleted was incorrect.")
			done.Invoke()
		}()
	})

	qunit.Test("Delete", func(assert qunit.QUnitAssert) {
		qunit.Expect(1)
		done := assert.Call("async")
		go func() {
			deletedTodo := &Todo{
				Id: 1,
			}
			err := model.Delete(deletedTodo)
			assert.Ok(err == nil, fmt.Sprintf("model.Update returned an error: %v", err))
			done.Invoke()
		}()
	})
}
