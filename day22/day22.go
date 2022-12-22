package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	log.SetOutput(io.Discard)
}

type coord struct {
	x, y int
}

type player struct {
	pos coord
	f   int
}

type wf func(npos coord, p player, bm []string) (coord, int)

var OFFSETS []coord = []coord{
	{1, 0},  // 0 = facing right
	{0, 1},  // 1 = facing down
	{-1, 0}, // 2 = facing left
	{0, -1}, // 3 = facing up
}

func (p player) turn(c string) player {
	switch c {
	case "R":
		p.f = (p.f + 1) % 4
	case "L":
		p.f--
		if p.f == -1 {
			p.f = 3
		}
	default:
		log.Fatal("Unrecognized command to turn:", c)
	}

	return p
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bm := make([]string, 0)
	p := player{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		if len(bm) == 0 {
			p.pos.x = strings.IndexByte(line, '.')
		}
		bm = append(bm, line)
	}

	if !scanner.Scan() {
		log.Fatalln("Couldn't scan instructions")
	}
	// Add spaces around R's and L's so we can easily split everything into
	// a list of commands to follow.
	ins := strings.Split(
		strings.ReplaceAll(
			strings.ReplaceAll(scanner.Text(), "L", " L "), "R", " R "), " ")

	partOne := followInstructions(p, ins, bm, wrapPartOne)
	log.Printf("Row:%d; Column:%d; Facing: %d\n", partOne.pos.y+1, partOne.pos.x+1, p.f)
	passOne := (partOne.pos.y+1)*1000 + (partOne.pos.x+1)*4 + partOne.f
	fmt.Println("[PartOne] Password:", passOne)

	partTwo := followInstructions(p, ins, bm, wrapPartTwo)
	log.Printf("Row:%d; Column:%d; Facing: %d\n", partTwo.pos.y+1, partTwo.pos.x+1, p.f)
	passTwo := (partTwo.pos.y+1)*1000 + (partTwo.pos.x+1)*4 + partTwo.f
	fmt.Println("[PartTwo] Password:", passTwo)
}

func followInstructions(p player, ins []string, bm []string, w wf) player {
	for _, d := range ins {
		if d == "R" || d == "L" {
			p = p.turn(d)
			log.Println("Turn", d, p.pos, p.f)
			continue
		}

		l, err := strconv.Atoi(d)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Walking", l)

		for i := 0; i < l; i++ {
			o := OFFSETS[p.f]
			npos := coord{
				p.pos.x + o.x,
				p.pos.y + o.y,
			}
			log.Println(npos, p.f, string(whichCube(npos)))

			npos, nf := w(npos, p, bm)
			log.Println(npos, nf, string(whichCube(npos)))

			c := bm[npos.y][npos.x]
			if c == '#' {
				break
			} else if c == '.' {
				p.pos = npos
				p.f = nf
			} else {
				log.Fatal("Invalid map element at", npos, ":", bm[npos.y][npos.x])
			}
		}
	}

	return p
}

func isWrapping(npos coord, bm []string) bool {
	return npos.y < 0 || npos.y >= len(bm) ||
		npos.x < 0 || npos.x >= len(bm[npos.y]) ||
		bm[npos.y][npos.x] == ' '
}

func wrapPartOne(npos coord, p player, bm []string) (coord, int) {
	if !isWrapping(npos, bm) {
		return npos, p.f
	}

	switch p.f {
	case 0:
		// We're going right, so we need to loop left
		npos.x = strings.IndexAny(bm[npos.y], ".#")
	case 1:
		// We're going down, so we need to loop up
		for ny := 0; ny < len(bm); ny++ {
			if npos.x < len(bm[ny]) && bm[ny][npos.x] != ' ' {
				npos.y = ny
				break
			}
		}
	case 2:
		// We're going left, so we need to loop right
		npos.x = strings.LastIndexAny(bm[npos.y], ".#")
	case 3:
		// We're going up, so we need to loop down
		for ny := len(bm) - 1; ny >= 0; ny-- {
			if npos.x < len(bm[ny]) && bm[ny][npos.x] != ' ' {
				npos.y = ny
				break
			}
		}
	default:
		log.Fatal("Invalid facing:", p.f)
	}

	return npos, p.f
}

// Not sure how to compute the cube folding programatically, so lets hardcode it
// For my input, the cube looks like:
//
// ..AABB
// ..AABB
// ..CC..
// ..CC..
// DDEE..
// DDEE..
// FF....
// FF....
// After using some advanced post-it origami techniques...
var cubePos map[byte]coord = map[byte]coord{
	'A': {1, 0},
	'B': {2, 0},
	'C': {1, 1},
	'D': {0, 2},
	'E': {1, 2},
	'F': {0, 3},
}

func inCube(pos, cube coord) bool {
	return pos.x >= cubeLeftEdge(cube) && pos.x <= cubeRightEdge(cube) &&
		pos.y >= cubeTopEdge(cube) && pos.y <= cubeBottomEdge(cube)
}

func cubeLeftEdge(cube coord) int {
	return cube.x * 50
}

