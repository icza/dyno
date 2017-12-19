# dyno

[![Build Status](https://travis-ci.org/icza/dyno.svg?branch=master)](https://travis-ci.org/icza/dyno)
[![GoDoc](https://godoc.org/github.com/icza/dyno?status.svg)](https://godoc.org/github.com/icza/dyno)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/dyno)](https://goreportcard.com/report/github.com/icza/dyno)
[![codecov](https://codecov.io/gh/icza/dyno/branch/master/graph/badge.svg)](https://codecov.io/gh/icza/dyno)

_Foreword: This package is in an **experimental** phase._

Package dyno is a utility to work with _dynamic objects_ at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON or YAML text into an `interface{}`
for example. When unmarshaling into `interface{}`, libraries usually choose
`map[string]interface{}` or `map[interface{}]interface{}` to represent objects,
and `[]interface{}` to represent arrays.

Package dyno supports dynamic objects that are a mixture of `interface{}`
slices and maps with `interface{}` values and `string` or `interface{}` keys
in any depth and combination.

When operating on a dynamic object, you designate a value you're interested
in by specifying a _path_. A path is a _navigation_; it is a series of keys
and indices that tells how to get to the value.
