package main

import (
	"encoding/json"
	"fmt"
	"github.com/albrow/forms"
	"github.com/albrow/negroni-json-recovery"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/martini-contrib/cors"
	"github.com/unrolled/render"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

type Todo struct {
	Id          int
	Title       string
	IsCompleted bool
}

type todosIndex map[int]*Todo
type todosList []*Todo

var (
	// todos stores all the todos as a map of id to *Todo
	todos = todosIndex{}
	// todosMutex protects access to the todos map
	todosMutex = sync.Mutex{}
	// todosCounter is incremented every time a new todo is created
	// it is used to set todo ids.
	todosCounter = 0
	// r is used to render responses
	r = render.New(render.Options{
		IndentJSON: true,
	})
)

const (
	statusUnprocessableEntity = 422
)

func main() {
	createInitialTodos()

	// Routes
	router := mux.NewRouter()
	router.HandleFunc("/todos", todosController.Index).Methods("GET")
	router.HandleFunc("/todos", todosController.Create).Methods("POST")
	router.HandleFunc("/todos/{id}", todosController.Read).Methods("GET")
	router.HandleFunc("/todos/{id}", todosController.Update).Methods("PUT")
	router.HandleFunc("/todos/{id}", todosController.Delete).Methods("DELETE")

	// Other middleware
	n := negroni.New(negroni.NewLogger())
	n.UseHandler(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	n.Use(recovery.JSONRecovery(true))
	recovery.StackDepth = 3
	recovery.IndentJSON = true
	recovery.Formatter = func(errMsg string, stack []byte, file string, line int, fullMessages bool) interface{} {
		return map[string]string{
			"error": errMsg,
		}
	}

	// Router must always come last
	n.UseHandler(router)

	// Start the server
	n.Run(":3000")
}

func createInitialTodos() {
	createTodo("Write a frontend framework in Go")
	createTodo("???")
	createTodo("Profit!")
}

func createTodo(title string) *Todo {
	todosMutex.Lock()
	defer todosMutex.Unlock()
	id := todosCounter
	todosCounter++
	todo := &Todo{
		Id:    id,
		Title: title,
	}
	todos[id] = todo
	return todo
}

// Todos Controller and its methods
type todosControllerType struct{}

var todosController = todosControllerType{}

func (todosControllerType) Index(w http.ResponseWriter, req *http.Request) {
	r.JSON(w, http.StatusOK, todos)
}

func (todosControllerType) Create(w http.ResponseWriter, req *http.Request) {
	// Parse data and do validations
	todoData, err := forms.Parse(req)
	if err != nil {
		panic(err)
	}
	val := todoData.Validator()
	val.Require("Title")
	if val.HasErrors() {
		r.JSON(w, statusUnprocessableEntity, val.ErrorMap())
		return
	}

	// Create the todo and render response
	todo := createTodo(todoData.Get("Title"))
	r.JSON(w, http.StatusOK, todo)
}

func (todosControllerType) Read(w http.ResponseWriter, req *http.Request) {
	urlParams := mux.Vars(req)
	idString := urlParams["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}

	r.JSON(w, http.StatusOK, todos[id])
}

func (todosControllerType) Update(w http.ResponseWriter, req *http.Request) {
	// Get the existing todo from the map or render an error
	// if it wasn't found
	urlParams := mux.Vars(req)
	idString := urlParams["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}
	todo, found := todos[id]
	if !found {
		msg := fmt.Sprintf("Could not find todo with id = %d", id)
		r.JSON(w, http.StatusNotFound, map[string]string{
			"error": msg,
		})
		return
	}

	// Update the todo with the data in the request
	todoData, err := forms.Parse(req)
	if err != nil {
		panic(err)
	}
	todosMutex.Lock()
	if todoData.KeyExists("Title") {
		todo.Title = todoData.Get("Title")
	}
	if todoData.KeyExists("IsCompleted") {
		todo.IsCompleted = todoData.GetBool("IsCompleted")
	}
	todosMutex.Unlock()

	// Render response
	r.JSON(w, http.StatusOK, todo)
}

func (todosControllerType) Delete(w http.ResponseWriter, req *http.Request) {
	// Get the id from the url parameters
	urlParams := mux.Vars(req)
	idString := urlParams["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}

	// Delete the todo and render a response
	todosMutex.Lock()
	delete(todos, id)
	todosMutex.Unlock()
	r.JSON(w, http.StatusOK, struct{}{})
}

// Make todosIndex satisfy the json.Marshaller interface
// It will return a json array of todos sorted by id
func (t todosIndex) MarshalJSON() ([]byte, error) {
	todosList := todosList{}
	for _, todo := range t {
		todosList = append(todosList, todo)
	}
	sort.Sort(todosList)
	return json.Marshal(todosList)
}

// Make todoList satisfy sort.Interface
func (tl todosList) Len() int {
	return len(tl)
}

func (tl todosList) Less(i, j int) bool {
	return tl[i].Id < tl[j].Id
}

func (tl todosList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}
