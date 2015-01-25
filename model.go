package humble

import (
	"encoding/json"
	"fmt"
	"honnef.co/go/js/xhr"
	"reflect"
)

type Model interface {
	GetId() string
	UrlRoot() string
}

type modelsType struct{}

var Models = modelsType{}

// GetAll expects a pointer to a slice of poitners to some concrete type
// which implements Model (e.g., *[]*Todo). GetAll will send a GET request to
// a RESTful server and scan the results into models. It expects a json array
// of json objects from the server, where each object represents a single Model
// of some concrete type. It will use the UrlRoot() method of the models to
// figure out which url to send the GET request to.
func (*modelsType) GetAll(models interface{}) error {
	// We expect some pointer to a slice of models. Like *[]Todo
	// First Elem() givew us []Todo
	// Second Elem() gives us Todo
	modelsType := reflect.TypeOf(models).Elem()
	// TODO: type checking!
	modelTypePtr := modelsType.Elem()
	modelType := modelTypePtr.Elem()
	// reflect.New returns *Todo
	newModelVal := reflect.New(modelType)
	newModelInterface := newModelVal.Interface()
	newModel := newModelInterface.(Model)
	// TODO: check for a failed type assertion!
	urlRoot := newModel.UrlRoot()

	req := xhr.NewRequest("GET", urlRoot)
	req.Timeout = 1000 // one second, in milliseconds
	req.ResponseType = "text"
	err := req.Send(nil)
	if err != nil {
		return fmt.Errorf("Something went wrong with GET request to %s. %s", urlRoot, err.Error())
	}
	err = json.Unmarshal([]byte(req.Response.String()), models)
	if err != nil {
		return err
	}

	return nil
}

// Delete expects a pointer some concrete type which implements Model (e.g., *[]*Todo).
// DELETE will send a DELETE request to a RESTful server. It expects an empty json
// object from the server if the request was successful, and will not attempt to do anything
// with the response. It will use the UrlRoot() and GetId() methods of the models to determine
// which url to send the DELETE request to. Typically, the full url will look something
// like "http://hostname.com/todos/123"
func (*modelsType) Delete(model Model) error {
	fullURL := model.UrlRoot() + "/" + model.GetId()
	req := xhr.NewRequest("DELETE", fullURL)
	req.Timeout = 1000 // one second, in milliseconds
	req.ResponseType = "text"
	err := req.Send(nil)
	if err != nil {
		return fmt.Errorf("Something went wrong with DELETE request to %s. %s", fullURL, err.Error())
	}
	return nil
}
