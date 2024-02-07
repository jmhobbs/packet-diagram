package main

import "testing"

func Test_ParseRuleLine(t *testing.T) {
	tests := []struct {
		input               string
		expectedDescription string
		expectedSkip        bool
		expectedMatcher     matcher
	}{
		{
			input:               "5 Read Five Bytes",
			expectedDescription: "Read Five Bytes",
			expectedSkip:        false,
			expectedMatcher:     readBytes{5},
		},
		{
			input:               "-5 Skip Five Bytes",
			expectedDescription: "Skip Five Bytes",
			expectedSkip:        true,
			expectedMatcher:     readBytes{5},
		},
		{
			input:               ">0a Read Until 0a",
			expectedDescription: "Read Until 0a",
			expectedSkip:        false,
			expectedMatcher:     readUntil{0xa},
		},
		{
			input:               "->0a Skip Until 0a",
			expectedDescription: "Skip Until 0a",
			expectedSkip:        true,
			expectedMatcher:     readUntil{0xa},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			r, err := parseRuleLine(test.input)
			if err != nil {
				t.Error(err)
				t.Fail()
			}

			if r.Description != test.expectedDescription {
				t.Errorf("Expected %q, got %q", test.expectedDescription, r.Description)
			}

			if r.Skip != test.expectedSkip {
				t.Errorf("Expected %t, got %t", test.expectedSkip, r.Skip)
			}

			if r.Matcher != test.expectedMatcher {
				t.Errorf("Expected %v, got %v", test.expectedMatcher, r.Matcher)
			}
		})
	}
}
