/*
Package dyno is a utility to work with dynamic objects at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON text into an `interface{}` for example.
*/
package dyno

import "fmt"

// M is a dynamic map object. It has all the utility methods you need.
type M map[string]interface{}

// S is a dynamic slice object. It has all the utility methods you need.
type S []interface{}

// Value returns a value denoted by the path.
// Path may be empty in which case (nil, nil) is returned.
func Value(v interface{}, path ...interface{}) (interface{}, error) {
	switch node := v.(type) {
	case map[string]interface{}:
		return M(node).Value(path...)
	case []interface{}:
		return S(node).Value(path...)
	default:
		return nil, fmt.Errorf("invalid node type (expecting map or slice, got: %T): %v", node, node)
	}
}

// Value returns a value denoted by the path.
// Path may be empty in which case (nil, nil) is returned.
func (m M) Value(path ...interface{}) (interface{}, error) {
	if len(path) == 0 {
		return nil, nil
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

	return Value(value, path[1:]...)
}

// Value returns a value denoted by the path.
// Path may be empty in which case (nil, nil) is returned.
func (s S) Value(path ...interface{}) (interface{}, error) {
	if len(path) == 0 {
		return nil, nil
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

	return Value(value, path[1:]...)
}
