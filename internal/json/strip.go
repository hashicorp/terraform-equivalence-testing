// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package json

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	wildcard = "*"
)

func Field(path string) []Step {
	var field []Step
	for _, step := range strings.Split(path, ".") {
		field = append(field, Step{Step: step})
	}
	return field
}

// Strip mutates the input data by removing all the required fields.
//
// Check out the strip_test.go test cases for examples of the accepted format
// for each field.
func Strip(fields [][]Step, data interface{}) (interface{}, error) {
	for _, field := range fields {
		var err error
		data, err = strip(field, data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func strip(steps []Step, current interface{}) (interface{}, error) {
	if current == nil {
		return nil, nil
	}

	if len(steps) == 1 {
		return stripLeaf(steps[0], current)
	}

	return stripNode(steps, current)
}

func stripLeaf(part Step, current interface{}) (interface{}, error) {
	switch leaf := current.(type) {
	case map[string]interface{}:
		return stripMapLeaf(part, leaf), nil
	case []interface{}:
		return stripSliceLeaf(part, leaf)
	default:
		return nil, fmt.Errorf("unrecognised json type: %T", leaf)
	}
}

func stripMapLeaf(part Step, current map[string]interface{}) map[string]interface{} {
	switch part.Step {
	case wildcard:
		remaining := make(map[string]interface{})
		for key, value := range current {
			if !part.applyFilter(value) {
				remaining[key] = value
			}
		}
		return remaining
	default:
		if next, ok := current[part.Step]; ok {
			if !part.applyFilter(next) {
				return current
			}
		}
		delete(current, part.Step)
		return current
	}
}

func stripSliceLeaf(part Step, current []interface{}) ([]interface{}, error) {
	switch part.Step {
	case wildcard:
		remaining := make([]interface{}, 0)
		for _, item := range current {
			if !part.applyFilter(item) {
				remaining = append(remaining, item)
			}
		}
		return remaining, nil
	default:
		ix, err := strconv.Atoi(part.Step)
		if err != nil {
			return nil, fmt.Errorf("must specify an integer when referencing json arrays, instead specified %s", part)
		}

		if ix < 0 || ix >= len(current) {
			return nil, fmt.Errorf("index %d out of bounds for array of length %d", ix, len(current))
		}

		if !part.applyFilter(current[ix]) {
			return current, nil
		}

		return append(current[:ix], current[ix+1:]...), nil
	}
}

func stripNode(parts []Step, current interface{}) (interface{}, error) {
	switch node := current.(type) {
	case map[string]interface{}:
		return stripMapNode(parts, node)
	case []interface{}:
		return stripSliceNode(parts, node)
	default:
		return nil, fmt.Errorf("unrecognized json type: %T", node)
	}
}

func stripMapNode(parts []Step, current map[string]interface{}) (map[string]interface{}, error) {
	switch parts[0].Step {
	case wildcard:
		ret := map[string]interface{}{}
		for key, value := range current {
			if !parts[0].applyFilter(value) {
				ret[key] = value
				continue
			}

			var err error
			if ret[key], err = strip(parts[1:], value); err != nil {
				return nil, err
			}
		}
		return ret, nil
	default:
		if _, ok := current[parts[0].Step]; !ok {
			// If the JSON object doesn't have this path, just skip it.
			return current, nil
		}

		if !parts[0].applyFilter(current[parts[0].Step]) {
			return current, nil
		}

		var err error
		if current[parts[0].Step], err = strip(parts[1:], current[parts[0].Step]); err != nil {
			return nil, err
		}
		return current, nil
	}
}

func stripSliceNode(parts []Step, current []interface{}) ([]interface{}, error) {
	switch parts[0].Step {
	case wildcard:
		ret := make([]interface{}, 0)
		for _, item := range current {
			if !parts[0].applyFilter(item) {
				ret = append(ret, item)
				continue
			}

			stripped, err := strip(parts[1:], item)
			if err != nil {
				return nil, err
			}
			ret = append(ret, stripped)
		}
		return ret, nil
	default:
		ix, err := strconv.Atoi(parts[0].Step)
		if err != nil {
			return nil, fmt.Errorf("must specify an integer when referencing json arrays, instead specified %s", parts[0])
		}

		if ix < 0 || ix >= len(current) {
			return nil, fmt.Errorf("index %d out of bounds for array of length %d", ix, len(current))
		}

		if !parts[0].applyFilter(current[ix]) {
			return current, nil
		}

		if current[ix], err = strip(parts[1:], current[ix]); err != nil {
			return nil, err
		}
		return current, nil
	}
}
