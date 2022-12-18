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

type cube struct {
	x, y, z int
}

const SCANMAP = false

var offsets []cube = []cube{
	{-1, 0, 0},
	{1, 0, 0},
	{0, -1, 0},
	{0, 1, 0},
	{0, 0, -1},
	{0, 0, 1},
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cm := make(map[cube]bool)
	sfc := 0

	minC := cube{math.MaxInt16, math.MaxInt16, math.MaxInt16}
	maxC := cube{math.MinInt16, math.MinInt16, math.MinInt16}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		c := cube{}
		n, err := fmt.Sscanf(line, "%d,%d,%d", &c.x, &c.y, &c.z)
		if n != 3 {
			log.Fatal(err)
		}

		if c.x < minC.x {
			minC.x = c.x
		}
		if c.x > maxC.x {
			maxC.x = c.x
		}
		if c.y < minC.y {
			minC.y = c.y
		}
		if c.y > maxC.y {
			maxC.y = c.y
		}
		if c.z < minC.z {
			minC.z = c.z
		}
		if c.z > maxC.z {
			maxC.z = c.z
		}

		sfc += 6
		cm[c] = true

		for _, o := range offsets {
			if cm[cube{c.x + o.x, c.y + o.y, c.z + o.z}] {
				sfc -= 2
			}
		}
	}

	fmt.Println("[PartOne] Free Surfaces:", sfc)

	ext := findExternalCubes(cm, minC, maxC)

	for i := minC.x; i <= maxC.x; i++ {
		for j := minC.y; j <= maxC.y; j++ {
			for k := minC.z; k <= maxC.z; k++ {
				tc := cube{i, j, k}

				if !cm[tc] {
					continue
				}

				for _, o := range offsets {
					c := cube{tc.x + o.x, tc.y + o.y, tc.z + o.z}
					if cm[c] {
						// tc is connected with an adjacente cube, so we already dealth
						// with this case before
						continue
					}

					outsideFace := true
					// If not, we'll shoot a straight "ray" to see if we are "inside" or
					// "outside".
					for c.x >= minC.x && c.x <= maxC.x &&
						c.y >= minC.y && c.y <= maxC.y &&
						c.z >= minC.z && c.z <= maxC.z {

						if ext[c] {
							// If we found a exterior spot, we know this face is connected
							// to the exterior and can break now
							break
						}
						if cm[c] {
							outsideFace = false
							break
						}

						c = cube{c.x + o.x, c.y + o.y, c.z + o.z}
					}
					if !outsideFace {
						sfc -= 1
					}
				}

			}
		}
	}

	fmt.Println("[PartTwo] Exterior Surfaces:", sfc)

	if SCANMAP {
		for k := minC.z; k <= maxC.z; k++ {
			for i := minC.x; i <= maxC.x; i++ {
				for j := minC.y; j <= maxC.y; j++ {
					if cm[cube{i, j, k}] {
						fmt.Print("#")
					} else {
						fmt.Print(".")
					}
				}
				fmt.Println()
			}

			fmt.Println()
			fmt.Println()
		}
	}
}

func findExternalCubes(cm map[cube]bool, minC, maxC cube) map[cube]bool {

	mark := make(map[cube]int)
	q := make([]cube, 0, 1024)
	q = append(q, minC, maxC)

	for len(q) > 0 {
		c := q[0]
		q = q[1:]

		if mark[c] == 2 {
			continue
		}
		mark[c] = 2

		for _, o := range offsets {
			t := cube{c.x + o.x, c.y + o.y, c.z + o.z}
			if c.x < minC.x || c.x > maxC.x ||
				c.y < minC.y || c.y > maxC.y ||
				c.z < minC.z || c.z > maxC.z {
				continue
			}

			if cm[t] || mark[t] > 0 {
				continue
			}
			mark[t] = 1
			q = append(q, t)
		}
	}

	exCubes := make(map[cube]bool)
	for k := range mark {
		exCubes[k] = true
	}
	log.Println("External Cubes:", len(exCubes))

	return exCubes
}
