package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type monkey struct {
	v          int
	m1, m2, op string
}

func (m monkey) getValue(mm map[string]monkey) int {
	if m.op == "" {
		return m.v
	}

	log.Printf("Resolving %#v\n", m)
	m1v := mm[m.m1].getValue(mm)
	m2v := mm[m.m2].getValue(mm)
	log.Printf("m1: %d; m2: %d\n", m1v, m2v)
	switch m.op {
	case "+":
		return m1v + m2v
	case "-":
		return m1v - m2v
	case "*":
		return m1v * m2v
	case "/":
		return m1v / m2v
	default:
		log.Fatalln("Unrecognized op", m.op, "on monkey", m)
	}
	return 0
}

func init() {
	log.SetOutput(io.Discard)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mm := make(map[string]monkey)
	pm := make(map[string]string) // Map each monkey to its parent
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		m := monkey{}
		var name string

		if strings.ContainsAny(line, "+-/*") {
			n, err := fmt.Sscanf(line, "%s %s %s %s", &name, &m.m1, &m.op, &m.m2)
			if n != 4 {
				log.Fatalln("Issue scanning op monkey `", line, "`; err:", err)
			}
			// Remove trailing ":"
			name = name[:len(name)-1]
			mm[name] = m
			pm[m.m1] = name
			pm[m.m2] = name
		} else {
			n, err := fmt.Sscanf(line, "%s %d", &name, &m.v)
			if n != 2 {
				log.Fatalln("Issue scanning val monkey", line, "err:", err)
			}
			// Remove trailing ":"
			name = name[:len(name)-1]
			mm[name] = m
		}
	}

	fmt.Println("[PartOne] Answer from root:", mm["root"].getValue(mm))

	fmt.Println("[PartOne] Humn to match root:", humnToMatchRoot(mm, pm))
}

func humnToMatchRoot(mm map[string]monkey, pm map[string]string) int {
	// First we build the path from humn to root
	st := []string{}
	name := "humn"
	for {
		st = append(st, name)
		name = pm[name]
		if name == "" {
			break
		}
	}
	log.Println(st)

	// Now we start descending from root back to human, doing the operation
	// backwards in each node until we arrive at humn.
	//
	// On the first iteration of the loop below (when "root" is on top of the)
	// stack, tv will be set to the value we need to calculate
	// On each other iteration - until there is only "humn" in the stack, we
	// update tv to hold the target value of the following node.
	tv := 0
	for len(st) > 1 {
		cur := st[len(st)-1]
		next := st[len(st)-2]
		log.Printf("Cur: %s, Next: %s, TV: %d\n", cur, next, tv)
		st = st[:len(st)-1]

		var oppV int
		var or int
		if mm[cur].m1 == next {
			oppV = mm[mm[cur].m2].getValue(mm)
			or = 0
		} else {
			oppV = mm[mm[cur].m1].getValue(mm)
			or = 1
		}

		op := mm[cur].op
		if cur == "root" {
			op = "="
		}

		switch op {
		case "+":
			tv = tv - oppV
		case "-":
			if or == 0 {
				tv = tv + oppV
			} else {
				tv = oppV - tv
			}
		case "*":
			tv = tv / oppV
		case "/":
			if or == 0 {
				tv = tv * oppV
			} else {
				tv = oppV / tv
			}
		case "=":
			tv = oppV
		}
	}

	return tv
}
