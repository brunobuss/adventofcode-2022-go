package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dp := 0
	bp := 0
	for {
		runsack1, runsack2, runsack3, s := readGroup(scanner)
		if !s {
			break
		}
		dp += diffPriority(runsack1) + diffPriority(runsack2) + diffPriority(runsack3)
		bp += badgePriority(runsack1, runsack2, runsack3)
	}
	fmt.Println("Diff Priority Sum: ", dp)
	fmt.Println("Badge Priority Sum: ", bp)
}

// Read groups of 3 runsacks
func readGroup(scanner *bufio.Scanner) (string, string, string, bool) {
	var runsacks []string
	for i := 0; i < 3; i++ {
		if !scanner.Scan() {
			return "", "", "", false
		}
		line := scanner.Text()
		var r string
		n, err := fmt.Sscan(line, &r)
		if n != 1 {
			log.Fatal("Error reading line: ", err)
		}
		runsacks = append(runsacks, r)
	}
	return runsacks[0], runsacks[1], runsacks[2], true
}

func priority(item byte) byte {
	if item >= 'a' && item <= 'z' {
		return (item - 'a') + 1
	} else if item >= 'A' && item <= 'Z' {
		return (item - 'A') + 27
	}
	log.Fatalln("Invalid item: ", item)
	return 0
}

func diffPriority(runsack string) int {
	h := len(runsack) / 2
	var cOne map[byte]int = map[byte]int{}
	for _, v := range runsack[:h] {
		cOne[byte(v)]++
	}
	for _, v := range runsack[h:] {
		if cOne[byte(v)] != 0 {
			return int(priority(byte(v)))
		}
	}
	log.Fatalln("Duplicated item not found: ", runsack)
	return 0
}

func badgePriority(runsack1, runsack2, runsack3 string) int {
	var cOne map[byte]int = map[byte]int{}
	for _, v := range runsack1 {
		cOne[byte(v)]++
	}
	var cTwo map[byte]int = map[byte]int{}
	for _, v := range runsack2 {
		cTwo[byte(v)]++
	}
	for _, v := range runsack3 {
		if cOne[byte(v)] != 0 && cTwo[byte(v)] != 0 {
			return int(priority(byte(v)))
		}
	}
	log.Fatalln("Badge Not Found: ", cOne, cTwo, runsack3)
	return 0
}
