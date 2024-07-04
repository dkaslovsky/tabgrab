package main

import "testing"

func TestCheckPrefixMismatch(t *testing.T) {
	tests := map[string]struct {
		prefixes     prefixSet
		targetPrefix string
		expected     bool
	}{
		"empty prefix set": {
			prefixes:     newPrefixSet(),
			targetPrefix: "",
			expected:     false,
		},
		"multiples prefixes in set": {
			prefixes:     newPrefixSet("-", "*"),
			targetPrefix: "",
			expected:     false,
		},
		"unique prefix in set from common prefix set and target matches": {
			prefixes:     newPrefixSet("-"),
			targetPrefix: "-",
			expected:     false,
		},
		"unique prefix in set from common prefix set and target starts with a match": {
			prefixes:     newPrefixSet("-"),
			targetPrefix: "- ",
			expected:     false,
		},
		"unique prefix in set from common prefix set and target does not match": {
			prefixes:     newPrefixSet("-"),
			targetPrefix: "*",
			expected:     true,
		},
		"unique prefix in set from common prefix set and empty target does not match": {
			prefixes:     newPrefixSet("-"),
			targetPrefix: "",
			expected:     true,
		},
		"unique prefix in set not in common prefix set and empty target does not match": {
			prefixes:     newPrefixSet("h"),
			targetPrefix: "",
			expected:     false,
		},
		"unique prefix in set not in common prefix set and target does not match": {
			prefixes:     newPrefixSet("h"),
			targetPrefix: "-",
			expected:     true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if result := checkPrefixMismatch(test.prefixes, test.targetPrefix); result != test.expected {
				t.Errorf("expected %t, result %t", test.expected, result)
			}
		})
	}
}
