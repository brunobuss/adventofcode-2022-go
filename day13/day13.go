package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

func init() {
	log.SetOutput(io.Discard)
}

type packet struct {
	value int
	items []packet
}

var DEBUG_IDENT = ""

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln("Error opening file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 1
	diffSum := 0
	allPackets := []packet{}
	for {
		scanner.Scan()
		firstLine := scanner.Text()
		scanner.Scan()
		secondLine := scanner.Text()

		le, _ := parse(firstLine)
		allPackets = append(allPackets, *le)
		re, _ := parse(secondLine)
		allPackets = append(allPackets, *re)

		log.Println("== Pair", i, "==")

		r := compare(*le, *re)

		switch r {
		case -1:
			diffSum += i
		case 1:
		case 0:
			log.Fatalln("Same input entry", i)
		}

		if !scanner.Scan() {
			break
		}
		i++
	}

	log.Println("Processed", i, "checks")
	fmt.Println("Index Sum of right order packages:", diffSum)

	fd, _ := parse("[[2]]")
	sd, _ := parse("[[6]]")
	allPackets = append(allPackets, *fd, *sd)
	sort.Slice(allPackets, func(i, j int) bool {
		return compare(allPackets[i], allPackets[j]) == -1
	})

	divMul := 1
	for i, p := range allPackets {
		if p.toString() == "[[2]]" || p.toString() == "[[6]]" {
			divMul *= (i + 1)
		}
	}

	fmt.Println("Decoder key for the distress signal:", divMul)
}

func newValueEntry(v int) packet {
	return packet{value: v, items: nil}
}

func newEmptyListEntry() packet {
	return packet{value: 0, items: []packet{}}
}

func newListEntry(items []packet) packet {
	return packet{value: 0, items: items}
}

func (e packet) isValue() bool {
	return e.items == nil
}

func parse(s string) (*packet, string) {
	log.Println("Parsing `", s, "`")
	if len(s) == 0 {
		return nil, s
	} else if s[0] == ']' {
		return nil, s[1:]
	}

	if s[0] == '[' {
		ne := new(packet)
		ne.items = make([]packet, 0)
		s = s[1:]
		for {
			var se *packet
			se, s = parse(s)
			if se != nil {
				ne.items = append(ne.items, *se)
				if s[0] == ',' {
					s = s[1:]
				}
			} else {
				break
			}
		}
		return ne, s
	}

	var v int
	n, err := fmt.Sscan(s, &v)
	if n != 1 {
		log.Fatalln("Error parsing number from `", s, "`: ", err)
	}
	log.Println("Parsed ", v, " from `", s, "`")
	ve := newValueEntry(v)
	vl := len(strconv.Itoa(v))
	log.Println("Returning `", s[vl:], "` from `", s, "`")
	return &ve, s[vl:]
}

// Returns -1 if right order, 1 if wrong order, 0 same order
func compare(l, r packet) int {

	log.Println(DEBUG_IDENT, "- Compare", l.toString(), "vs", r.toString())
	DEBUG_IDENT += "  "
	defer debugUp()

	if l.isValue() && r.isValue() {
		switch {
		case l.value < r.value:
			log.Println(DEBUG_IDENT, "- Left side is smaller, so inputs are in the right order")
			return -1
		case l.value > r.value:
			log.Println(DEBUG_IDENT, "- Right side is smaller, so inputs are *not* in the right order")
			return 1
		default:
			return 0
		}
	} else if !l.isValue() && !r.isValue() {
		for i := 0; i < len(l.items); i++ {
			if i >= len(r.items) {
				// Right list short of elements
				log.Println(DEBUG_IDENT, "- Right side run out of items, so inputs are *not* in the right order")
				return 1
			}
			r := compare(l.items[i], r.items[i])
			if r != 0 {
				return r
			}
		}
		if len(l.items) == len(r.items) {
			// Whole list is the same
			return 0
		}
		// Here means that lists were equal up-to length of left one
		log.Println(DEBUG_IDENT, "- Left side run out of items, so inputs are in the right order")
		return -1
	} else if l.isValue() {
		// We know l is a value and r is a list
		nl := newListEntry([]packet{l})
		log.Println(DEBUG_IDENT, "- Mixed types; convert left to", nl.toString(), "and retry comparison")
		return compare(nl, r)
	} else {
		// We know l is a list and r is a value
		nr := newListEntry([]packet{r})
		log.Println(DEBUG_IDENT, "- Mixed types; convert right to", nr.toString(), "and retry comparison")
		return compare(l, nr)
	}
}

func (e packet) toString() string {
	if e.isValue() {
		return strconv.Itoa(e.value)
	} else if len(e.items) == 0 {
		return "[]"
	}

	s := e.items[0].toString()
	for _, it := range e.items[1:] {
		s += "," + it.toString()
	}

	return "[" + s + "]"
}

func debugUp() {
	DEBUG_IDENT = DEBUG_IDENT[0 : len(DEBUG_IDENT)-2]
}
