package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	curCal := 0
	elfCals := []int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			elfCals = append(elfCals, curCal)
			curCal = 0
		} else {
			var cal int
			n, err := fmt.Sscan(line, &cal)
			if n != 1 {
				log.Fatal(err)
			}
			curCal += cal
		}
	}
	elfCals = append(elfCals, curCal)
	sort.Sort(sort.Reverse(sort.IntSlice(elfCals)))
	fmt.Println("Highest Calories: ", elfCals[0])
	fmt.Println("Highest 3 Calories: ", (elfCals[0] + elfCals[1] + elfCals[2]))
}
