package main

import (
	"github.com/JohannWeging/jasmine"
	"github.com/soroushjp/humble/model"
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
	jasmine.Describe("ReadAll", func() {
		jasmine.It("gets all the existing todos", func() {
			expectedModels := []*Todo{
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
			go func() {
				gotModels := []*Todo{}
				err := model.ReadAll(&gotModels)
				jasmine.Expect(err).ToEqual(nil)
				jasmine.Expect(gotModels).ToEqual(expectedModels)
			}()
		})
	})

}
