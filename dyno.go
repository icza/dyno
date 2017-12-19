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

// SGet returns a value denoted by the path consisting string elements.
//
// SGet is an optimized and specialised version of the generic Get.
// The path may only contain map keys (no slice indices), and each value associated
// with the keys (being the path elements) must also be maps with string keys,
// except the value asssociated with the last path element.
//
// If path is empty or nil, m is returned (which will be of type MS).
func (m MS) SGet(path ...string) (interface{}, error) {
	if len(path) == 0 {
		return m, nil
	}

	lastIdx := len(path) - 1
	var value interface{}
	var ok bool

	for i, key := range path {
		if value, ok = m[key]; !ok {
			return nil, fmt.Errorf("missing key: %s", key)
		}
		if i == lastIdx {
			break
		}
		m2, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid node type (expected map with string keys, got: %T): %v", value, value)
		}
		m = MS(m2)
	}

	return value, nil
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
