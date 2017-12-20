/*
Package dyno is a utility to work with dynamic objects at ease.

Primary goal is to easily handle dynamic objects and arrays (and a mixture of these)
that are the result of unmarshaling a JSON or YAML text into an interface{}
for example. When unmarshaling into interface{}, libraries usually choose
either map[string]interface{} or map[interface{}]interface{} to represent objects,
and []interface{} to represent arrays. Package dyno supports a mixture of
these in any depth and combination.

When operating on a dynamic object, you designate a value you're interested
in by specifying a path. A path is a navigation; it is a series of map keys
and slice indices that tells how to get to the value.

*/
package dyno

import "fmt"

// Get returns a value denoted by the path.
//
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

// SGet returns a value denoted by the path consisting of only string keys.
//
// SGet is an optimized and specialized version of the general Get.
// The path may only contain string map keys (no slice indices), and
// each value associated with the keys (being the path elements) must also
// be maps with string keys, except the value asssociated with the last
// path element.
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
//
// The last element of the path must be a map key or a slice index, and the
// preceeding path must denote a map or a slice respectively which must already exist.
//
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

// SSet sets a map element with string key type, denoted by the path
// consisting of only string keys.
//
// SSet is an optimized and specialized version of the general Set. The
// path may only contain string map keys (no slice indices), and each
// value associated with the keys (being the path elements) must also be
// a maps with string keys, except the value associated with the last path
// element.
//
// The map denoted by the preceeding path before the last path element
// must already exist.
//
// Path cannot be empty or nil, else an error is returned.
func SSet(m map[string]interface{}, value interface{}, path ...string) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	i := len(path) - 1 // The last index
	if len(path) > 1 {
		v, err := SGet(m, path[:i]...)
		if err != nil {
			return err
		}

		var ok bool
		m, ok = v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map with string keys node, got: %T (path element idx: %d)", value, i)
		}
	}

	m[path[i]] = value
	return nil
}

// Append appends a value to a slice denoted by the path.
//
// The slice denoted by path must already exist.
//
// Path cannot be empty or nil, else an error is returned.
func Append(v interface{}, value interface{}, path ...interface{}) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	node, err := Get(v, path...)
	if err != nil {
		return err
	}

	s, ok := node.([]interface{})
	if !ok {
		return fmt.Errorf("expected slice node, got: %T (path element idx: %d)", node, len(path))
	}

	// Must set the new slice value:
	return Set(v, append(s, value), path...)
}
