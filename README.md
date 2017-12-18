# dyno

[![Build Status](https://travis-ci.org/icza/dyno.svg?branch=master)](https://travis-ci.org/icza/dyno)
[![GoDoc](https://godoc.org/github.com/icza/dyno?status.svg)](https://godoc.org/github.com/icza/dyno)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/dyno)](https://goreportcard.com/report/github.com/icza/dyno)
[![codecov](https://codecov.io/gh/icza/dyno/branch/master/graph/badge.svg)](https://codecov.io/gh/icza/dyno)

_Foreword: This package is in an **experimental** phase._

Package dyno is a utility to work with _dynamic objects_ at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON text into an `interface{}` for example.
