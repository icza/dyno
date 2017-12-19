/*
Package dyno is a utility to work with dynamic objects at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON text into an `interface{}` for example.
*/
package dyno

import "fmt"

// MS is a dynamic map object with string keys, armed with utility methods.
type MS map[string]interface{}

// MI is a dynamic map object with interface{} keys, armed with utility methods.
type MI map[interface{}]interface{}

// S is a dynamic slice object, armed with utility methods.
type S []interface{}

// Get returns a value denoted by the path.
// If path is empty or nil, v is returned.
func Get(v interface{}, path ...interface{}) (interface{}, error) {
	if len(path) == 0 {
		return v, nil
	}

	switch node := v.(type) {
	case map[string]interface{}:
		return MS(node).Get(path...)
	case []interface{}:
		return S(node).Get(path...)
	case map[interface{}]interface{}:
		return MI(node).Get(path...)
	default:
		return nil, fmt.Errorf("invalid node type (expected map or slice, got: %T): %v", node, node)
	}
}

// Get returns a value denoted by the path.
// If path is empty or nil, m is returned (which will be of type MS).
func (m MS) Get(path ...interface{}) (interface{}, error) {
	if len(path) == 0 {
		return m, nil
	}

	key, ok := path[0].(string)
	if !ok {
		return nil, fmt.Errorf("element is not string: %v", path[0])
	}

	value, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("missing key: %s", key)
	}

	if len(path) == 1 {
		return value, nil
	}

	return Get(value, path[1:]...)
}

// Get returns a value denoted by the path.
// If path is empty or nil, m is returned (which will be of type MI).
func (m MI) Get(path ...interface{}) (interface{}, error) {
	if len(path) == 0 {
		return m, nil
	}

	value, ok := m[path[0]]
	if !ok {
		return nil, fmt.Errorf("missing key: %s", path[0])
	}

	if len(path) == 1 {
		return value, nil
	}

	return Get(value, path[1:]...)
}

// Get returns a value denoted by the path.
// If path is empty or nil, s is returned (which will be of type S).
func (s S) Get(path ...interface{}) (interface{}, error) {
	if len(path) == 0 {
		return s, nil
	}

	idx, ok := path[0].(int)
	if !ok {
		return nil, fmt.Errorf("element is not int (type: %T): %v", path[0], path[0])
	}

	if idx < 0 || idx >= len(s) {
		return nil, fmt.Errorf("index out of range: %d", idx)
	}
	value := s[idx]

	if len(path) == 1 {
		return value, nil
	}

	return Get(value, path[1:]...)
}
