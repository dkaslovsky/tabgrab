package main

var commonPrefixes = newPrefixSet("-", "*")

type prefixSet map[byte]struct{}

func newPrefixSet(prefixes ...string) prefixSet {
	set := make(prefixSet)
	for _, p := range prefixes {
		set.addFrom(p)
	}
	return set
}

func (p prefixSet) addFrom(s string) {
	if len(s) == 0 {
		return
	}
	p[s[0]] = struct{}{}
}

func (p prefixSet) contains(s string) bool {
	if len(s) == 0 {
		return false
	}
	_, found := p[s[0]]
	return found
}

func (p prefixSet) pop() (string, bool) {
	for item := range p {
		return string(item), true
	}
	return "", false
}

func checkPrefixMismatch(prefixes prefixSet, targetPrefix string) bool {
	// No mismatch without a single unique prefix in the set
	if len(prefixes) != 1 {
		return false
	}

	uniqPrefix, exists := prefixes.pop()
	if !exists {
		return false // Do not return a mismatch if a prefix cannot be popped from the set
	}

	// Mismatch if the target prefix is empty and the set's unique prefix is a common prefix
	if targetPrefix == "" {
		return commonPrefixes.contains(uniqPrefix)
	}

	return uniqPrefix[0] != targetPrefix[0]
}
