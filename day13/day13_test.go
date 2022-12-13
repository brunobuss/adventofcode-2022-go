package main

import (
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	cases := []struct {
		l, r     string
		expected int
	}{
		{
			"[1,1,3,1,1]",
			"[1,1,5,1,1]",
			-1,
		},
		{
			"[[1],[2,3,4]]",
			"[[1],4]",
			-1,
		},
		{
			"[9]",
			"[[8,7,6]]",
			1,
		},
		{
			"[[4,4],4,4]",
			"[[4,4],4,4,4]",
			-1,
		},
		{
			"[7,7,7,7]",
			"[7,7,7]",
			1,
		},
		{
			"[]",
			"[3]",
			-1,
		},
		{
			"[3]",
			"[3]",
			0,
		},
		{
			"[]",
			"[]",
			0,
		},
		{
			"[[[]]]",
			"[[]]",
			1,
		},
		{
			"[[]]",
			"[[[]]]",
			-1,
		},
		{
			"[1,[2,[3,[4,[5,6,7]]]],8,9]",
			"[1,[2,[3,[4,[5,6,0]]]],8,9]",
			1,
		},
	}

	for _, c := range cases {
		le, _ := parse(c.l)
		re, _ := parse(c.r)
		if result := compare(*le, *re); result != c.expected {
			t.Errorf("Expected %v but got %v for:\nleft  = %v\nright = %v\n", c.expected, result, c.l, c.r)
		}
	}

}

func TestParse(t *testing.T) {
	cases := []struct {
		input    string
		expected packet
	}{
		{
			"13",
			newValueEntry(13),
		},
		{
			"[]",
			newEmptyListEntry(),
		},
		{
			"[[]]",
			newListEntry([]packet{
				newEmptyListEntry(),
			}),
		},
		{
			"[13]",
			newListEntry([]packet{
				newValueEntry(13),
			}),
		},
		{
			"[13,5,33]",
			newListEntry([]packet{
				newValueEntry(13),
				newValueEntry(5),
				newValueEntry(33),
			}),
		},
		{
			"[[13]]",
			newListEntry([]packet{
				newListEntry([]packet{
					newValueEntry(13),
				}),
			}),
		},
		{
			"[[13,5,33]]",
			newListEntry([]packet{
				newListEntry([]packet{
					newValueEntry(13),
					newValueEntry(5),
					newValueEntry(33),
				}),
			}),
		},
		{
			"[[13],[5],[33]]",
			newListEntry([]packet{
				newListEntry([]packet{
					newValueEntry(13),
				}),
				newListEntry([]packet{
					newValueEntry(5),
				}),
				newListEntry([]packet{
					newValueEntry(33),
				}),
			}),
		},
		{
			"[[13],5,[33,[]],[99],8]",
			newListEntry([]packet{
				newListEntry([]packet{
					newValueEntry(13),
				}),
				newValueEntry(5),
				newListEntry([]packet{
					newValueEntry(33),
					newEmptyListEntry(),
				}),
				newListEntry([]packet{
					newValueEntry(99),
				}),
				newValueEntry(8),
			}),
		},
	}

	for _, c := range cases {
		if result, _ := parse(c.input); reflect.DeepEqual(result, c.expected) {
			t.Errorf("From `%s` expected %v but got %v", c.input, c.expected, result)
		}
	}
}
