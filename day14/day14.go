package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x, y int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rocks := make(map[coord]bool)
	for scanner.Scan() {
		line := scanner.Text()
		cs := parse(line)
		fillRockMap(rocks, cs)
	}

	sandMap := make(map[coord]bool)
	lb := getLowerBoundMap(rocks)
	countSand(coord{}, sandMap, rocks, lb)
	fmt.Println("Sands count (Part One):", len(sandMap))

	sandMap = make(map[coord]bool)
	bt := findHighestY(rocks) + 2
	countSandWithBottom(coord{}, sandMap, rocks, bt)
	fmt.Println("Sands count (Part Two):", len(sandMap))
}

func parse(s string) []coord {
	ps := strings.Split(s, " -> ")
	coords := []coord{}
	for _, p := range ps {
		cs := strings.Split(p, ",")
		if len(cs) != 2 {
			log.Fatalln("Expected 2 integers, but had", len(cs), "from", s)
		}
		xv, err := strconv.Atoi(cs[0])
		if err != nil {
			log.Fatalln(err)
		}
		yv, err := strconv.Atoi(cs[1])
		if err != nil {
			log.Fatalln(err)
		}

		coords = append(coords, coord{
			x: xv - 500,
			y: yv,
		})
	}
	return coords
}

func fillRockMap(m map[coord]bool, cs []coord) {
	cur := cs[0]
	for _, next := range cs[1:] {
		m[cur] = true
		for cur != next {
			if cur.x == next.x {
				if cur.y < next.y {
					cur.y++
				} else {
					cur.y--
				}
			} else {
				if cur.x < next.x {
					cur.x++
				} else {
					cur.x--
				}
			}
			m[cur] = true
		}
		cur = next
	}
}

func getLowerBoundMap(m map[coord]bool) map[int]int {
	lb := make(map[int]int)
	for k := range m {
		if k.y > lb[k.x] {
			lb[k.x] = k.y
		}
	}
	return lb
}

func findHighestY(m map[coord]bool) int {
	h := 0
	for k := range m {
		if k.y > h {
			h = k.y
		}
	}
	return h
}

func countSand(sp coord, sm map[coord]bool, rks map[coord]bool, lb map[int]int) bool {
	if rks[sp] || sm[sp] {
		return true
	} else if sp.y > lb[sp.x] {
		return false
	}

	below := coord{sp.x, sp.y + 1}
	if r := countSand(below, sm, rks, lb); !r {
		return false
	}
	diagLeft := coord{sp.x - 1, sp.y + 1}
	if r := countSand(diagLeft, sm, rks, lb); !r {
		return false
	}
	diagRight := coord{sp.x + 1, sp.y + 1}
	if r := countSand(diagRight, sm, rks, lb); !r {
		return false
	}

	sm[sp] = true
	return true
}

func countSandWithBottom(sp coord, sm map[coord]bool, rks map[coord]bool, bt int) {
	if rks[sp] || sm[sp] || sp.y >= bt {
		return
	}

	below := coord{sp.x, sp.y + 1}
	countSandWithBottom(below, sm, rks, bt)
	diagLeft := coord{sp.x - 1, sp.y + 1}
	countSandWithBottom(diagLeft, sm, rks, bt)
	diagRight := coord{sp.x + 1, sp.y + 1}
	countSandWithBottom(diagRight, sm, rks, bt)

	sm[sp] = true
}
