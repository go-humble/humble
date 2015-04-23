package model

import (
	"encoding/json"
	"fmt"
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

// Create expects a pointer some concrete type which implements Model (e.g., *Todo).
// It will send a POST request to the RESTful server. It expects a JSON containing the
// created object from the server if the request was successful, and will set the fields of
// model with the data in the response object. It will use the RootURL() method of
// the model to determine which url to send the POST request to.
func Create(model Model) error {
	fullURL := model.RootURL()
	encodedString := encodeModelFields(model)
	return sendRequestWithDataAndUnmarshal("POST", fullURL, encodedString, model)
}

// Read will send a GET request to a RESTful server to get the model by the given id,
// then it will scan the results into model. It expects a json object which contains all
// the fields for the requested model. Read will use the RootURL() method of the model to
// figure out which url to send the GET request to. Typically the full url will look something
// like "http://hostname.com/todos/123"
func Read(id string, model Model) error {
	fullURL := model.RootURL() + "/" + id
	return sendRequestWithoutDataAndUnmarshal("GET", fullURL, model)
}

// ReadAll expects a pointer to a slice of poitners to some concrete type
// which implements Model (e.g., *[]*Todo). ReadAll will send a GET request to
// a RESTful server and scan the results into models. It expects a json array
// of json objects from the server, where each object represents a single Model
// of some concrete type. It will use the RootURL() method of the models to
// figure out which url to send the GET request to.
func ReadAll(models interface{}) error {
	rootURL, err := getURLFromModels(models)
	if err != nil {
		return err
	}
	return sendRequestWithoutDataAndUnmarshal("GET", rootURL, models)
}

// getURLFromModels returns the url that should be used for the type that corresponds
// to models. It does this by instantiating a new model of the correct type and then
// calling RootURL on it. models should be a pointer to a slice of models.
func getURLFromModels(models interface{}) (string, error) {
	// Check the type of models
	typ := reflect.TypeOf(models)
	switch {
	// Make sure its a pointer
	case typ.Kind() != reflect.Ptr:
		return "", fmt.Errorf("models must be a pointer to a slice of models. %T is not a pointer.", models)
	// Make sure its a pointer to a slice
	case typ.Elem().Kind() != reflect.Slice:
		return "", fmt.Errorf("models must be a pointer to a slice of models. %T is not a pointer to a slice", models)
	// Make sure the type of the elements of the slice implement Model
	case !typ.Elem().Elem().Implements(reflect.TypeOf([]Model{}).Elem()):
		return "", fmt.Errorf("models must be a pointer to a slice of models. The elem type %T does not implement model", typ.Elem().Elem())
	}
	// modelType is the type of the elements of models
	modelType := typ.Elem().Elem()
	// Ultimately, we need to be able to instantiate a new object of a type that
	// implements Model so that we can call RootURL on it. The trouble is that
	// reflect.New only works for things that are not pointers, and the type of
	// the elements of models could be pointers. To solve for this, we are going
	// to get the Elem of modelType if it is a pointer and keep track of the number
	// of times we get the Elem. So if modelType is *Todo, we'll call Elem once to
	// get the type Todo.
	numDeref := 0
	for modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		numDeref += 1
	}
	// Now that we have the underlying type that is not a pointer, we can instantiate
	// a new object with reflect.New.
	newModelVal := reflect.New(modelType).Elem()
	// Now we need to iteratively get the address of the object we created exactly
	// numDeref times to get to a type that implements Model. Note that Addr is the
	// inverse of Elem.
	for i := 0; i < numDeref; i++ {
		newModelVal = newModelVal.Addr()
	}
	// Now we can use a type assertion to convert the object we instantiated to a Model
	newModel := newModelVal.Interface().(Model)
	// Finally, once we have a Model we can get what we wanted by calling RootURL
	return newModel.RootURL(), nil
}

// Update expects a pointer some concrete type which implements Model (e.g., *Todo), with a model.Id
// that matches a stored object on the server. It will send a PUT request to the RESTful server.
// It expects a JSON containing the updated object from the server if the request was successful,
// and will set the fields of model with the data in the response object.
// It will use the RootURL() method of the model to determine which url to send the PUT request to.
func Update(model Model) error {
	fullURL := model.RootURL() + "/" + model.GetId()
	encodedString := encodeModelFields(model)
	return sendRequestWithDataAndUnmarshal("PUT", fullURL, encodedString, model)
}

// Delete expects a pointer some concrete type which implements Model (e.g., *Todo).
// It will send a DELETE request to a RESTful server. It expects an empty json
// object from the server if the request was successful, and will not attempt to do anything
// with the response. It will use the RootURL() and GetId() methods of the model to determine
// which url to send the DELETE request to. Typically, the full url will look something
// like "http://hostname.com/todos/123"
func Delete(model Model) error {
	fullURL := model.RootURL() + "/" + model.GetId()
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("Something went wrong building DELETE request to %s: %s", fullURL, err.Error())
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		return fmt.Errorf("Something went wrong with DELETE request to %s: %s", fullURL, err.Error())
	}
	return nil
}

// sendRequestWithoutDataAndUnmarshal creates a request with the given method and url, sends
// it using the default client, and then marshals the json response into v.
func sendRequestWithoutDataAndUnmarshal(method string, url string, v interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return fmt.Errorf("Something went wrong building %s request to %s: %s", method, url, err.Error())
	}
	return sendRequestAndUnmarshal(req, v)
}

// sendRequestWithDataAndUnmarshal creates a request with the given method, url, and data, sends
// it using the default client, and then marshals the json response into v. It uses the Content-Type
// header application/x-www-form-urlencoded.
func sendRequestWithDataAndUnmarshal(method string, url string, data string, v interface{}) error {
	req, err := buildRequestWithData(method, url, data)
	if err != nil {
		return err
	}
	return sendRequestAndUnmarshal(req, v)
}

// buildRequestWithData creates and returns a request with the given method, url, and data. It
// also sets the Content-Type header to application/x-www-form-urlencoded.
func buildRequestWithData(method, url, data string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Something went wrong building %s request to %s: %s", method, url, err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// sendRequestAndUnmarshal sends req using http.DefaultClient and then marshals the response into v.
// TODO: do something if the response status code is non-200.
func sendRequestAndUnmarshal(req *http.Request, v interface{}) error {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Something went wrong with %s request to %s: %s", req.Method, req.URL.String(), err.Error())
	}
	return unmarshalResponse(res, v)
}

// unmarshalResponse reads the data from the body of res and then uses the json package to
// unmarshal the data into v.
func unmarshalResponse(res *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Couldn't read response to %s: %s", res.Request.URL.String(), err.Error())
	}
	return json.Unmarshal(body, v)
}

// encodeModelFields returns the fields of model represented as a url-encoded string.
// Suitable for POST requests with a content type of application/x-www-form-urlencoded
func encodeModelFields(model Model) string {
	// TODO: Do stronger type checking to prevent errors
	encodedString := ""
	modelValPtr := reflect.ValueOf(model)
	modelVal := modelValPtr.Elem()
	modelValType := modelVal.Type()
	for i := 0; i < modelValType.NumField(); i++ {
		field := modelVal.Field(i)
		fieldName := modelValType.Field(i).Name
		value := fmt.Sprint(field.Interface())
		encodedString += fieldName + "=" + url.QueryEscape(value)
		if i != modelValType.NumField()-1 {
			encodedString += "&"
		}
	}
	return encodedString
}
