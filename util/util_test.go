package util

import (
	"testing"
)

func TestValidChecker(t *testing.T) {
	var TestCases = []struct {
		input    string
		expected string
	}{
		{"01234567", "permanent"},
		{"12345678", "permanent"},
		{"asdfqwer", "temporary"},
		{"0asdf123", "temporary"},
		{"0", ""},
		{"a", ""},
		{"asdf", "temporary"},
		{"asd", "temporary"},
		{"asdfqwerasdf", ""},
		{"1000000000", ""},
		{"dafsdf?", ""},
		{"once", "temporary"},
	}
	for i, TestCase := range TestCases {
		output, err := ValidChecker(TestCase.input)
		if output != TestCase.expected {
			t.Fatalf("Test %d | input: %s, expected: %s, output: %s\n", i, TestCase.input, TestCase.expected, output)
		}
		if err != nil {
			t.Logf("[%s: %s]\n", err, TestCase.input)
		}
	}
}
