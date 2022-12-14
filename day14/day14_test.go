package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input    string
		expected []coord
	}{
		{
			input: "498,4 -> 498,6",
			expected: []coord{
				{-2, 4}, {-2, 6},
			},
		},
		{
			input: "498,4 -> 498,6 -> 496,6",
			expected: []coord{
				{-2, 4}, {-2, 6}, {-4, 6},
			},
		},
		{
			input: "503,4 -> 502,4 -> 502,9 -> 494,9",
			expected: []coord{
				{3, 4}, {2, 4}, {2, 9}, {-6, 9},
			},
		},
	}

	for _, c := range cases {
		if result := parse(c.input); !reflect.DeepEqual(result, c.expected) {
			t.Errorf("Parsing `%s` expected %v but got %v\n", c.input, c.expected, result)
		}
	}
}

func TestFillRockMap(t *testing.T) {
	cases := []struct {
		input       []coord
		expectedLen int
	}{
		{
			[]coord{{0, 0}, {4, 0}},
			5,
		},
		{
			[]coord{{0, 0}, {-4, 0}},
			5,
		},
		{
			[]coord{{0, 0}, {0, 4}},
			5,
		},
		{
			[]coord{{0, 0}, {0, -4}},
			5,
		},
		{
			[]coord{{0, 0}, {0, 2}, {2, 2}},
			5,
		},
		{
			[]coord{{0, 0}, {0, 4}, {2, 4}, {2, 2}, {-2, 2}},
			12,
		},
	}

	for _, c := range cases {
		m := make(map[coord]bool)
		fillRockMap(m, c.input)
		if len(m) != c.expectedLen {
			t.Errorf("Filling with %v, expected %d but got %d\n", c.input, c.expectedLen, len(m))
		}
	}
}
