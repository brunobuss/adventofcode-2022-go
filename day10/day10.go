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
		log.Fatalln(err)
	}
	defer file.Close()

	x := 1
	c := 1
	s := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		il := 1  // Noop only takes 1 cycle
		inc := 0 // Noop does not increment x

		if line[0] == 'a' {
			// ADDX
			n, err := fmt.Sscanf(line, "addx %d", &inc)
			if n != 1 {
				log.Fatalf("Error parsing `%s`: %s\n", line, err)
			}
			il = 2
		}

		for i := 0; i < il; i++ {
			if c == 20 || (c-20)%40 == 0 {
				log.Printf("During cycle %d, X=%d.\n", c, x)
				s += c * x
			}

			// We use (c-1) because the pixel is supposed to be 0-based
			// when comparing with the register X, but our logic here
			// starts with c = 1 to match the challenge statement.
			if absDiff(x, (c-1)%40) <= 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
			if c%40 == 0 {
				fmt.Println()
			}
			c += 1
		}
		x += inc
	}
	fmt.Println("Sum of signal strength: ", s)
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
