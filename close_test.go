package main

import "testing"

func TestMatch(tt *testing.T) {
	tests := map[string]struct {
		url          string
		matchVals    []string
		nonMatchVals []string
		expected     bool
	}{
		"empty matchVals and nonMatchVals": {
			url:          "foo.bar.baz",
			matchVals:    []string{},
			nonMatchVals: []string{},
			expected:     false,
		},
		"empty string matchVals and nonMatchVals": {
			url:          "foo.bar.baz",
			matchVals:    []string{""},
			nonMatchVals: []string{""},
			expected:     false,
		},
		"single matchVal that matches": {
			url:          "foo.bar.baz",
			matchVals:    []string{"bar"},
			nonMatchVals: []string{},
			expected:     true,
		},
		"single matchVal that does not match": {
			url:          "foo.bar.baz",
			matchVals:    []string{"xyz"},
			nonMatchVals: []string{},
			expected:     false,
		},
		"single matchVal that matches with empty nonMatchVal": {
			url:          "foo.bar.baz",
			matchVals:    []string{"bar"},
			nonMatchVals: []string{""},
			expected:     true,
		},
		"multiple matchVals that match": {
			url:          "foo.bar.baz",
			matchVals:    []string{"bar", "fo"},
			nonMatchVals: []string{},
			expected:     true,
		},
		"multiple matchVals with only one matching": {
			url:          "foo.bar.baz",
			matchVals:    []string{"bar", "xyz"},
			nonMatchVals: []string{},
			expected:     false,
		},
		"multiple matchVals with only one matching and order reversed": {
			url:          "foo.bar.baz",
			matchVals:    []string{"xyz", "bar"},
			nonMatchVals: []string{},
			expected:     false,
		},
		"single nonMatchVal that matches": {
			url:          "foo.bar.baz",
			matchVals:    []string{},
			nonMatchVals: []string{"bar"},
			expected:     false,
		},
		"single nonMatchVal that does not match": {
			url:          "foo.bar.baz",
			matchVals:    []string{},
			nonMatchVals: []string{"xyz"},
			expected:     true,
		},
		"single nonMatchVal that does not match with empty matchVal": {
			url:          "foo.bar.baz",
			matchVals:    []string{""},
			nonMatchVals: []string{"xyz"},
			expected:     true,
		},
		"multiple nonMatchVals that match": {
			url:          "foo.bar.baz",
			matchVals:    []string{},
			nonMatchVals: []string{"bar", "fo"},
			expected:     false,
		},
		"multiple nonMatchVals that do not match": {
			url:          "foo.bar.baz",
			matchVals:    []string{},
			nonMatchVals: []string{"xyz", "aaa"},
			expected:     true,
		},
		"multiple nonMatchVals with only one matching": {
			url:          "foo.bar.baz",
			matchVals:    []string{},
			nonMatchVals: []string{"xyz", "bar"},
			expected:     false,
		},
		"multiple nonMatchVals with only one matching and order reversed": {
			url:          "foo.bar.baz",
			matchVals:    []string{},
			nonMatchVals: []string{"bar", "xyz"},
			expected:     false,
		},
		"matching matchVal with matching nonMatchVal": {
			url:          "foo.bar.baz",
			matchVals:    []string{"bar"},
			nonMatchVals: []string{"foo"},
			expected:     false,
		},
		"matching matchVal with non-matching nonMatchVal": {
			url:          "foo.bar.baz",
			matchVals:    []string{"bar"},
			nonMatchVals: []string{"xyz"},
			expected:     true,
		},
		"same values for matchVal and nonMatchVal and both match": {
			url:          "foo.bar.baz",
			matchVals:    []string{"bar"},
			nonMatchVals: []string{"bar"},
			expected:     false,
		},
		"same values for matchVal and nonMatchVal and neither match": {
			url:          "foo.bar.baz",
			matchVals:    []string{"xyz"},
			nonMatchVals: []string{"xyz"},
			expected:     false,
		},
	}

	for name, test := range tests {
		tt.Run(name, func(t *testing.T) {
			result := match(test.url, test.matchVals, test.nonMatchVals)
			if result != test.expected {
				t.Errorf("expected %t, result %t", test.expected, result)
			}
		})
	}
}
