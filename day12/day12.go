package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type pos struct {
	x, y int
}

func init() {
	log.SetOutput(io.Discard)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	grid := []string{}
	var e pos
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for i := 0; i < len(line); i++ {
			if line[i] == 'E' {
				e.x = i
				e.y = len(grid)
			}
		}
		grid = append(grid, line)
	}

	shortToS := findFromExit(e, grid, true)
	fmt.Println("Part One solution: ", shortToS)
	shortToAnyA := findFromExit(e, grid, false)
	fmt.Println("Part Two solution: ", shortToAnyA)
}

func findFromExit(e pos, g []string, onlyS bool) int {
	q := []pos{e}
	v := map[pos]bool{}
	d := map[pos]int{}

	for len(q) > 0 {
		c := q[0]
		q = q[1:]
		if v[c] {
			continue
		}
		v[c] = true

		if g[c.y][c.x] == 'S' || (g[c.y][c.x] == 'a' && !onlyS) {
			return d[c]
		}

		for _, m := range genMoves(c, g) {
			if !v[m] {
				// We can safely update here because we're doing a
				// breadth-first search, otherwise we should check if
				// the distance is actually improving or not
				d[m] = d[c] + 1
				q = append(q, m)
			}
		}
	}

	log.Fatalf("Unable to find path to start\n")
	return -1
}

func genMoves(cur pos, g []string) []pos {
	candidates := []pos{}
	moves := []pos{}

	if cur.x > 0 {
		candidates = append(candidates, pos{cur.x - 1, cur.y})
	}
	if cur.x < len(g[0])-1 {
		candidates = append(candidates, pos{cur.x + 1, cur.y})
	}
	if cur.y > 0 {
		candidates = append(candidates, pos{cur.x, cur.y - 1})
	}
	if cur.y < len(g)-1 {
		candidates = append(candidates, pos{cur.x, cur.y + 1})
	}

	for _, c := range candidates {
		// Note that we want to move *from* the candidate *into* the current position
		if canTraverse(g[c.y][c.x], g[cur.y][cur.x]) {
			moves = append(moves, c)
		}
	}

	return moves
}

func canTraverse(orig, dest byte) bool {
	return hValue(dest) <= hValue(orig)+1
}

func hValue(b byte) byte {
	switch b {
	case 'S':
		return 'a'
	case 'E':
		return 'z'
	default:
		return b
	}
}
