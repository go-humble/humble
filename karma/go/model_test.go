package main

import (
	"github.com/JohannWeging/jasmine"
)

func main() {
	// This test is just checks that the code cross-compiled correctly and can
	// be executed by the karma test runner.
	jasmine.Describe("Tests", func() {
		jasmine.It("can be loaded", func() {
			jasmine.Expect(true).ToBe(true)
		})
	})
}
