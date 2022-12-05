package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

// 3 for sample, 9 for input
const STACKS = 9

// This could have been a cmd line flag :)
const PART_TWO = true

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stacks := make([][]byte, STACKS)
	scanner := bufio.NewScanner(file)
	// First we read the containers/stacks
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		// This input might look a bit like a PITA to parse, but because it's
		// nicely alligned we know the offset where each letter should suppose
		// to be: first stack is in 1, second in 5, third in 9, ...
		offset := 1
		stack := 0
		for offset < len(line) {
			if unicode.IsDigit(rune(line[offset])) {
				// This is the line with the stacks numbers, so we can skip it
				break
			}

			item := line[offset]
			if item != ' ' {
				stacks[stack] = append(stacks[stack], item)
			}
			stack++
			offset += 4
		}
	}

	// In the step before, we built the stacks with the "top" at the
	// beginning of the lists, so lets flip them all so we can use them
	// like stacks
	for i := range stacks {
		reverse(stacks[i])
	}

	// Now we replicate the moves
	for scanner.Scan() {
		line := scanner.Text()

		var q, from, to int
		n, err := fmt.Sscanf(line, "move %d from %d to %d", &q, &from, &to)
		if n != 3 {
			log.Fatal("Error reading line: ", err)
		}

		// Make the indexes 0-based:
		from--
		to--

		// Instead of moving one by one, we can just reverse the top q elements
		// in the from stack and just copy them all to the new stack
		s, e := len(stacks[from])-q, len(stacks[from])
		cp := stacks[from][s:e]
		if !PART_TWO {
			reverse(cp)
		}
		stacks[to] = append(stacks[to], cp...)
		stacks[from] = stacks[from][:s]
	}

	fmt.Print("Top of Stacks: ")
	for _, v := range stacks {
		fmt.Printf("%c", v[len(v)-1])
	}
	fmt.Println()
}

func reverse(stack []byte) {
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}
}
