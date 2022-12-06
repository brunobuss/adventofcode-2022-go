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
	for scanner.Scan() {
		line := scanner.Text()

		s := findMarkerOffset([]byte(line), 4)
		m := findMarkerOffset([]byte(line), 14)
		fmt.Println("Start of Packet offset: ", s)
		fmt.Println("Start of Message offset: ", m)
	}
}

func findMarkerOffset(buffer []byte, diff int) int {
	cm := map[byte]int{}
	for i, v := range buffer {
		cm[v]++
		if i < diff-1 {
			// Not enough characters to check for the marker yet
			continue
		} else if i >= diff {
			// We have already `diff + 1` entries added to the map, so we need to
			// remove the oldest one. If there is only one entry for that character
			// we delete the whole map entry, so the len() check below works
			// as intended. Otherwise we only reduce the count
			o := buffer[i-diff]
			if cm[o] == 1 {
				delete(cm, o)
			} else {
				cm[o]--
			}
		}

		if len(cm) == diff {
			// We have `diff`` different entries, which means we have `diff`
			// different characters so we have found our marker. Return the
			// next index as the buffer offset past marker
			return i + 1
		}
	}

	return -1
}