func cubeRightEdge(cube coord) int {
	return (cube.x+1)*50 - 1
}

func cubeTopEdge(cube coord) int {
	return cube.y * 50
}

func cubeBottomEdge(cube coord) int {
	return (cube.y+1)*50 - 1
}

func whichCube(p coord) byte {
	for k, v := range cubePos {
		if inCube(p, v) {
			return k
		}
	}
	return '?'
}

func wrapPartTwo(npos coord, p player, bm []string) (coord, int) {
	if !isWrapping(npos, bm) {
		return npos, p.f
	}

	if inCube(p.pos, cubePos['A']) {
		switch p.f {
		case 0:
			// Right of A if left of B
		case 1:
			// Bottom of A is top of C
		case 2:
			// Left of A is left of D
			return coord{
				cubeLeftEdge(cubePos['D']),
				cubeBottomEdge(cubePos['D']) - (p.pos.y - cubeTopEdge(cubePos['A'])),
			}, 0
		case 3:
			// Top of A is left of F
			return coord{
				cubeLeftEdge(cubePos['F']),
				cubeTopEdge(cubePos['F']) + (p.pos.x - cubeLeftEdge(cubePos['A'])),
			}, 0
		}
	} else if inCube(p.pos, cubePos['B']) {
		switch p.f {
		case 0:
			// Right of B is right of E
			return coord{
				cubeRightEdge(cubePos['E']),
				cubeBottomEdge(cubePos['E']) - (p.pos.y - cubeTopEdge(cubePos['B'])),
			}, 2
		case 1:
			// Bottom of B is right of C
			return coord{
				cubeRightEdge(cubePos['C']),
				cubeTopEdge(cubePos['C']) + (p.pos.x - cubeLeftEdge(cubePos['B'])),
			}, 2
		case 2:
			// Left of B is right of A
		case 3:
			// Top of B is bottom of F
			return coord{
				cubeLeftEdge(cubePos['F']) + (p.pos.x - cubeLeftEdge(cubePos['B'])),
				cubeBottomEdge(cubePos['F']),
			}, 3
		}
	} else if inCube(p.pos, cubePos['C']) {
		switch p.f {
		case 0:
			// Right of C is bottom of B
			return coord{
				cubeLeftEdge(cubePos['B']) + (p.pos.y - cubeTopEdge(cubePos['C'])),
				cubeBottomEdge(cubePos['B']),
			}, 3
		case 1:
			// Bottom of C is top of E
		case 2:
			// Left of C is top of D
			return coord{
				cubeLeftEdge(cubePos['D']) + (p.pos.y - cubeTopEdge(cubePos['C'])),
				cubeTopEdge(cubePos['D']),
			}, 1
		case 3:
			// Top of C is bottom of A
		}
	} else if inCube(p.pos, cubePos['D']) {
		switch p.f {
		case 0:
			// Right of D is left of E
		case 1:
			// Bottom of D is top of F
		case 2:
			// Left of D is left of A
			return coord{
				cubeLeftEdge(cubePos['A']),
				cubeBottomEdge(cubePos['A']) - (p.pos.y - cubeTopEdge(cubePos['D'])),
			}, 0
		case 3:
			// Top of D is left of C
			return coord{
				cubeLeftEdge(cubePos['C']),
				cubeTopEdge(cubePos['C']) + (p.pos.x - cubeLeftEdge(cubePos['D'])),
			}, 0
		}
	} else if inCube(p.pos, cubePos['E']) {
		switch p.f {
		case 0:
			// Right of E is right of B
			return coord{
				cubeRightEdge(cubePos['B']),
				cubeTopEdge(cubePos['B']) + (cubeBottomEdge(cubePos['E']) - p.pos.y),
			}, 2
		case 1:
			// Bottom of E is right of F
			return coord{
				cubeRightEdge(cubePos['F']),
				cubeTopEdge(cubePos['F']) + (p.pos.x - cubeLeftEdge(cubePos['E'])),
			}, 2
		case 2:
			// Left of E is right of D
		case 3:
			// Top of E is bottom of C
		}
	} else if inCube(p.pos, cubePos['F']) {
		switch p.f {
		case 0:
			// Right of F is bottom of E
			return coord{
				cubeLeftEdge(cubePos['E']) + (p.pos.y - cubeTopEdge(cubePos['F'])),
				cubeBottomEdge(cubePos['E']),
			}, 3
		case 1:
			// Bottom of F is top of B
			return coord{
				cubeLeftEdge(cubePos['B']) + (p.pos.x - cubeLeftEdge(cubePos['F'])),
				cubeTopEdge(cubePos['B']),
			}, 1
		case 2:
			// Left of F is top of A
			return coord{
				cubeLeftEdge(cubePos['A']) + (p.pos.y - cubeTopEdge(cubePos['F'])),
				cubeTopEdge(cubePos['A']),
			}, 1
		case 3:
			// Top of F is bottom of D
		}
	} else {
		log.Fatalln("Where the heck we're?", p.pos)
	}

	return npos, p.f
}
