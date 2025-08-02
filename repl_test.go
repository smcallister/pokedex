package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	// Test cases.
    cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " HELLO ",
			expected: []string{"hello"},
		},
		{
			input:    "     ",
			expected: []string{},
		},
	}

	// Run through each test case.
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected %d words but got %d", len(c.expected), len(actual))
			t.Fail()
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected %s but got %s", expectedWord, word)
				t.Fail()
			}
		}
	}
}