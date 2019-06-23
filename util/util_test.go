package util

import (
	"testing"
)

func TestValidChecker(t *testing.T) {
	var TestCases = []struct {
		input    string
		expected string
	}{
		{"01234567", ""},
		{"12345678", "permanent"},
		{"asdfqwer", "temporary"},
		{"0asdf123", "temporary"},
		{"0", ""},
		{"a", ""},
		{"asdf", "temporary"},
		{"asd", "temporary"},
		{"asdfqwerasdf", ""},
		{"1000000000", ""},
	}
	for i, TestCase := range TestCases {
		output, err := ValidChecker(TestCase.input)
		if output != TestCase.expected {
			t.Fatalf("Test %d | input: %s, expected: %s, output: %s\n", i, TestCase.input, TestCase.expected, output)
		}
		if err != nil {
			t.Logf("[%s]\n", err)
		}
	}
}

func Test_generator(t *testing.T) {
	for i := 0; i < 8; i++ {
		str, err := generator(8)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	}
}
