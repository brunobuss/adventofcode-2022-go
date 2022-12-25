package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	scanner := bufio.NewScanner(file)
	snafuTotal := "0"
	snafuTotalAsInt := 0
	for scanner.Scan() {
		line := scanner.Text()
		snafuTotal = snafuSum(line, snafuTotal)
		n := snafuToInt(line)
		snafuTotalAsInt += n
	}

	log.Println("Total in decimal:", snafuTotalAsInt)
	log.Println("Sanity Check", snafuToInt(snafuTotal))
	fmt.Println("[Part One] Total in Snafu:", snafuTotal)
}

// Convert SNAFU string s into an integer
func snafuToInt(s string) int {
	p := 1
	n := 0
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case '=':
			n += p * -2
		case '-':
			n += p * -1
		case '0':
		case '1':
			n += p
		case '2':
			n += p * 2
		default:
			log.Fatalln("Unrecognized SNAFU digit", s[i], "in", s)
		}
		p *= 5
	}
	return n
}

// Sum two SNAFU strings s1 and s2
func snafuSum(s1, s2 string) string {
	if len(s1) > len(s2) {
		s2 = strings.Repeat("0", len(s1)-len(s2)) + s2
	} else {
		s1 = strings.Repeat("0", len(s2)-len(s1)) + s1
	}

	r := ""
	c := 0
	for i := len(s1) - 1; i >= 0; i-- {
		s := snafuDigitToInt(s1[i]) + snafuDigitToInt(s2[i]) + c
		sd, nc := intToSnafuDigitAndCarry(s)
		r = sd + r
		c = nc
	}

	if c != 0 {
		sd, _ := intToSnafuDigitAndCarry(c)
		r = sd + r
	}

	return r
}

func snafuDigitToInt(c byte) int {
	switch c {
	case '=':
		return -2
	case '-':
		return -1
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	default:
		log.Fatalln("Unrecognized SNAFU digit", c)
	}
	return 0
}

// This could be generalized, but for the purposes of snafuSum(),
// we only need to handle -5 to 5 as it includes the sum of 2 -2's
// plus a -1 of carry up to the sum of 2 +2's plus a 1 of carry.
func intToSnafuDigitAndCarry(d int) (string, int) {
	switch d {
	case -5:
		return "0", -1
	case -4:
		return "1", -1
	case -3:
		return "2", -1
	case -2:
		return "=", 0
	case -1:
		return "-", 0
	case 0:
		return "0", 0
	case 1:
		return "1", 0
	case 2:
		return "2", 0
	case 3:
		return "=", 1
	case 4:
		return "-", 1
	case 5:
		return "0", 1
	default:
		log.Fatalln("Can't parse", d, "into snafu digit plus carry")
	}
	return "0", 0
}
