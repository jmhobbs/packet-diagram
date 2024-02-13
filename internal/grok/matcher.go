package grok

import "regexp"

var ruleMatcher *regexp.Regexp

func init() {
	/*
		Group 1: Negation for N bytes
		Group 2: N bytes
		Group 3: Read/Skip marker for scan
		Group 4: Target byte for scan
	*/
	ruleMatcher = regexp.MustCompile("(-)?([0-9]+)|(->|>)([0-9a-fA-F]{1,2})")
}

type matcher interface {
	Match([]byte) (matched []byte, length int, err error)
}
