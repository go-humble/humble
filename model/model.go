package model

import (
	"encoding/json"
	"fmt"
	"honnef.co/go/js/xhr"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Model interface {
	GetId() string
	RootURL() string
}

// ReadAll expects a pointer to a slice of poitners to some concrete type
// which implements Model (e.g., *[]*Todo). ReadAll will send a GET request to
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
	// Create request
	res, err := http.Get(urlRoot)
	if err != nil {
		return fmt.Errorf("Something went wrong with GET request to %s: %s", urlRoot, err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Couldn't read response to %s: %s", urlRoot, err.Error())
	}
	if err := json.Unmarshal(body, models); err != nil {
		return fmt.Errorf("Failed to unmarshal response into object, with Error: %s.\nResponse was: %s", err, string(body))
	}
	return nil
}

// Read will send a GET request to a RESTful server to get the model by the given id,
// then it will scan the results into model. It expects a json object which contains all
// the fields for the requested model. Read will use the RootURL() method of the model to
// figure out which url to send the GET request to. Typically the full url will look something
// like "http://hostname.com/todos/123"
func Read(id string, model Model) error {
	fullURL := model.RootURL() + "/" + id
	res, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("Something went wrong with GET request to %s: %s", fullURL, err.Error())
	}
	if err := unmarshalResponse(res, model); err != nil {
		return err
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
	encodedString := encodeModelFields(model)
	res, err := http.Post(fullURL, "application/x-www-form-urlencoded", strings.NewReader(encodedString))
	if err != nil {
		return fmt.Errorf("Something went wrong with POST request to %s. %s", fullURL, err.Error())
	}
	if err := unmarshalResponse(res, model); err != nil {
		return err
	}
	return nil
}

// Update expects a pointer some concrete type which implements Model (e.g., *Todo), with a model.Id
// that matches a stored object on the server. It will send a PUT request to the RESTful server.
// It expects a JSON containing the updated object from the server if the request was successful,
// and will set the fields of model with the data in the response object.
// It will use the RootURL() method of the model to determine which url to send the PUT request to.
func Update(model Model) error {
	fullURL := model.RootURL() + "/" + model.GetId()
	encodedString := encodeModelFields(model)
	req, err := http.NewRequest("PUT", fullURL, strings.NewReader(encodedString))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Something went wrong with PUT request to %s. %s", fullURL, err.Error())
	}
	if err := unmarshalResponse(res, model); err != nil {
		return err
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

// unmarshalResponse reads the data from the body of res and then uses the json package to
// unmarshal the data into model.
func unmarshalResponse(res *http.Response, model Model) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Couldn't read response to %s: %s", res.Request.URL.String(), err.Error())
	}
	if err := json.Unmarshal(body, model); err != nil {
		return fmt.Errorf("Failed to unmarshal response into model, with error: %s.\nResponse was: %s", err, string(body))
	}
	return nil
}

// encodeModelFields returns the fields of model represented as a url-encoded string.
// Suitable for POST requests with a content type of application/x-www-form-urlencoded
func encodeModelFields(model Model) string {
	// TODO: Do stronger type checking to prevent errors
	result := ""
	modelValPtr := reflect.ValueOf(model)
	modelVal := modelValPtr.Elem()
	modelValType := modelVal.Type()
	for i := 0; i < modelValType.NumField(); i++ {
		field := modelVal.Field(i)
		fieldName := modelValType.Field(i).Name
		value := fmt.Sprint(field.Interface())
		result += fieldName + "=" + url.QueryEscape(value)
		if i != modelValType.NumField()-1 {
			result += "&"
		}
	}
	return result
}
