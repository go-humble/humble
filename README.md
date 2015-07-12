Humble
======

Humble is a collection of loosely-coupled tools designed to build client-side
and hybrid web applications using go and
[gopherjs](https://github.com/gopherjs/gopherjs). This repository contains no
code, but serves as an introduction to the Humble Toolkit and a central place
for creating new issues that are not related to any specific sub-package.


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
significantly as new issues are discovered and new features are added. Since
the API will be rapidly changing, we recommend using a dependency tool such as
[godep](https://github.com/tools/godep) to ensure that your code does not break
unexpectedly.


Packages
--------

All packages are organized under the
[Go-Humble Organization on GitHub](https://github.com/go-humble). They are
designed to be either used together or separately.

### [Detect](https://github.com/go-humble/detect)

Detect is a tiny go package for detecting whether code is running on the server
or browser. It is intended to be used in hybrid go applications.

### [Examples](https://github.com/go-humble/examples)

Examples contains several examples of how to use the Humble Toolkit to build
real applications.

### [Rest](https://github.com/go-humble/rest)

Rest is a small package for sending requests to a RESTful API and unmarshaling
the response. Rest sends requests using CRUD semantics. It supports requests
with a Content-Type of either application/x-www-form-urlencoded or
application/json and parses json responses from the server.

### [View](https://github.com/go-humble/view)

View is a small package for organizing view-related code. View includes a View
interface and some helper functions for operating on views (e.g. Append,
Replace, Remove, etc.).

### [Temple](https://github.com/go-humble/temple)

Temple is a library and a command line tool for sanely managing go templates,
with the ability to share them between the server and browser.

### [Router](https://github.com/go-humble/router)

Router is an easy-to-use router which runs in the browser. It supports url
parameters and uses history.pushState, gracefully falling back to url hashes
if needed.


Where is the Old Code?
----------------------

If you're looking for the files that used to be in this repository, they have
all been moved to stand-alone packages. Check out the
[Humble Toolkit](https://github.com/go-humble) on github to view all the packages!
