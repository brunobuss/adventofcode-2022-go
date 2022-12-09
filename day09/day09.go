package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Coords struct {
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

	tm := map[int]map[int]bool{
		0: {
			0: true,
		},
	}
	tq := 1

	knots := make([]Coords, 10)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		d := byte(line[0])
		q, err := strconv.Atoi(line[2:])
		if err != nil {
			log.Fatalln("Error converting to int: ", line[2:])
		}

		log.Printf("Processing `%s`\n", line)

		for i := 0; i < q; i++ {
			switch d {
			case 'U':
				knots[0].y++
			case 'D':
				knots[0].y--
			case 'R':
				knots[0].x++
			case 'L':
				knots[0].x--
			default:
				log.Fatalln("Unrecognized command: ", d)
			}

			for i := 1; i < len(knots); i++ {
				knots[i] = moveTail(knots[i-1], knots[i])
			}

			log.Printf("Status (%d of %d): %v\n", i+1, q, knots)

			// Check if last knot visited this space already, if not
			// mark as visited and increase count of visited spaces
			lk := knots[len(knots)-1]
			if _, e := tm[lk.x]; !e {
				tm[lk.x] = map[int]bool{}
			}
			if !tm[lk.x][lk.y] {
				tm[lk.x][lk.y] = true
				tq++
			}
		}

		log.Printf("After `%s`, we got: %v\n", line, knots)
	}

	fmt.Println("Last knot visited spaces: ", tq)
}

func moveTail(h, t Coords) Coords {
	if (absDiff(h.x, t.x) + absDiff(h.y, t.y)) > 4 {
		log.Fatalf("Illegal State for coords H(%v) and T(%v)\n", h, t)
	}

	if absDiff(h.x, t.x) <= 1 && absDiff(h.y, t.y) <= 1 {
		// No move necessary, H still in the square around T.
		return t
	}

	if absDiff(h.x, t.x) == 0 {
		// Need to move vertically
		if h.y > t.y {
			t.y++
		} else {
			t.y--
		}
	} else if absDiff(h.y, t.y) == 0 {
		// Need to move horizontally
		if h.x > t.x {
			t.x++
		} else {
			t.x--
		}
	} else {
		// Need to move diagonally
		if h.y > t.y {
			t.y++
		} else {
			t.y--
		}

		if h.x > t.x {
			t.x++
		} else {
			t.x--
		}
	}
	return t
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
