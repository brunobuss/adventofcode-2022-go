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

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	grid := make([][]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, 0, len(grid))
		for _, r := range line {
			row = append(row, int(r-'0'))
		}
		grid = append(grid, row)
	}

	l := len(grid)
	nvc := 0
	hs := 0
	for i := 1; i < l-1; i++ {
		for j := 1; j < l-1; j++ {
			s := 1

			vup := true
			for k := i - 1; k >= 0 && vup; k-- {
				if grid[k][j] >= grid[i][j] {
					vup = false
					s *= absDiff(i, k)
				}
			}
			if vup {
				s *= absDiff(i, 0)
			}

			vdown := true
			for k := i + 1; k < l && vdown; k++ {
				if grid[k][j] >= grid[i][j] {
					vdown = false
					s *= absDiff(i, k)
				}
			}
			if vdown {
				s *= absDiff(i, l-1)
			}

			vleft := true
			for k := j - 1; k >= 0 && vleft; k-- {
				if grid[i][k] >= grid[i][j] {
					vleft = false
					s *= absDiff(j, k)
				}
			}
			if vleft {
				s *= absDiff(j, 0)
			}

			vright := true
			for k := j + 1; k < l && vright; k++ {
				if grid[i][k] >= grid[i][j] {
					vright = false
					s *= absDiff(j, k)
				}
			}
			if vright {
				s *= absDiff(j, l-1)
			}

			// Tree is not visible only if they are not visible
			// from all 4 directions
			nv := !vup && !vdown && !vleft && !vright
			log.Printf("g(%d, %d) = %d; v(%d, %d) = %v\n",
				i, j, grid[i][j], i, j, nv)
			if nv {
				nvc++
			}

			log.Printf("S(%d, %d) = %d\n", i, j, s)
			if s > hs {
				log.Printf("New HS (%d, %d) = %d\n", i, j, s)
				hs = s
			}
		}
	}

	vc := l*l - nvc
	fmt.Println("Visible from outside grid: ", vc)
	fmt.Println("Highest Scenic Score: ", hs)
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
