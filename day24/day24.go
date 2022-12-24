package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

type coord struct {
	x, y int
}

type blizzard struct {
	direction byte
}

type state struct {
	pos coord
	r   int
}

type memoret struct {
	v int
	s bool
}

func (b blizzard) moveOffset() coord {
	switch b.direction {
	case 1:
		return coord{0, -1}
	case 2:
		return coord{1, 0}
	case 3:
		return coord{0, 1}
	case 4:
		return coord{-1, 0}
	}
	log.Fatalln("Invalid blizzard direction", b)
	return coord{0, 0}
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

	scanner := bufio.NewScanner(file)
	bm := make(map[coord][]blizzard)
	mwidth := 0
	mheight := 0

	for scanner.Scan() {
		line := scanner.Text()

		for i, r := range line {
			b := blizzard{}
			switch r {
			case '^':
				b.direction = 1
			case '>':
				b.direction = 2
			case 'v':
				b.direction = 3
			case '<':
				b.direction = 4
			}
			if b.direction > 0 {
				bm[coord{i, mheight}] = []blizzard{b}
			}
		}

		mwidth = len(line)
		mheight++
	}

	start := coord{1, 0}
	end := coord{mwidth - 2, mheight - 1}
	size := coord{mwidth, mheight}

	log.Println("size", size)

	// Every rmax rounds, because of the cyclic nature of the blizzards
	// they return to the same spot at in the beginning.
	// So we can precompute all the blizzard positions.
	lcm := lcm(size.x-2, size.y-2)
	pbm := make([]map[coord][]blizzard, lcm)
	pbm[0] = bm
	for r := 1; r < lcm; r++ {
		pbm[r] = moveBlizzards(size, pbm[r-1])
	}

	memoGoing := make(map[state]memoret)
	rFirstTrip, s := explore(0, lcm, start, end, size, pbm, memoGoing)
	fmt.Println("[PartOne] Rounds until exit:", rFirstTrip, s)

	memoBack := make(map[state]memoret)
	rBackTrip, s := explore(rFirstTrip, lcm, end, start, size, pbm, memoBack)
	fmt.Println("[PartTwo] Rounds to come back:", rBackTrip, s)

	rSecondTrip, s := explore(rFirstTrip+rBackTrip, lcm, start, end, size, pbm, memoGoing)
	fmt.Println("[PartTwo] Rounds to go again:", rSecondTrip, s)

	fmt.Println("[PartTwo] Total Rounds:", rFirstTrip+rBackTrip+rSecondTrip, s)
}

func explore(r, rmax int, pos, dest, size coord, pbm []map[coord][]blizzard, memo map[state]memoret) (int, bool) {
	log.Println("R", r, "P", pos)
	if pos == dest {
		// Arrived at the exit
		return 0, true
	}

	state := state{pos, r}
	if v, e := memo[state]; e {
		return v.v, v.s
	}

	// r > 1000 is to avoid stackoverflow
	// If needed (wrong answer) we can convert this recursion into a
	// loop based with stack
	if !canMove(pos, size) || r > 1000 || len(pbm[r%rmax][pos]) > 0 {
		// Impossible move
		return 0, false
	}

	nr := r + 1
	b := -1

	var e int
	var s bool
	var ss bool
	// Go South
	e, s = explore(nr, rmax, coord{pos.x, pos.y + 1}, dest, size, pbm, memo)
	if s && (b == -1 || e < b) {
		b = e
		ss = true
	}
	// Go East
	e, s = explore(nr, rmax, coord{pos.x + 1, pos.y}, dest, size, pbm, memo)
	if s && (b == -1 || e < b) {
		b = e
		ss = true
	}
	// Stay
	e, s = explore(nr, rmax, coord{pos.x, pos.y}, dest, size, pbm, memo)
	if s && (b == -1 || e < b) {
		b = e
		ss = true
	}
	// Go West
	e, s = explore(nr, rmax, coord{pos.x - 1, pos.y}, dest, size, pbm, memo)
	if s && (b == -1 || e < b) {
		b = e
		ss = true
	}
	// Go North
	e, s = explore(nr, rmax, coord{pos.x, pos.y - 1}, dest, size, pbm, memo)
	if s && (b == -1 || e < b) {
		b = e
		ss = true
	}

	memo[state] = memoret{b + 1, ss}
	return b + 1, ss
}

func moveBlizzards(size coord, bm map[coord][]blizzard) map[coord][]blizzard {
	nbm := make(map[coord][]blizzard)
	for k, v := range bm {
		for _, b := range v {
			o := b.moveOffset()
			npos := coord{k.x + o.x, k.y + o.y}
			if npos.x == 0 {
				npos.x = size.x - 2
			} else if npos.x == size.x-1 {
				npos.x = 1
			}

			if npos.y == 0 {
				npos.y = size.y - 2
			} else if npos.y == size.y-1 {
				npos.y = 1
			}

			if _, e := nbm[npos]; !e {
				nbm[npos] = []blizzard{}
			}
			nbm[npos] = append(nbm[npos], b)
		}
	}
	return nbm
}

func canMove(pos, size coord) bool {
	if pos.x <= 0 || pos.x >= size.x-1 {
		// First or second columns, all walls
		return false
	}
	if pos.y < 0 || pos.y >= size.y {
		// Out of boudns
		return false
	}
	if pos.y == 0 && pos.x != 1 {
		// First row, but not entrance
		return false
	}
	if pos.y == size.y-1 && pos.x != size.x-2 {
		// Last row, but not exit
		return false
	}
	return true
}

func lcm(a, b int) int {
	mfa := make(map[int]int)
	for _, f := range factors(a) {
		mfa[f]++
	}
	mfb := make(map[int]int)
	for _, f := range factors(b) {
		mfb[f]++
	}

	for k, v := range mfb {
		if v > mfa[k] {
			mfa[k] = v
		}
	}

	r := 1
	for k, v := range mfa {
		for i := 0; i < v; i++ {
			r *= k
		}
	}

	return r
}

func factors(a int) []int {
	f := make([]int, 0)
	r := a
	for i := 2; i <= int(math.Sqrt(float64(a)))+1 && r > 1; {
		if r%i == 0 {
			f = append(f, i)
			r = r / i
		} else {
			i++
		}
	}

	if r > 1 {
		f = append(f, r)
	}

	return f
}
