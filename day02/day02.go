package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const OPP_ROCK = "A"
const OPP_PAPER = "B"
const OPP_SCISSORS = "C"

const MY_ROCK = "X"
const MY_PAPER = "Y"
const MY_SCISSORS = "Z"

const SHOULD_LOSE = "X"
const SHOULD_TIE = "Y"
const SHOULD_WIN = "Z"

const ROCK_POINTS = 1
const PAPER_POINTS = 2
const SCISSORS_POINTS = 3

const WIN_POINTS = 6
const TIE_POINTS = 3
const LOSS_POINTS = 0

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pointsPartOne, pointsPartTwo := 0, 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var opp, my string
		n, err := fmt.Sscanln(line, &opp, &my)
		if n != 2 {
			log.Fatal(err)
		}

		pointsPartOne += resultPartOne(opp, my)
		pointsPartTwo += resultPartTwo(opp, my)
	}

	fmt.Println("Total points (first strategy): ", pointsPartOne)
	fmt.Println("Total points (second strategy): ", pointsPartTwo)
}

func resultPartOne(opp, my string) int {
	switch opp {
	case OPP_ROCK:
		switch my {
		case MY_ROCK:
			return TIE_POINTS + ROCK_POINTS
		case MY_PAPER:
			return WIN_POINTS + PAPER_POINTS
		case MY_SCISSORS:
			return LOSS_POINTS + SCISSORS_POINTS
		}
	case OPP_PAPER:
		switch my {
		case MY_ROCK:
			return LOSS_POINTS + ROCK_POINTS
		case MY_PAPER:
			return TIE_POINTS + PAPER_POINTS
		case MY_SCISSORS:
			return WIN_POINTS + SCISSORS_POINTS
		}
	case OPP_SCISSORS:
		switch my {
		case MY_ROCK:
			return WIN_POINTS + ROCK_POINTS
		case MY_PAPER:
			return LOSS_POINTS + PAPER_POINTS
		case MY_SCISSORS:
			return TIE_POINTS + SCISSORS_POINTS
		}
	}

	log.Fatalln("Unrecognized play (", opp, ",", my, ")")
	return 0
}

func resultPartTwo(opp, my string) int {
	switch my {
	case SHOULD_LOSE:
		switch opp {
		case OPP_ROCK:
			return resultPartOne(opp, MY_SCISSORS)
		case OPP_PAPER:
			return resultPartOne(opp, MY_ROCK)
		case OPP_SCISSORS:
			return resultPartOne(opp, MY_PAPER)
		}
	case SHOULD_TIE:
		switch opp {
		case OPP_ROCK:
			return resultPartOne(opp, MY_ROCK)
		case OPP_PAPER:
			return resultPartOne(opp, MY_PAPER)
		case OPP_SCISSORS:
			return resultPartOne(opp, MY_SCISSORS)
		}
	case SHOULD_WIN:
		switch opp {
		case OPP_ROCK:
			return resultPartOne(opp, MY_PAPER)
		case OPP_PAPER:
			return resultPartOne(opp, MY_SCISSORS)
		case OPP_SCISSORS:
			return resultPartOne(opp, MY_ROCK)
		}
	}

	log.Fatalln("Unrecognized strategy 2 (", opp, ",", my, ")")
	return 0
}
