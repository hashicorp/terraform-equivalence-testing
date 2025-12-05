// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package json

import "strconv"

// Step represents a step in the path to a field that should be stripped from
// the input data.
type Step struct {
	Step   string
	Filter []Filter
}

// Filter represents a filter that should be validated before a Field is
// stripped.
type Filter struct {
	Path  []string
	Value interface{}
}

func (step Step) applyFilter(data interface{}) bool {
	for _, f := range step.Filter {
		if !filter(f.Path, f.Value, data) {
			return false
		}
	}
	return true
}

func filter(parts []string, target interface{}, data interface{}) bool {
	if data == nil {
		return false
	}

	if len(parts) == 0 {
		return target == data
	}

	switch data := data.(type) {
	case map[string]interface{}:
		return filter(parts[1:], target, data[parts[0]])
	case []interface{}:
		ix, err := strconv.Atoi(parts[0])
		if err != nil {
			return false
		}
		if ix < 0 || ix >= len(data) {
			return false
		}
		return filter(parts[1:], target, data[ix])
	default:
		return false
	}
}
