// Package bind provides functions for watching for changes to go objects/variables
// and binding go objects/variables to different elements in the DOM. In order to work
// correctly, it requires watch.js: https://github.com/melanke/Watch.JS/
package bind

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"reflect"
)

// Watch allows you to listen for and react to changes to v. If v changes, f will
// be called. v must be a pointer.
func Watch(v interface{}, f func()) error {
	// Make sure the type is a pointer
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("First argument to Watch must be a pointer. Got %T", v)
	}
	// Invoke watch.js with the given parameters
	js.Global.Call("watch", js.InternalObject(v), f)
	return nil
}

// WatchField allows you to listen for and react to changes for a specific field in
// v. If that field changes, f will be called. If other fields in v change, f will not
// be called. v must be a pointer to a struct, and fieldName must be the name of one of
// the fields of v.
func WatchField(v interface{}, fieldName string, f func()) error {
	typ := reflect.TypeOf(v)
	// Make sure the type is a pointer to a struct
	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("First argument to WatchField must be a pointer. Got %T", v)
	}
	structTyp := typ.Elem()
	if structTyp.Kind() != reflect.Struct {
		return fmt.Errorf("First argument to WatchField must be a pointer to a struct. Got %T, which is a pointer to %s", v, structTyp.Kind())
	}
	// Make sure there is some field in v matching fieldName
	attrFound := false
	for i := 0; i < structTyp.NumField(); i++ {
		field := structTyp.Field(i)
		if field.Name == fieldName {
			attrFound = true
			break
		}
	}
	if !attrFound {
		return fmt.Errorf("Error in WatchField: type %T does not have field named %s", v, fieldName)
	}
	// Invoke watch.js with the given parameters
	js.Global.Call("watch", js.InternalObject(v), fieldName, f)
	return nil
}
