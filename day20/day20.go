package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	log.SetOutput(io.Discard)
}

// We use this entry struct to distinguish two elements with same
// value, so we can easily query in a map the position
type entry struct {
	n   int
	app int
}

const DEC_KEY = 811589153

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	v := make([]entry, 0, 5000)
	app := make(map[int]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		e := entry{}
		n, err := fmt.Sscanf(line, "%d", &e.n)
		if n != 1 {
			log.Fatal(err)
		}
		e.app = app[e.n]
		app[e.n]++
		v = append(v, e)
	}

	partOne := make([]entry, len(v))
	copy(partOne, v)
	r1 := mix(partOne, 1, 1)
	fmt.Println("[PartOne] Sum Groove Coordinates:", r1)

	r10 := mix(v, DEC_KEY, 10)
	fmt.Println("[PartTwo] Sum Groove Coordinates:", r10)
}

func mix(v []entry, key, turns int) int {
	o := make([]entry, len(v))
	copy(o, v)
	log.Println("o:", o)
	log.Println("v:", v)

	pos := make(map[entry]int, 5000)
	for i, n := range v {
		pos[n] = i
	}

	for k := 0; k < turns; k++ {
		for i := 0; i < len(v); i++ {
			e := o[i]
			cp := pos[e]
			n := e.n * key
			log.Println(e.n, "moves:")
			log.Println(v)
			for s := 0; s < abs(n%(len(v)-1)); s++ {
				var np int
				if n > 0 {
					np = (cp + 1) % len(v)

					pos[v[cp]] = (pos[v[cp]] + 1) % len(v)
					pos[v[np]]--
					if pos[v[np]] == -1 {
						pos[v[np]] = len(v) - 1
					}

				} else {
					np = cp - 1
					if np == -1 {
						np = len(v) - 1
					}

					pos[v[cp]]--
					if pos[v[cp]] == -1 {
						pos[v[cp]] = len(v) - 1
					}
					pos[v[np]] = (pos[v[np]] + 1) % len(v)
				}
				// Swap
				v[cp], v[np] = v[np], v[cp]
				cp = np
				log.Println(v)
			}
		}
	}

	z := pos[entry{}]
	log.Println(z)
	log.Println(v[(z+1000)%len(v)], v[(z+2000)%len(v)], v[(z+3000)%len(v)])
	log.Println(
		v[(z+1000)%len(v)].n*key,
		v[(z+2000)%len(v)].n*key,
		v[(z+3000)%len(v)].n*key)
	return (v[(z+1000)%len(v)].n + v[(z+2000)%len(v)].n + v[(z+3000)%len(v)].n) * key
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return 0 - a
}
