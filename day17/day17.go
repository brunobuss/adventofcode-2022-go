package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const DEBUG = false

func init() {
	log.SetOutput(io.Discard)
}

type coord struct {
	x, y int
}

type rockShape struct {
	height      int
	fillOffsets []coord // Coord offset from top-left that contains a rock piece
}

var SHAPES []rockShape = []rockShape{
	// ####
	{
		height: 1,
		fillOffsets: []coord{
			{0, 0}, {1, 0}, {2, 0}, {3, 0},
		},
	},
	// .#.
	// ###
	// .#.
	{
		height: 3,
		fillOffsets: []coord{
			{1, 0}, {0, -1}, {1, -1}, {2, -1}, {1, -2},
		},
	},
	// ..#
	// ..#
	// ###
	{
		height: 3,
		fillOffsets: []coord{
			{2, 0}, {2, -1}, {0, -2}, {1, -2}, {2, -2},
		},
	},
	// #
	// #
	// #
	// #
	{
		height: 4,
		fillOffsets: []coord{
			{0, 0}, {0, -1}, {0, -2}, {0, -3},
		},
	},
	// ##
	// ##
	{
		height: 2,
		fillOffsets: []coord{
			{0, 0}, {1, 0}, {0, -1}, {1, -1},
		},
	},
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		log.Fatal(err)
	}
	gasLine := scanner.Text()

	tOne := simulate(2022, gasLine)
	fmt.Println("[PartOne] Top after 2022 rocks:", tOne)

	tTwo := simulate(1000000000000, gasLine)
	fmt.Println("[PartTwo] Top after 1000000000000 rocks:", tTwo)
}

// Simulate rockLimit rocks and returns the height of the topmost piece
func simulate(rockLimit int, shifts string) int {
	shiftI := 0
	prevTop := 0
	top := 0
	settledRocks := make(map[coord]bool)

	// I don't think this is the minimal state for the cycle detection,
	// but seems to work ...
	type memoState struct {
		rock, bs, es int
		finalOffset  coord
		topDiff      int
	}
	memoRound := make(map[memoState]int)
	topMemo := make(map[memoState]int)
	topDiffCycle := []int{}

	for rockSpawned := 0; rockSpawned < rockLimit; rockSpawned++ {
		log.Printf("=== Rock %d - Shape %d ===\n", rockSpawned+1, rockSpawned%len(SHAPES))
		currentRock := SHAPES[rockSpawned%len(SHAPES)]
		startTop := top + 3 + currentRock.height
		rockPos := coord{2, startTop}

		var st memoState

		// Where we began to shift for this rock
		bs := shiftI % len(shifts)

		for {
			// First we try to shift
			shiftPos := rockPos
			sD := shifts[shiftI%len(shifts)]
			switch sD {
			case '<':
				shiftPos.x--
			case '>':
				shiftPos.x++
			default:
				log.Fatalln("Unrecorgnized shift:", sD)
			}
			log.Printf("Attemp shift (%c) from %v to %v\n", sD, rockPos, shiftPos)
			if !collision(shiftPos, currentRock, settledRocks) {
				log.Println("Shift done")
				rockPos = shiftPos

			}
			shiftI++

			// Now we try to make the rock go down one step
			downPos := coord{rockPos.x, rockPos.y - 1}
			if !collision(downPos, currentRock, settledRocks) {
				log.Printf("Rock going down from %v to %v\n", rockPos, downPos)
				rockPos = downPos
			} else {
				for _, o := range currentRock.fillOffsets {
					settledRocks[rockPos.plusOffset(o)] = true
				}
				prevTop = top
				if rockPos.y > top {
					top = rockPos.y
				}
				break
			}
		}

		// Where we end the shifts for this rock
		es := shiftI % len(shifts)
		fo := coord{
			rockPos.x - 2,
			rockPos.y - startTop,
		}
		st = memoState{
			rockSpawned % len(SHAPES),
			bs, es,
			fo,
			top - prevTop,
		}
		if v, e := memoRound[st]; e {
			topDiffCycle = append(topDiffCycle, top-prevTop)

			log.Printf("(%d) Rock %d is of shape %d and just like rock %d (%d ago), took shifts %d to %d and moved %v.\nWe now have top %d, back then we had %d (dS: %d/dC:%d).\n",
				len(topDiffCycle), rockSpawned, st.rock, v, rockSpawned-v, bs, es, fo, top, topMemo[st], top-prevTop, top-topMemo[st])

			memoRound[st] = rockSpawned

			// If we are in a cycle for sure, then we can change the logic to jump ahead
			if len(topDiffCycle) == rockSpawned-v+1 {
				cycleSize := rockSpawned - v
				topInc := (top - topMemo[st]) / 2
				log.Printf("Stopping now: rock=%d, cycle=%d, topInc=%d\n", rockSpawned, cycleSize, topInc)

				qCycles := (rockLimit - (rockSpawned)) / cycleSize
				rem := (rockLimit - (rockSpawned)) % cycleSize

				futRock := rockSpawned + qCycles*cycleSize
				futTop := top + qCycles*topInc

				log.Printf("Looking ahead %d cycles: rock=%d, top=%d, reminder=%d\n", qCycles, futRock, futTop, rem)
				rockSpawned = futRock
				top = futTop

				// If we are targetting an exactly cycle match, it seems we double count
				// the first element of the cycle (since we only fully detect we're in the
				// cycle after passing over it for the second time)
				// Same reason why we start from i=1 if we still got some remaining steps.
				if rem == 0 {
					top -= topDiffCycle[0]
				} else {
					for i := 1; i < rem; i++ {
						rockSpawned++
						top += topDiffCycle[i]
						log.Printf("R:%d T:%d\n", rockSpawned, top)
					}
				}

			}
		} else {
			topDiffCycle = []int{}
			memoRound[st] = rockSpawned
			topMemo[st] = top
		}
	}

	if DEBUG {
		for t := top; t > 0; t-- {
			fmt.Print("|")
			for x := 0; x < 7; x++ {
				if settledRocks[coord{x, t}] {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println("|")
		}
		fmt.Println("+-------+")
	}

	return top
}

func collision(pos coord, rock rockShape, settledRocks map[coord]bool) bool {
	for _, o := range rock.fillOffsets {
		nC := pos.plusOffset(o)
		if nC.y <= 0 || nC.x < 0 || nC.x >= 7 || settledRocks[nC] {
			return true
		}
	}
	return false
}

func (c coord) plusOffset(o coord) coord {
	return coord{
		c.x + o.x,
		c.y + o.y,
	}
}
