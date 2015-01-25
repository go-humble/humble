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
		return fmt.Errorf("Something went wrong with Get request to %s. %s", urlRoot, err.Error())
	}
	err = json.Unmarshal([]byte(req.Response.String()), models)
	if err != nil {
		return err
	}

	return nil
}
