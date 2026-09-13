package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello World",
			expected: []string{"hello", "world"},
		},

		{
			input:    "  cHArmander Bulbasaur pikachu ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    " something where",
			expected: []string{"something", "where"},
		},
	}

	for idx, c := range cases {
		fmt.Printf("Test nr.%v", idx+1)
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			fmt.Printf(" failed 👎")
			t.Errorf("Word count of %d doesn't match expected word count of %d", len(actual), len(c.expected))

		} else {
			fmt.Printf(" lengths match, passed👌\n")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				fmt.Printf(" failed 👎")
				t.Errorf("%s doesn't equal %s", word, expectedWord)
			} else {
				fmt.Printf("%s matches %s, passed👌\n", word, expectedWord)
			}
		}
	}
}
