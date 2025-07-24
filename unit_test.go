package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{input: "Hello World",
			expected: []string{"Hello", "World"},
		},
		{
			input:    "Hello World wowowowwoow",
			expected: []string{"Hello", "World", "wowowowwoow"},
		},
	}
	for _, cas := range cases {
		output := cleanInput(cas.input)
		if len(output) != len(cas.expected) {
			t.Errorf("Test case failed, Length does not match")
		}
		for i := range output {
			wordactual := cas.expected[i]
			wordgot := output[i]
			if wordgot != wordactual {
				t.Errorf("test case failed words dont match")
			}
		}
	}

}
