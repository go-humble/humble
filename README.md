Humble
======

Humble is a collection of loosely-coupled tools designed to build client-side
and hybrid web applications using go and
[gopherjs](https://github.com/gopherjs/gopherjs).

Humble is designed for writing front-end code, and is entirely back-end
agnostic. You can use Humble in combination with any back-end server written in
any language. If you do choose to write your back-end in go as well, Humble
offers some tools to make it easy to share code between the server and browser.

This repository contains no code, but serves as an introduction to Humble and a
central place for creating new issues and feature requests that are not related
to any specific sub-package.

How it Works
------------

Humble allows you to write front-end code in pure go, which you can then
compile to javascript with [gopherjs](https://github.com/gopherjs/gopherjs) and
run in the browser. Humble is pure go. It feels like go, compiles with the
standard go tools, and follows go idioms when possible.

[GopherJS](https://github.com/gopherjs/gopherjs) supports all the main features
of the go language, including goroutines and channels. Most of the standard
library is also supported (see the
[compatibility table](https://github.com/gopherjs/gopherjs/blob/master/doc/packages.md)).


Why Write Client-Side Code in Go?
---------------------------------

Ultimately, Humble is not for everyone and may not be suitable for all projects.
Many developers will be perfectly happy writing front-end code in javascript,
and the javascript ecosystem is vast and thriving. It is not our goal to replace
javascript or convince every developer to switch. However, we recognize that
javascript is not everyone's favorite language, and Humble offers developers,
especially those already familiar with go, a viable alternative.

Go offers several benefits over javascript for writing client-side code:

1. **Built-In Type-Safety**. Go is a type-safe, compiled language which means
   that certain classes of mistakes can be caught before you even run your code.
   It also makes it possible for text editors to do static analysis and more
   intelligent autocomplete, increasing your productivity. While projects like
   [TypeScript](http://www.typescriptlang.org/) exist, go offers type-safety as
   a core feature, and *all* go code, whether part of the standard library or a
   third-party package, supports it.
2. **Robust Standard Library**. Go comes with an incredibly robust
   [standard library](https://golang.org/pkg/), and almost all of it is
   supported in gopherjs (see
   the [compatibility table](https://github.com/gopherjs/gopherjs/blob/master/doc/packages.md)).
3. **The Ability to Build Hybrid Applications**. If you already have a server
   written in go, you can use Humble to write the front-end in go too. Having
   your entire codebase in one language reduces cognitive load and increases
   maintainability. It is even possible to share code between the server and
   browser, just like you can with javascript and node.js.
4. **Sane Concurrency Patterns**. Go is one of few modern languages to implement
   [CSP concurrency patterns](https://en.wikipedia.org/wiki/Communicating_sequential_processes).
   Goroutines and channels offer an intuitive and standardized way to deal with
   concurrency and asynchronous code, and are fully supported in gopherjs.
5. **Phenomenal Tooling**. Go offers some of the best tooling of any modern
   language. There is standardized documentation on
   [godoc.org](http://godoc.org/), builtin
   [testing, benchmarking](http://golang.org/pkg/testing/), and
   [profiling](http://blog.golang.org/profiling-go-programs), and even tools for
   [detecting race conditions](http://blog.golang.org/race-detector) and
   [displaying test coverage](https://blog.golang.org/cover).


Development Status
------------------

Humble is brand new and is under active development.

All sub-packages are well-tested and are even tested in real browsers when
applicable. We do our best to respond to critical issues as quickly as possible.
As such, Humble can be considered safe for use in side-projects, experiments,
and non-critical production applications. At this time, we do not recommend
using Humble for critical production applications.

Humble uses semantic versioning, but offers no guarantees of backwards
compatibility until version 1.0. It is likely that the API will change
significantly as new issues are discovered and new features are added. As such,
we recommend using a dependency tool such as
[godep](https://github.com/tools/godep) to ensure that your code does not break
unexpectedly.


Packages
--------

In contrast with front-end javascript frameworks such as Angular and Ember,
Humble is much more conservative. It doesn't enforce any kind of structure on
your applications, but tries to provide all the tools you need for the most
common use cases. The packages that make up Humble are loosely-coupled, which
means they work well together but can be used separately too. Humble can be
used with [gopherjs bindings](https://github.com/gopherjs/gopherjs/wiki/bindings),
such as [jQuery](https://github.com/gopherjs/jquery). It is even possible to
use Humble together
[with existing javascript code](https://github.com/gopherjs/gopherjs#interacting-with-the-dom).

### [Detect](https://github.com/go-humble/detect)

Detect is a tiny go package for detecting whether code is running on the server
or browser. It is intended to be used in hybrid go applications.

### [Examples](https://github.com/go-humble/examples)

Examples contains several examples of how to use Humble to build real
applications.

### [Form](https://github.com/go-humble/form)

Form is a package for validating and serializing html forms in the browser. It
supports a variety of validations on form inputs and binding forms to arbitrary
go structs.

### [Locstor](https://github.com/go-humble/locstor)

Locstor provides localStorage bindings. In addition to being able to interact
with the localStorage API directly, you can create a DataStore object for
storing and retrieving arbitrary go data structures, not just strings.

### [Rest](https://github.com/go-humble/rest)

Rest is a small package for sending requests to a RESTful API and unmarshaling
the response. Rest sends requests using CRUD semantics. It supports requests
with a Content-Type of either application/x-www-form-urlencoded or
application/json and parses json responses from the server.

### [Router](https://github.com/go-humble/router)

Router is an easy-to-use router which runs in the browser. It supports url
parameters and uses history.pushState, gracefully falling back to url hashes
if needed.

### [Temple](https://github.com/go-humble/temple)

Temple is a library and a command line tool for sanely managing go templates,
with the ability to share them between the server and browser.

### [View](https://github.com/go-humble/view)

View is a small package for organizing view-related code. View includes a View
interface and some helper functions for operating on views (e.g. Append,
Replace, Remove, etc.).


Where is the Old Code?
----------------------

If you're looking for the files that used to be in this repository, they have
all been moved to stand-alone packages. Check out the
[Go-Humble Organization Page](https://github.com/go-humble) on github to view
all the packages!
