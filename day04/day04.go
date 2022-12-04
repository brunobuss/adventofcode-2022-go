package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type CampRange struct {
	start, end int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	contained := 0
	intersect := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var first, second CampRange
		n, err := fmt.Sscanf(line,
			"%d-%d,%d-%d",
			&first.start, &first.end, &second.start, &second.end)
		if n != 4 {
			log.Fatal("Error reading line: ", err)
		}

		if first.contains(second) || second.contains(first) {
			contained++
		}
		if first.intersect(second) {
			intersect++
		}
	}
	fmt.Println("Fully contained ranges: ", contained)
	fmt.Println("Intersected ranges: ", intersect)
}

func (r CampRange) contains(o CampRange) bool {
	return r.start <= o.start && r.end >= o.end
}

func (r CampRange) intersect(o CampRange) bool {
	return !((o.start < r.start && o.end < r.start) || (o.start > r.end && o.end > r.end))
}
