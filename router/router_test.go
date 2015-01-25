package router

import (
	"reflect"
	"testing"
)

// routeTestCases is format for test case. Each test case takes possible route paths, the URL hash provided
// and the expected route path and query parameters.
var routeTestCases = []struct {
	paths          []string
	hash           string
	expectedPath   string
	expectedParams map[string]string
}{
	// Test route to /
	{
		paths:          []string{"/home", "/about", "/"},
		hash:           "#/",
		expectedPath:   "/",
		expectedParams: nil,
	},
	// Test basic literal
	{
		paths:          []string{"/home", "/about"},
		hash:           "#/home",
		expectedPath:   "/home",
		expectedParams: nil,
	},
	// Test trailing slash in hash
	{
		paths:          []string{"/home", "/about"},
		hash:           "#/about/",
		expectedPath:   "/about",
		expectedParams: nil,
	},
	// Test query param
	{
		paths:        []string{"/home", "/home/{homeId}", "/about"},
		hash:         "#/home/55",
		expectedPath: "/home/{homeId}",
		expectedParams: map[string]string{
			"homeId": "55",
		},
	},
	// Test trailing slash after query param in hash
	{
		paths:        []string{"/home", "/about", "/about/{aboutId}"},
		hash:         "#/about/55/",
		expectedPath: "/about/{aboutId}",
		expectedParams: map[string]string{
			"aboutId": "55",
		},
	},
	// Test multiple query params
	{
		paths:        []string{"/home", "/home/{homeId}", "/about", "/home/{homeId}/image/{imageSize}/{imageType}/jpg"},
		hash:         "#/home/55/image/800x600/panoramic/jpg",
		expectedPath: "/home/{homeId}/image/{imageSize}/{imageType}/jpg",
		expectedParams: map[string]string{
			"homeId":    "55",
			"imageSize": "800x600",
			"imageType": "panoramic",
		},
	},
	// Test tie breaker
	{
		paths:          []string{"/home/all", "/home/{homeId}", "/about"},
		hash:           "#/home/all",
		expectedPath:   "/home/all",
		expectedParams: nil,
	},
}

func TestRouter(t *testing.T) {
	for _, tc := range routeTestCases {
		gotPath := ""
		gotParams := map[string]string{}
		r := New()
		for _, path := range tc.paths {
			handler := generateTestHandler(path, &gotPath, &gotParams)
			r.HandleFunc(path, handler)
		}
		r.hashChanged(tc.hash)
		if gotPath != tc.expectedPath {
			t.Errorf("Failed for hash=%s. Expected path: %s, Got path: %s", tc.hash, tc.expectedPath, gotPath)
		}
		if tc.expectedParams != nil {
			if !reflect.DeepEqual(gotParams, tc.expectedParams) {
				t.Errorf("Failed for hash=%s. Expected params: %s, Got params: %s", tc.hash, tc.expectedParams, gotParams)
			}
		}
	}
}

func generateTestHandler(path string, gotPath *string, gotParams *map[string]string) Handler {
	return func(params map[string]string) {
		*gotPath = path
		*gotParams = params
	}
}
