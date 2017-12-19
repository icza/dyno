/*
Package dyno is a utility to work with dynamic objects at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON text into an interface{} for example.
*/
package dyno

import "fmt"

// Get returns a value denoted by the path.
// If path is empty or nil, v is returned.
func Get(v interface{}, path ...interface{}) (interface{}, error) {
	for i, el := range path {
		switch node := v.(type) {
		case map[string]interface{}:
			key, ok := el.(string)
			if !ok {
				return nil, fmt.Errorf("expected string path element, got: %T (element idx: %d)", el, i)
			}
			v, ok = node[key]
			if !ok {
				return nil, fmt.Errorf("missing key: %s (path element idx: %d)", key, i)
			}

		case map[interface{}]interface{}:
			var ok bool
			v, ok = node[el]
			if !ok {
				return nil, fmt.Errorf("missing key: %v (path element idx: %d)", el, i)
			}

		case []interface{}:
			idx, ok := el.(int)
			if !ok {
				return nil, fmt.Errorf("expected int path element, got: %T (path element idx: %d)", el, i)
			}
			if idx < 0 || idx >= len(node) {
				return nil, fmt.Errorf("index out of range: %d (path element idx: %d)", idx, i)
			}
			v = node[idx]

		default:
			return nil, fmt.Errorf("expected map or slice node, got: %T (path element idx: %d)", node, i)
		}
	}

	return v, nil
}

// SGet returns a value denoted by the path consisting only string keys.
//
// SGet is an optimized and specialized version of the general Get.
// The path may only contain map keys (no slice indices), and each value associated
// with the keys (being the path elements) must also be maps with string keys,
// except the value asssociated with the last path element.
//
// If path is empty or nil, m is returned.
func SGet(m map[string]interface{}, path ...string) (interface{}, error) {
	if len(path) == 0 {
		return m, nil
	}

	lastIdx := len(path) - 1
	var value interface{}
	var ok bool

	for i, key := range path {
		if value, ok = m[key]; !ok {
			return nil, fmt.Errorf("missing key: %s (path element idx: %d)", key, i)
		}
		if i == lastIdx {
			break
		}
		m2, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected map with string keys node, got: %T (path element idx: %d)", value, i)
		}
		m = m2
	}

	return value, nil
}

// Set sets a map or slice element denoted by the path.
// The map or slice whose element is to be set must already exist.
// Path cannot be empty or nil, else an error is returned.
func Set(v interface{}, value interface{}, path ...interface{}) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	i := len(path) - 1 // The last index
	if len(path) > 1 {
		var err error
		v, err = Get(v, path[:i]...)
		if err != nil {
			return err
		}
	}

	el := path[i]

	switch node := v.(type) {
	case map[string]interface{}:
		key, ok := el.(string)
		if !ok {
			return fmt.Errorf("expected string path element, got: %T (element idx: %d)", el, i)
		}
		node[key] = value

	case map[interface{}]interface{}:
		node[el] = value

	case []interface{}:
		idx, ok := el.(int)
		if !ok {
			return fmt.Errorf("expected int path element, got: %T (path element idx: %d)", el, i)
		}
		if idx < 0 || idx >= len(node) {
			return fmt.Errorf("index out of range: %d (path element idx: %d)", idx, i)
		}
		node[idx] = value

	default:
		return fmt.Errorf("expected map or slice node, got: %T (path element idx: %d)", node, i)
	}

	return nil
}
