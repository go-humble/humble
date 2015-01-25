package model

import (
	"encoding/json"
	"fmt"
	"honnef.co/go/js/xhr"
	"net/url"
	"reflect"
)

type Model interface {
	GetId() string
	RootURL() string
}

// ReadAll expects a pointer to a slice of poitners to some concrete type
// which implements Model (e.g., *[]*Todo). GetAll will send a GET request to
// a RESTful server and scan the results into models. It expects a json array
// of json objects from the server, where each object represents a single Model
// of some concrete type. It will use the RootURL() method of the models to
// figure out which url to send the GET request to.
func ReadAll(models interface{}) error {
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
	urlRoot := newModel.RootURL()

	req := xhr.NewRequest("GET", urlRoot)
	req.Timeout = 1000 // one second, in milliseconds
	req.ResponseType = "text"
	err := req.Send(nil)
	if err != nil {
		return fmt.Errorf("Something went wrong with GET request to %s. %s", urlRoot, err.Error())
	}
	err = json.Unmarshal([]byte(req.Response.String()), models)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal response into object, with Error: %s.\nResponse was: %s", err, req.Response.String())
	}

	return nil
}

// Delete expects a pointer some concrete type which implements Model (e.g., *Todo).
// It will send a DELETE request to a RESTful server. It expects an empty json
// object from the server if the request was successful, and will not attempt to do anything
// with the response. It will use the RootURL() and GetId() methods of the model to determine
// which url to send the DELETE request to. Typically, the full url will look something
// like "http://hostname.com/todos/123"
func Delete(model Model) error {
	fullURL := model.RootURL() + "/" + model.GetId()
	req := xhr.NewRequest("DELETE", fullURL)
	req.Timeout = 1000 // one second, in milliseconds
	req.ResponseType = "text"
	err := req.Send(nil)
	if err != nil {
		return fmt.Errorf("Something went wrong with DELETE request to %s. %s", fullURL, err.Error())
	}
	return nil
}

// Create expects a pointer some concrete type which implements Model (e.g., *Todo).
// It will send a POST request to the RESTful server. It expects a JSON containing the
// created object from the server if the request was successful, and will set the fields of
// model with the data in the response object. It will use the RootURL() method of
// the model to determine which url to send the POST request to.
func Create(model Model) error {
	fullURL := model.RootURL()
	bodyString := ""
	// TODO: Do stronger type checking to prevent errors
	//Use reflect to identify fields of our model and convert to URL-encoded formdata string
	modelValPtr := reflect.ValueOf(model)
	modelVal := modelValPtr.Elem()
	modelValType := modelVal.Type()
	for i := 0; i < modelValType.NumField(); i++ {
		field := modelVal.Field(i)
		fieldName := modelValType.Field(i).Name
		value := fmt.Sprint(field.Interface())
		bodyString += fieldName + "=" + url.QueryEscape(value)
		if i != modelValType.NumField()-1 {
			bodyString += "&"
		}
	}
	// Create our x-www-form-urlencoded POST request
	req := xhr.NewRequest("POST", fullURL)
	req.Timeout = 1000 //one second, in milliseconds
	req.ResponseType = "text"
	req.SetRequestHeader("Content-Type", "application/x-www-form-urlencoded")
	//Send our request
	err := req.Send(bodyString)
	if err != nil {
		return err
	}
	// Unmarshal our server response object into our model
	err = json.Unmarshal([]byte(req.Response.String()), model)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal response into object, with Error: %s.\nResponse was: %s", err, req.Response.String())
	}
	return nil
}

// Update expects a pointer some concrete type which implements Model (e.g., *Todo), with a model.Id
// that matches a stored object on the server. It will send a PUT request to the RESTful server.
// It expects a JSON containing the updated object from the server if the request was successful,
// and will set the fields of model with the data in the response object.
// It will use the RootURL() method of the model to determine which url to send the PUT request to.
func Update(model Model) error {
	//Set our request URL to be root URL/Id, eg. example.com/api/todos/4
	fullURL := model.RootURL() + "/" + model.GetId()
	bodyString := ""
	// TODO: Do stronger type checking to prevent errors
	//Use reflect to identify fields of our model and convert to URL-encoded formdata string
	modelValPtr := reflect.ValueOf(model)
	modelVal := modelValPtr.Elem()
	modelValType := modelVal.Type()
	for i := 0; i < modelValType.NumField(); i++ {
		field := modelVal.Field(i)
		fieldName := modelValType.Field(i).Name
		value := fmt.Sprint(field.Interface())
		bodyString += fieldName + "=" + url.QueryEscape(value)
		if i != modelValType.NumField()-1 {
			bodyString += "&"
		}
	}
	// Create our x-www-form-urlencoded PUT request
	req := xhr.NewRequest("PUT", fullURL)
	req.Timeout = 1000 // one second, in milliseconds
	req.ResponseType = "text"
	req.SetRequestHeader("Content-Type", "application/x-www-form-urlencoded")
	// Send our request
	err := req.Send(bodyString)
	if err != nil {
		return fmt.Errorf("Something went wrong with PUT request to %s. %s", fullURL, err.Error())
	}
	// Unmarshal our server response object into our model
	err = json.Unmarshal([]byte(req.Response.String()), model)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal response into object, with Error: %s.\nResponse was: %s", err, req.Response.String())
	}
	return nil
}
