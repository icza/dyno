# dyno

[![Build Status](https://travis-ci.org/icza/dyno.svg?branch=master)](https://travis-ci.org/icza/dyno)
[![GoDoc](https://godoc.org/github.com/icza/dyno?status.svg)](https://godoc.org/github.com/icza/dyno)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/dyno)](https://goreportcard.com/report/github.com/icza/dyno)
[![codecov](https://codecov.io/gh/icza/dyno/branch/master/graph/badge.svg)](https://codecov.io/gh/icza/dyno)

_Foreword: This package is in an **experimental** phase and is a work-in-progress._

Package dyno is a utility to work with _dynamic objects_ at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON or YAML text into an `interface{}`
for example. When unmarshaling into `interface{}`, libraries usually choose
either `map[string]interface{}` or `map[interface{}]interface{}` to represent objects,
and `[]interface{}` to represent arrays. Package dyno supports a mixture of
these in any depth and combination.

When operating on a dynamic object, you designate a value you're interested
in by specifying a _path_. A path is a _navigation_; it is a series of map keys
and `int` slice indices that tells how to get to the value.

Should you need to marshal a dynamic object to JSON which contains maps with
`interface{}` key type (which is not supported by `encoding/json`), you may use
the `ConvertMapI2MapS` converter function.

The implementation does not use reflection at all, so performance is rather good.

Let's see a simple example editing a JSON text to mask out a password. This is
a simplified version of the [`Example_jsonEdit`](https://godoc.org/github.com/icza/dyno#example-package--JsonEdit) example function:

	src := `{"login":{"password":"secret","user":"test"},"name":"compA"}`
	var v interface{}
	if err := json.Unmarshal([]byte(src), &v); err != nil {
		panic(err)
	}
	// Edit (mask out) password:
	if err = dyno.Set(v, "xxx", "login", "password"); err != nil {
		fmt.Printf("Failed to set password: %v\n", err)
	}
	edited, err := json.Marshal(v)
	fmt.Printf("Edited JSON: %s, error: %v\n", edited, err)

Output will be:

	Edited JSON: {"login":{"password":"xxx","user":"test"},"name":"compA"}, error: <nil>
