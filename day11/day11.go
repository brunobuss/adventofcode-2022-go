package main

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
)

type inspection func(w int) int

type monkeyTest struct {
	tv int
	td int
	fd int
}

type monkey struct {
	items []int
	ic    int
	ins   inspection
	mt    monkeyTest
}

func init() {
	log.SetOutput(io.Discard)
}

func main() {
	input := make([]monkey, len(minput))

	copy(input, minput)
	r20 := simulate(input, 20, true)
	fmt.Printf("After 20 rounds with relief: %d\n", r20)

	copy(input, minput)
	r10000 := simulate(input, 10000, false)
	fmt.Printf("After 10000 rounds without relief: %d\n", r10000)
}

func simulate(ms []monkey, rounds int, relief bool) int {
	// We multiply all monkeys test value together to create an upper
	// bound on the value of the worry. This way, when not opperating
	// with relief, we can just mod over this value since every test
	// the monkies does to decide where to throw the item will keep working
	// as intended.
	// This is effectively the Least Common Multiple, since all the test
	// values in the input are distinct prime numbers.
	lcm := 1
	for _, m := range ms {
		lcm *= m.mt.tv
	}
	log.Println("lcm: ", lcm)

	for r := 0; r < rounds; r++ {
		for i, m := range ms {
			ms[i].ic += len(m.items)
			for _, item := range m.items {
				worry := m.ins(item)

				if relief {
					worry /= 3
				} else {
					worry %= lcm
				}

				var tm int = 0
				if worry%m.mt.tv == 0 {
					tm = m.mt.td
				} else {
					tm = m.mt.fd
				}

				ms[tm].items = append(ms[tm].items, worry)
			}
			ms[i].items = []int{}
		}

		if (r+1)%1000 == 0 {
			log.Println("After round", r+1)
			for i, m := range ms {
				its := ""
				for _, w := range m.items {
					its += " " + strconv.Itoa(w)
				}
				log.Printf("Monkey %d (%d): %s\n", i, m.ic, its)
			}
		}
	}

	wl := make([]int, len(ms))
	for i, m := range ms {
		log.Printf("Monkey %d inspected items %d times.\n", i, m.ic)
		wl[i] = m.ic
	}
	sort.Sort(sort.Reverse(sort.IntSlice(wl)))
	return wl[0] * wl[1]
}

var minput []monkey = []monkey{
	{
		items: []int{53, 89, 62, 57, 74, 51, 83, 97},
		ic:    0,
		ins:   func(w int) int { return w * 3 },
		mt:    monkeyTest{13, 1, 5},
	},
	{
		items: []int{85, 94, 97, 92, 56},
		ic:    0,
		ins:   func(w int) int { return w + 2 },
		mt:    monkeyTest{19, 5, 2},
	},
	{
		items: []int{86, 82, 82},
		ic:    0,
		ins:   func(w int) int { return w + 1 },
		mt:    monkeyTest{11, 3, 4},
	},
	{
		items: []int{94, 68},
		ic:    0,
		ins:   func(w int) int { return w + 5 },
		mt:    monkeyTest{17, 7, 6},
	},

	{
		items: []int{83, 62, 74, 58, 96, 68, 85},
		ic:    0,
		ins:   func(w int) int { return w + 4 },
		mt:    monkeyTest{3, 3, 6},
	},
	{
		items: []int{50, 68, 95, 82},
		ic:    0,
		ins:   func(w int) int { return w + 8 },
		mt:    monkeyTest{7, 2, 4},
	},
	{
		items: []int{75},
		ic:    0,
		ins:   func(w int) int { return w * 7 },
		mt:    monkeyTest{5, 7, 0},
	},
	{
		items: []int{92, 52, 85, 89, 68, 82},
		ic:    0,
		ins:   func(w int) int { return w * w },
		mt:    monkeyTest{2, 0, 1},
	},
}
