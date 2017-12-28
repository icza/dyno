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
and int slice indices that tells how to get to the value.

Should you need to marshal a dynamic object to JSON which contains maps with
interface{} key type (which is not supported by encoding/json), you may use
the ConvertMapI2MapS converter function.

The implementation does not uses reflection at all, so performance is rather good.

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
				return nil, fmt.Errorf("expected string path element, got: %T (path element idx: %d)", el, i)
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

// GetInt returns an int value denoted by the path.
//
// If path is empty or nil, v is returned as an int.
func GetInt(v interface{}, path ...interface{}) (int, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("expected int value, got: %T", v)
	}
	return i, nil
}

// GetInteger returns an int64 value denoted by the path.
//
// This function accepts many different types and converts them to int64, namely:
//   -integer types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
//    (which implies the aliases byte and rune too)
//   -floating point types (float64, float32)
//   -string (fmt.Sscan() will be used for parsing)
//   -any type with an Int64() (int64, error) method (e.g. json.Number)
//
// If path is empty or nil, v is returned as an int64.
func GetInteger(v interface{}, path ...interface{}) (int64, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}

	switch i := v.(type) {
	case int64:
		return i, nil
	case int:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int16:
		return int64(i), nil
	case int8:
		return int64(i), nil
	case uint:
		return int64(i), nil
	case uint64:
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint16:
		return int64(i), nil
	case uint8:
		return int64(i), nil
	case float64:
		return int64(i), nil
	case float32:
		return int64(i), nil
	case string:
		var n int64
		_, err := fmt.Sscan(i, &n)
		return n, err
	case interface {
		Int64() (int64, error)
	}:
		return i.Int64()
	default:
		return 0, fmt.Errorf("expected some form of integer number, got: %T", v)
	}
}

// GetFloat64 returns a float64 value denoted by the path.
//
// If path is empty or nil, v is returned as a float64.
func GetFloat64(v interface{}, path ...interface{}) (float64, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("expected float64 value, got: %T", v)
	}
	return f, nil
}

// GetFloating returns a float64 value denoted by the path.
//
// This function accepts many different types and converts them to float64, namely:
//  -floating point types (float64, float32)
//  -integer types (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
//   (which implies the aliases byte and rune too)
//  -string (fmt.Sscan() will be used for parsing)
//  -any type with a Float64() (float64, error) method (e.g. json.Number)
//
// If path is empty or nil, v is returned as an int64.
func GetFloating(v interface{}, path ...interface{}) (float64, error) {
	v, err := Get(v, path...)
	if err != nil {
		return 0, err
	}

	switch f := v.(type) {
	case float64:
		return f, nil
	case float32:
		return float64(f), nil
	case int64:
		return float64(f), nil
	case int:
		return float64(f), nil
	case int32:
		return float64(f), nil
	case int16:
		return float64(f), nil
	case int8:
		return float64(f), nil
	case uint:
		return float64(f), nil
	case uint64:
		return float64(f), nil
	case uint32:
		return float64(f), nil
	case uint16:
		return float64(f), nil
	case uint8:
		return float64(f), nil
	case string:
		var n float64
		_, err := fmt.Sscan(f, &n)
		return n, err
	case interface {
		Float64() (float64, error)
	}:
		return f.Float64()
	default:
		return 0, fmt.Errorf("expected some form of floating point number, got: %T", v)
	}
}

// GetString returns a string value denoted by the path.
//
// If path is empty or nil, v is returned as a string.
func GetString(v interface{}, path ...interface{}) (string, error) {
	v, err := Get(v, path...)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("expected string value, got: %T", v)
	}
	return s, nil
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
			return fmt.Errorf("expected string path element, got: %T (path element idx: %d)", el, i)
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

// ConvertMapI2MapS walks the given dynamic object recursively, and
// converts maps with interface{} key type to maps with string key type.
// This function comes handy if you want to marshal a dynamic object into
// JSON where maps with interface{} key type are not allowed.
//
// Recursion is implemented into values of the following types:
//   -map[interface{}]interface{}
//   -map[string]interface{}
//   -[]interface{}
//
// When converting map[interface{}]interface{} to map[string]interface{},
// fmt.Sprint() with default formatting is used to convert the key to a string key.
func ConvertMapI2MapS(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string: // Fast check if it's already a string
				m[k2] = ConvertMapI2MapS(v2)
			default:
				m[fmt.Sprint(k)] = ConvertMapI2MapS(v2)
			}
		}
		v = m

	case []interface{}:
		for i, v2 := range x {
			x[i] = ConvertMapI2MapS(v2)
		}

	case map[string]interface{}:
		for k, v2 := range x {
			x[k] = ConvertMapI2MapS(v2)
		}
	}

	return v
}
