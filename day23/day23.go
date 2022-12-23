package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func init() {
	log.SetOutput(io.Discard)
}

type coord struct {
	x, y int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	em := make(map[coord]bool)
	ln := 0
	for scanner.Scan() {
		line := scanner.Text()
		for i, r := range line {
			if r == '#' {
				em[coord{i, ln}] = true
			}
		}
		ln++
	}

	log.Println("em", len(em))
	printMap(em)

	r := 0
	partOneAnswer := 0
	for {
		moves := make(map[coord]coord)
		movesC := make(map[coord]int)

		for e := range em {
			npos := e
			if wantToMove(e, em) {
				mo := findNextMoveOffset(r, e, em)
				npos = coord{e.x + mo.x, e.y + mo.y}
			}
			moves[e] = npos
			movesC[npos]++
		}

		log.Println("== Begin of Round", r+1, "==")
		log.Println("moves", len(moves))
		printMoveMap(em, moves)
		log.Println("movesC", len(movesC))

		em = make(map[coord]bool)
		mc := 0
		for e, move := range moves {
			if movesC[move] == 1 {
				// Only one elf wants to move there, so we can move
				em[move] = true
				if e != move {
					mc++
				}
			} else {
				// More than one, so we need to hold back
				em[e] = true
			}
		}

		log.Println("em", len(em))
		log.Println("mc", mc)
		log.Println("== End of Round", r+1, "==")
		printMap(em)

		if mc == 0 {
			// No elves moved, so we're done!
			break
		}

		if r == 9 {
			min, max := findMinMaxBoundaries(em)
			log.Println("Rect:", min, max)
			rectArea := (max.x - min.x + 1) * (max.y - min.y + 1)
			partOneAnswer = rectArea - len(em)
		}
		r++
	}

	fmt.Println("[PartOne] R10 empty tiles in rect area:", partOneAnswer)
	fmt.Println("[PartTwo] Stopped moving after round", r+1)
}

// Return true iff there is another elf touching them
func wantToMove(e coord, em map[coord]bool) bool {
	return em[coord{e.x - 1, e.y - 1}] || em[coord{e.x, e.y - 1}] || em[coord{e.x + 1, e.y - 1}] ||
		em[coord{e.x - 1, e.y + 1}] || em[coord{e.x, e.y + 1}] || em[coord{e.x + 1, e.y + 1}] ||
		em[coord{e.x - 1, e.y}] || em[coord{e.x + 1, e.y}]
}

func findNextMoveOffset(r int, e coord, em map[coord]bool) coord {
	for d := 0; d < 4; d++ {
		rd := (d + r) % 4
		switch rd {
		case 0:
			if canMoveNorth(e, em) {
				return coord{0, -1}
			}
		case 1:
			if canMoveSouth(e, em) {
				return coord{0, 1}
			}
		case 2:
			if canMoveWest(e, em) {
				return coord{-1, 0}
			}
		case 3:
			if canMoveEast(e, em) {
				return coord{1, 0}
			}
		}
	}
	return coord{}
}

func canMoveNorth(e coord, em map[coord]bool) bool {
	return !em[coord{e.x - 1, e.y - 1}] && !em[coord{e.x, e.y - 1}] && !em[coord{e.x + 1, e.y - 1}]
}

func canMoveSouth(e coord, em map[coord]bool) bool {
	return !em[coord{e.x - 1, e.y + 1}] && !em[coord{e.x, e.y + 1}] && !em[coord{e.x + 1, e.y + 1}]
}

func canMoveWest(e coord, em map[coord]bool) bool {
	return !em[coord{e.x - 1, e.y - 1}] && !em[coord{e.x - 1, e.y}] && !em[coord{e.x - 1, e.y + 1}]
}

func canMoveEast(e coord, em map[coord]bool) bool {
	return !em[coord{e.x + 1, e.y - 1}] && !em[coord{e.x + 1, e.y}] && !em[coord{e.x + 1, e.y + 1}]
}

func findMinMaxBoundaries(em map[coord]bool) (coord, coord) {
	minC, maxC := coord{math.MaxInt, math.MaxInt}, coord{math.MinInt, math.MinInt}
	for e := range em {
		minC.x = min(minC.x, e.x)
		minC.y = min(minC.y, e.y)

		maxC.x = max(maxC.x, e.x)
		maxC.y = max(maxC.y, e.y)
	}
	return minC, maxC
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func printMap(em map[coord]bool) {
	min, max := findMinMaxBoundaries(em)
	for y := min.y - 2; y < max.y+2; y++ {
		s := ""
		for x := min.x - 2; x < max.x+2; x++ {
			if em[coord{x, y}] {
				s += "#"
			} else {
				s += "."
			}
		}
		log.Println(s)
	}
}

func printMoveMap(em map[coord]bool, moves map[coord]coord) {
	min, max := findMinMaxBoundaries(em)
	for y := min.y - 2; y < max.y+2; y++ {
		s := ""
		for x := min.x - 2; x < max.x+2; x++ {
			e := coord{x, y}
			if !em[e] {
				s += "."
				continue
			}
			if moves[e].y > e.y {
				s += "V"
			} else if moves[e].y < e.y {
				s += "^"
			} else if moves[e].x > e.x {
				s += ">"
			} else if moves[e].x < e.x {
				s += "<"
			} else {
				s += "#"
			}
		}
		log.Println(s)
	}
}
