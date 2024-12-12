// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestStripJson(t *testing.T) {
	tcs := []struct {
		input    interface{}
		expected interface{}
		fields   [][]Step
	}{
		{
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
			fields:   [][]Step{},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			fields: [][]Step{},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
			},
			fields: [][]Step{
				Field("list"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{},
			},
			fields: [][]Step{
				Field("list.*"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": []interface{}{
						"one",
						"two",
					},
					"two": []interface{}{
						"one",
						"two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"two": []interface{}{
						"one",
						"two",
					},
				},
			},
			fields: [][]Step{
				Field("map.one"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": []interface{}{
						"one",
						"two",
					},
					"two": []interface{}{
						"one",
						"two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{},
			},
			fields: [][]Step{
				Field("map.*"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": []interface{}{
						"one",
						"two",
					},
					"two": []interface{}{
						"one",
						"two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"one": []interface{}{},
					"two": []interface{}{
						"one",
						"two",
					},
				},
			},
			fields: [][]Step{
				Field("map.one.*"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"two": "two",
					},
					map[string]interface{}{
						"two": "two",
					},
				},
			},
			fields: [][]Step{
				Field("map.one"),
				Field("list.*.one"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			expected: map[string]interface{}{
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			fields: [][]Step{
				Field("map"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
						"two": "two",
					},
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
				"list": []interface{}{
					map[string]interface{}{
						"two": "two",
					},
					map[string]interface{}{
						"one": "one",
					},
				},
			},
			fields: [][]Step{
				Field("list.0.one"),
				Field("list.1.two"),
			},
		},
		{
			input: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
			},
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"one": "one",
					"two": "two",
				},
			},
			fields: [][]Step{
				Field("other_map.one"),
			},
		},
		{
			input: []interface{}{
				map[string]interface{}{
					"value":  "one",
					"filter": "true",
				},
				map[string]interface{}{
					"value":  "two",
					"filter": "false",
				},
				map[string]interface{}{
					"value":  "three",
					"filter": "true",
				},
				map[string]interface{}{
					"value":  "four",
					"filter": "false",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"value":  "two",
					"filter": "false",
				},
				map[string]interface{}{
					"value":  "four",
					"filter": "false",
				},
			},
			fields: [][]Step{
				{
					{
						Step: wildcard,
						Filter: []Filter{
							{
								Path:  []string{"filter"},
								Value: "true",
							},
						},
					},
				},
			},
		},
		{
			input: []interface{}{
				map[string]interface{}{
					"value":  "one",
					"filter": "true",
				},
				map[string]interface{}{
					"value":  "two",
					"filter": "false",
				},
				map[string]interface{}{
					"value":  "three",
					"filter": "true",
				},
				map[string]interface{}{
					"value":  "four",
					"filter": "false",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"filter": "true",
				},
				map[string]interface{}{
					"value":  "two",
					"filter": "false",
				},
				map[string]interface{}{
					"filter": "true",
				},
				map[string]interface{}{
					"value":  "four",
					"filter": "false",
				},
			},
			fields: [][]Step{
				{
					{
						Step: wildcard,
						Filter: []Filter{
							{
								Path:  []string{"filter"},
								Value: "true",
							},
						},
					},
					{
						Step: "value",
					},
				},
			},
		},
		{
			input: []interface{}{
				map[string]interface{}{
					"value": "one",
				},
				map[string]interface{}{
					"value": "two",
				},
				map[string]interface{}{
					"value": "three",
				},
				map[string]interface{}{
					"value": "four",
				},
			},
			expected: []interface{}{
				map[string]interface{}{},
				map[string]interface{}{
					"value": "two",
				},
				map[string]interface{}{
					"value": "three",
				},
				map[string]interface{}{
					"value": "four",
				},
			},
			fields: [][]Step{
				{
					{
						Step: wildcard,
					},
					{
						Step: "value",
						Filter: []Filter{
							{
								Value: "one",
							},
						},
					},
				},
			},
		},
		{
			input: []interface{}{
				map[string]interface{}{
					"value": "one",
					"tags": map[string]interface{}{
						"env": "prod",
					},
				},
				map[string]interface{}{
					"value": "two",
					"tags": map[string]interface{}{
						"env": "dev",
					},
				},
				map[string]interface{}{
					"value": "three",
					"tags": map[string]interface{}{
						"env": "prod",
					},
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"value": "two",
					"tags": map[string]interface{}{
						"env": "dev",
					},
				},
			},
			fields: [][]Step{
				{
					{
						Step: wildcard,
						Filter: []Filter{
							{
								Path:  []string{"tags", "env"},
								Value: "prod",
							},
						},
					},
				},
			},
		},
	}
	for ix, tc := range tcs {
		t.Run(fmt.Sprintf("%d", ix), func(t *testing.T) {
			actual, err := Strip(tc.fields, tc.input)
			if err != nil {
				t.Fatalf("call to StripJson failed unexpectedly: %v", err)
			}

			actualStr, err := json.Marshal(actual)
			if err != nil {
				t.Fatalf("could not convert actual into bytes: %v", err)
			}

			expectedStr, err := json.Marshal(tc.expected)
			if err != nil {
				t.Fatalf("could not convert expected into bytes: %v", err)
			}

			if string(actualStr) != string(expectedStr) {
				t.Fatalf("actual does not equal expected\nexpected:\n\t%s\nactual:\n\t%s\n", string(expectedStr), string(actualStr))
			}
		})
	}
}
