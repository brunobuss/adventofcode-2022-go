package main

import "testing"

func TestFindStartOfPacketOffset(t *testing.T) {
	c := []struct {
		input    string
		expected int
	}{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 7},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 5},
		{"nppdvjthqldpwncqszvftbrmjlhg", 6},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 10},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 11},
	}

	for _, test := range c {
		if result := findMarkerOffset([]byte(test.input), 4); result != test.expected {
			t.Errorf("For buffer %s expected start offset is %d but got %d",
				test.input, test.expected, result)
		}
	}
}

func TestFindStartOfMessageOffset(t *testing.T) {
	c := []struct {
		input    string
		expected int
	}{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 19},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 23},
		{"nppdvjthqldpwncqszvftbrmjlhg", 23},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 29},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 26},
	}

	for _, test := range c {
		if result := findMarkerOffset([]byte(test.input), 14); result != test.expected {
			t.Errorf("For buffer %s expected start offset is %d but got %d",
				test.input, test.expected, result)
		}
	}
}
