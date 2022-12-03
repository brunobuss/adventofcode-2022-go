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
		rucksack1, rucksack2, rucksack3, s := readGroup(scanner)
		if !s {
			break
		}
		dp += diffPriority(rucksack1) + diffPriority(rucksack2) + diffPriority(rucksack3)
		bp += badgePriority(rucksack1, rucksack2, rucksack3)
	}
	fmt.Println("Diff Priority Sum: ", dp)
	fmt.Println("Badge Priority Sum: ", bp)
}

// Read groups of 3 rucksack
func readGroup(scanner *bufio.Scanner) (string, string, string, bool) {
	var rucksacks []string
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
		rucksacks = append(rucksacks, r)
	}
	return rucksacks[0], rucksacks[1], rucksacks[2], true
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

func diffPriority(rucksack string) int {
	h := len(rucksack) / 2
	var cOne map[byte]int = map[byte]int{}
	for _, v := range rucksack[:h] {
		cOne[byte(v)]++
	}
	for _, v := range rucksack[h:] {
		if cOne[byte(v)] != 0 {
			return int(priority(byte(v)))
		}
	}
	log.Fatalln("Duplicated item not found: ", rucksack)
	return 0
}

func badgePriority(rucksack1, rucksack2, rucksack3 string) int {
	var bc1 map[byte]int = map[byte]int{}
	for _, v := range rucksack1 {
		bc1[byte(v)]++
	}
	var bcBoth map[byte]int = map[byte]int{}
	for _, v := range rucksack2 {
		// Items not present in the first set, cannot be badge candidates
		// so we only add ones also present there.
		if bc1[byte(v)] != 0 {
			bcBoth[byte(v)]++
		}
	}
	for _, v := range rucksack3 {
		if bcBoth[byte(v)] != 0 {
			return int(priority(byte(v)))
		}
	}
	log.Fatalln("Badge Not Found: ", bc1, bcBoth, rucksack3)
	return 0
}
