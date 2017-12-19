package dyno

import (
	"reflect"
	"testing"
)

var (
	s = []interface{}{
		1, "a", 3.3, []interface{}{"inner", "inner2"},
	}
	mi = map[interface{}]interface{}{
		"x": 1,
		"y": 2,
		"z": map[interface{}]interface{}{
			3: "three",
		},
	}
	ms = map[string]interface{}{
		"a": 1,
		"p": map[string]interface{}{
			"x": 1,
			"y": 2,
		},
		"pi": mi,
		"ns": []interface{}{1.1, 2.2, 3.3},
		"b":  2,
		"s":  s,
	}
)

func TestGet(t *testing.T) {
	cases := []struct {
		title string        // Title of the test case
		v     interface{}   // Input dynamic object
		path  []interface{} // path whose value to get
		value interface{}   // Expected value
		isErr bool          // Tells if error is expected
	}{
		// Test success:
		{
			title: "nil path on map",
			v:     ms,
			path:  nil,
			value: ms,
		},
		{
			title: "nil path on slice",
			v:     s,
			path:  nil,
			value: s,
		},
		{
			title: "simple map element",
			v:     ms,
			path:  []interface{}{"a"},
			value: 1,
		},
		{
			title: "simple map element #2",
			v:     ms,
			path:  []interface{}{"s"},
			value: s,
		},
		{
			title: "nested map element",
			v:     ms,
			path:  []interface{}{"p", "x"},
			value: 1,
		},
		{
			title: "nested map (mi) element",
			v:     ms,
			path:  []interface{}{"pi", "x"},
			value: 1,
		},
		{
			title: "nested map (mi) element #2",
			v:     ms,
			path:  []interface{}{"pi", "z", 3},
			value: "three",
		},
		{
			title: "nested slice element",
			v:     s,
			path:  []interface{}{3, 1},
			value: "inner2",
		},
		{
			title: "map element and slice element",
			v:     ms,
			path:  []interface{}{"ns", 1},
			value: 2.2,
		},
		{
			title: "map element and slice element #2",
			v:     ms,
			path:  []interface{}{"s", 1},
			value: "a",
		},

		// Test errors:
		{
			title: "expected map or slice node error",
			v:     1,
			path:  []interface{}{"x"},
			isErr: true,
		},
		{
			title: "expected string path element error",
			v:     ms,
			path:  []interface{}{1},
			isErr: true,
		},
		{
			title: "missing key (ms) error",
			v:     ms,
			path:  []interface{}{"x"},
			isErr: true,
		},
		{
			title: "missing key (mi) error",
			v:     mi,
			path:  []interface{}{"a"},
			isErr: true,
		},
		{
			title: "expected int path element error",
			v:     ms,
			path:  []interface{}{"ns", "x"},
			isErr: true,
		},
		{
			title: "index out of range error (negative)",
			v:     ms,
			path:  []interface{}{"ns", -1},
			isErr: true,
		},
		{
			title: "index out of range error (too big)",
			v:     ms,
			path:  []interface{}{"ns", 11},
			isErr: true,
		},
	}

	for _, c := range cases {
		value, err := Get(c.v, c.path...)
		if !reflect.DeepEqual(value, c.value) {
			t.Errorf("[title: %s] Expected value: %v, got: %v", c.title, c.value, value)
		}
		if c.isErr != (err != nil) {
			t.Errorf("[title: %s] Expected error: %v, got: %v, err value: %v", c.title, c.isErr, err != nil, err)
		}
	}
}

func TestSGet(t *testing.T) {
	cases := []struct {
		title string                 // Title of the test case
		v     map[string]interface{} // Input map
		path  []string               // path whose value to get
		value interface{}            // Expected value
		isErr bool                   // Tells if error is expected
	}{
		// Test success:
		{
			title: "nil path on map",
			v:     ms,
			path:  nil,
			value: ms,
		},
		{
			title: "empty path on map",
			v:     ms,
			path:  []string{},
			value: ms,
		},
		{
			title: "simple map element",
			v:     ms,
			path:  []string{"a"},
			value: 1,
		},
		{
			title: "simple map element #2",
			v:     ms,
			path:  []string{"s"},
			value: s,
		},
		{
			title: "nested map element",
			v:     ms,
			path:  []string{"p", "x"},
			value: 1,
		},

		// Test errors:
		{
			title: "missing key error",
			v:     ms,
			path:  []string{"x"},
			isErr: true,
		},
		{
			title: "expected map with string keys node error",
			v:     ms,
			path:  []string{"pi", "x"},
			isErr: true,
		},
		{
			title: "expected map with string keys node error #2",
			v:     ms,
			path:  []string{"ns", "1"},
			isErr: true,
		},
	}

	for _, c := range cases {
		value, err := SGet(c.v, c.path...)
		if !reflect.DeepEqual(value, c.value) {
			t.Errorf("[title: %s] Expected value: %v, got: %v", c.title, c.value, value)
		}
		if c.isErr != (err != nil) {
			t.Errorf("[title: %s] Expected error: %v, got: %v, err value: %v", c.title, c.isErr, err != nil, err)
		}
	}
}

