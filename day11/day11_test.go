package main

import "testing"

func TestSimulatePartOne(t *testing.T) {
	input := make([]monkey, len(msample))
	copy(input, msample)
	if result := simulate(input, 20, true); result != 10605 {
		t.Errorf("For sample input, after 20 rounds with relief expected 10605 but got %d", result)
	}
}

func TestSimulatePartTwo(t *testing.T) {
	input := make([]monkey, len(msample))
	copy(input, msample)
	if result := simulate(input, 10000, false); result != 2713310158 {
		t.Errorf("For sample input, after 10000 rounds with relief expected 2713310158 but got %d", result)
	}
}

var msample []monkey = []monkey{
	{
		items: []int{79, 98},
		ic:    0,
		ins:   func(w int) int { return w * 19 },
		mt:    monkeyTest{23, 2, 3},
	},
	{
		items: []int{54, 65, 75, 74},
		ic:    0,
		ins:   func(w int) int { return w + 6 },
		mt:    monkeyTest{19, 2, 0},
	},
	{
		items: []int{79, 60, 97},
		ic:    0,
		ins:   func(w int) int { return w * w },
		mt:    monkeyTest{13, 1, 3},
	},
	{
		items: []int{74},
		ic:    0,
		ins:   func(w int) int { return w + 3 },
		mt:    monkeyTest{17, 0, 1},
	},
}
