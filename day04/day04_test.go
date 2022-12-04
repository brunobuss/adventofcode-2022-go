package main

import "testing"

type testCase struct {
	r1, r2   CampRange
	expected bool
}

func TestCampRangeContains(t *testing.T) {
	cases := []testCase{
		{CampRange{3, 7}, CampRange{1, 2}, false},
		{CampRange{3, 7}, CampRange{1, 3}, false},
		{CampRange{3, 7}, CampRange{1, 5}, false},

		{CampRange{3, 7}, CampRange{3, 3}, true},
		{CampRange{3, 7}, CampRange{3, 5}, true},
		{CampRange{3, 7}, CampRange{3, 7}, true},
		{CampRange{3, 7}, CampRange{5, 7}, true},
		{CampRange{3, 7}, CampRange{7, 7}, true},

		{CampRange{3, 7}, CampRange{5, 10}, false},
		{CampRange{3, 7}, CampRange{7, 10}, false},
		{CampRange{3, 7}, CampRange{9, 10}, false},
	}

	for _, test := range cases {
		if result := test.r1.contains(test.r2); result != test.expected {
			t.Errorf("Result of %v contains %v was %v, but expected %v",
				test.r1, test.r2, result, test.expected)
		}
	}
}

func TestCampRangeIntersect(t *testing.T) {
	cases := []testCase{
		{CampRange{3, 7}, CampRange{1, 2}, false},
		{CampRange{3, 7}, CampRange{1, 3}, true},
		{CampRange{3, 7}, CampRange{1, 5}, true},

		{CampRange{3, 7}, CampRange{3, 3}, true},
		{CampRange{3, 7}, CampRange{3, 5}, true},
		{CampRange{3, 7}, CampRange{3, 7}, true},
		{CampRange{3, 7}, CampRange{5, 7}, true},
		{CampRange{3, 7}, CampRange{7, 7}, true},

		{CampRange{3, 7}, CampRange{5, 10}, true},
		{CampRange{3, 7}, CampRange{7, 10}, true},
		{CampRange{3, 7}, CampRange{9, 10}, false},
	}

	for _, test := range cases {
		if result := test.r1.intersect(test.r2); result != test.expected {
			t.Errorf("Result of %v intersects %v was %v, but expected %v",
				test.r1, test.r2, result, test.expected)
		}
	}
}
