// package rest contains functions for sending http requests to a REST API, which
// can be used to create, read, update, and delete models from a server.
//
// TODO: add a really detailed package doc comment describing:
//   - The methods and urls that are used for each function
//   - The format in which models are encoded and what field types are supported
//   - What responses from the server should look like
//   - What happens if there is a non-200 response status code
package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// Model must be satisfied by all models. Satisfying this interface allows you to
// use the helper methods which send http requests to a REST API. They are used
// for e.g., creating a new model or getting an existing model from the server.
// Because of the way reflection is used to encode the data, a Model must have an
// underlying type of either a struct or a pointer to a struct.
type Model interface {
	// GetId returns a unique identifier for the model. It is used for determining
	// which URL to send a request to.
	GetId() string
	// RootURL returns the url for the REST resource corresponding to this model.
	// Typically it should look something like "http://examle.com/todos". Note that
	// the trailing slash should not be included.
	RootURL() string
}

// Create sends an http request to create the given model. It uses reflection to
// convert the fields of model to url-encoded data. Then it sends a POST request to
// model.RootURL() with the encoded data in the body and the Content-Type header set to
// application/x-www-form-urlencoded. It expects a JSON response containing the created
// object from the server if the request was successful, in which case it will mutate model
// by setting the fields to the values in the JSON response. Since model may be mutated,
// it should be a poitner.
func Create(model Model) error {
	fullURL := model.RootURL()
	encodedModelData, err := encodeModelFields(model)
	if err != nil {
		return err
	}
	return sendRequestAndUnmarshal("POST", fullURL, encodedModelData, model)
}

// Read sends an http request to read (or fetch) the model with the given id
// from the server. It sends a GET request to model.RootURL() + "/" + model.GetId().
// Read expects a JSON response containing the data for the requested model if the
// request was successful, in which case it will mutate model by setting the fields
// to the values in the JSON response. Since model may be mutated, it should be
// a pointer.
func Read(id string, model Model) error {
	fullURL := model.RootURL() + "/" + id
	return sendRequestAndUnmarshal("GET", fullURL, "", model)
}

// ReadAll sends an http request to read (or fetch) all the models of a particular
// type from the server (e.g. get all the todos). It sends a GET request to
// model.RootURL(). ReadAll expects a JSON response containing an array of objects,
// where each object contains data for one model. models must be a pointer to a slice
// of some type which implements Model. ReadAll will mutate models by growing or shrinking
// the slice as needed, and by setting the fields of each element to the values in the JSON
// response.
func ReadAll(models interface{}) error {
	rootURL, err := getURLFromModels(models)
	if err != nil {
		return err
	}
	return sendRequestAndUnmarshal("GET", rootURL, "", models)
}

// Update sends an http request to update the given model, i.e. to change some or all
// of the fields. It uses reflection to convert the fields of model to url-encoded data.
// Then it sends a PUT request to model.RootURL() with the encoded data in the body and
// the Content-Type header set to application/x-www-form-urlencoded. Update expects a
// JSON response containing the data for the updated model if the request was successful,
// in which case it will mutate model by setting the fields to the values in the JSON
// response. Since model may be mutated, it should be a pointer.
func Update(model Model) error {
	fullURL := model.RootURL() + "/" + model.GetId()
	encodedModelData, err := encodeModelFields(model)
	if err != nil {
		return err
	}
	return sendRequestAndUnmarshal("PUT", fullURL, encodedModelData, model)
}

// Delete sends an http request to delete the given model. It sends a DELETE request
// to model.RootURL() + "/" + model.GetId(). DELETE expects an empty JSON response
// if the request was successful, and it will not mutate model.
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
		return "", fmt.Errorf("models must be a pointer to a slice of models. The elem type %s does not implement model", typ.Elem().Elem().String())
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

// sendRequestAndUnmarshal constructs a request with the given method, url, and
// data. If data is an empty string, it will construct a request without any
// data in the body. If data is a non-empty string, it will send it as the body
// of the request and set the Content-Type header to
// application/x-www-form-urlencoded. Then sendRequestAndUnmarshal sends the
// request using http.DefaultClient and marshals the response into v using the json
// package.
// TODO: do something if the response status code is non-200.
func sendRequestAndUnmarshal(method string, url string, data string, v interface{}) error {
	// Build the request
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("Something went wrong building %s request to %s: %s", method, url, err.Error())
	}
	// Set the Content-Type header only if data was provided
	if data != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	// Specify that we want json as the response type. This is especially useful
	// for applications which share things between client and server
	req.Header.Set("Accept", "application/json")
	// Send the request using the default client
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Something went wrong with %s request to %s: %s", req.Method, req.URL.String(), err.Error())
	}
	// Unmarshal the response into v
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Couldn't read response to %s: %s", res.Request.URL.String(), err.Error())
	}
	return json.Unmarshal(body, v)
}

// encodeModelFields returns the fields of model represented as a url-encoded string.
// Suitable for POST requests with a content type of application/x-www-form-urlencoded.
// It returns an error if model is a nil pointer or if it is not a struct or a pointer
// to a struct. Any fields that are nil will not be added to the url-encoded string.
func encodeModelFields(model Model) (string, error) {
	modelVal := reflect.ValueOf(model)
	// dereference the pointer until we reach the underlying struct value.
	for modelVal.Kind() == reflect.Ptr {
		if modelVal.IsNil() {
			return "", errors.New("Error encoding model as url-encoded data: model was a nil pointer.")
		}
		modelVal = modelVal.Elem()
	}
	// Make sure the type of model after dereferencing is a struct.
	if modelVal.Kind() != reflect.Struct {
		return "", fmt.Errorf("Error encoding model as url-encoded data: model must be a struct or a pointer to a struct.")
	}
	encodedFields := []string{}
	for i := 0; i < modelVal.Type().NumField(); i++ {
		field := modelVal.Type().Field(i)
		fieldValue := modelVal.FieldByName(field.Name)
		encodedField, err := encodeField(field, fieldValue)
		if err != nil {
			if _, ok := err.(nilFieldError); ok {
				// If there was a nil field, continue without adding the field
				// to the encoded data.
				continue
			}
			// We should return any other kind of error
			return "", err
		}
		encodedFields = append(encodedFields, field.Name+"="+encodedField)
	}
	return strings.Join(encodedFields, "&"), nil
}

type nilFieldError struct{}

func (nilFieldError) Error() string {
	return "field was nil"
}

// encodeField converts a field with the given value to a string. It returns an error
// if field has a type which is unsupported. It returns a special error (nilFieldError)
// if a field has a value of nil. The supported types are int and its variants (int64,
// int32, etc.), uint and its variants (uint64, uint32, etc.), bool, string, and []byte.
func encodeField(field reflect.StructField, value reflect.Value) (string, error) {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			// Skip nil fields
			return "", nilFieldError{}
		}
		value = value.Elem()
	}
	switch v := value.Interface().(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, bool:
		return fmt.Sprint(v), nil
	case string:
		return url.QueryEscape(v), nil
	case []byte:
		return url.QueryEscape(string(v)), nil
	default:
		return "", fmt.Errorf("Error encoding model as url-encoded data: Don't know how to convert %v of type %T to a string.", v, v)
	}
}