func TestSet(t *testing.T) {
	cases := []struct {
		title string        // Title of the test case
		v     interface{}   // Input dynamic object
		value interface{}   // Value to set
		path  []interface{} // path whose value to set
		exp   interface{}   // Expected result
		isErr bool          // Tells if error is expected
	}{
		// Test success:
		{
			title: "add new map element",
			v:     map[string]interface{}{},
			value: 1,
			path:  []interface{}{"a"},
			exp:   map[string]interface{}{"a": 1},
		},
		{
			title: "change existing map element",
			v:     map[string]interface{}{"a": 1},
			value: 2,
			path:  []interface{}{"a"},
			exp:   map[string]interface{}{"a": 2},
		},
		{
			title: "change existing slice element",
			v:     []interface{}{"a", 1},
			value: 2,
			path:  []interface{}{1},
			exp:   []interface{}{"a", 2},
		},
		{
			title: "change existing map (mi) element",
			v:     map[interface{}]interface{}{1: "one"},
			value: "two",
			path:  []interface{}{1},
			exp:   map[interface{}]interface{}{1: "two"},
		},
		{
			title: "change existing nested map element",
			v: map[string]interface{}{
				"a": map[string]interface{}{"b": 1},
			},
			value: 2,
			path:  []interface{}{"a", "b"},
			exp: map[string]interface{}{
				"a": map[string]interface{}{"b": 2},
			},
		},
		{
			title: "replace existing element with a value of different type",
			v: map[string]interface{}{
				"a": map[string]interface{}{"b": 1},
			},
			value: 2,
			path:  []interface{}{"a"},
			exp:   map[string]interface{}{"a": 2},
		},
		{
			title: "change existing element in map-slice-map",
			v: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"b": 1},
				},
			},
			value: 2,
			path:  []interface{}{"a", 0, "b"},
			exp: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{"b": 2},
				},
			},
		},

		// Test errors:
		{
			title: "change existing map element",
			v:     map[string]interface{}{"a": 1},
			value: 2,
			path:  []interface{}{},
			exp:   map[string]interface{}{"a": 1},
			isErr: true,
		},
		{
			title: "internal Get call returns error",
			v:     map[string]interface{}{"a": 1},
			value: 2,
			path:  []interface{}{"b", "c"},
			exp:   map[string]interface{}{"a": 1},
			isErr: true,
		},
		{
			title: "expected string path element error",
			v:     map[string]interface{}{"a": 1},
			value: 2,
			path:  []interface{}{1},
			exp:   map[string]interface{}{"a": 1},
			isErr: true,
		},
		{
			title: "expected int path element error",
			v:     []interface{}{"a", 1},
			value: 2,
			path:  []interface{}{"a"},
			exp:   []interface{}{"a", 1},
			isErr: true,
		},
		{
			title: "index out of range error (negative)",
			v:     []interface{}{"a", 1},
			value: 2,
			path:  []interface{}{-1},
			exp:   []interface{}{"a", 1},
			isErr: true,
		},
		{
			title: "index out of range error (too big)",
			v:     []interface{}{"a", 1},
			value: 2,
			path:  []interface{}{11},
			exp:   []interface{}{"a", 1},
			isErr: true,
		},
		{
			title: "expected map or slice node error",
			v:     1,
			value: 2,
			path:  []interface{}{"x"},
			exp:   1,
			isErr: true,
		},
	}

	for _, c := range cases {
		err := Set(c.v, c.value, c.path...)
		if !reflect.DeepEqual(c.v, c.exp) {
			t.Errorf("[title: %s] Expected value: %v, got: %v", c.title, c.exp, c.v)
		}
		if c.isErr != (err != nil) {
			t.Errorf("[title: %s] Expected error: %v, got: %v, err value: %v", c.title, c.isErr, err != nil, err)
		}
	}
}
