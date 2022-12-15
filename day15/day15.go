package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

const SAMPLE = false

type coords struct {
	x, y int
}

type sensor struct {
	c  coords
	md int
}

type beacon struct {
	c coords
}

type intRange struct {
	begin, end int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sl := []sensor{}
	// Use a map for the beacons, in order to de-dup same beacon observed
	// by multiple sensors
	bm := make(map[beacon]bool)
	for scanner.Scan() {
		line := scanner.Text()
		var s sensor
		var b beacon
		n, err := fmt.Sscanf(
			line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&s.c.x, &s.c.y, &b.c.x, &b.c.y)
		if n != 4 {
			log.Fatalln(err)
		}

		s.md = distManhatan(s.c, b.c)
		sl = append(sl, s)
		bm[b] = true
	}

	bl := []beacon{}
	for k := range bm {
		bl = append(bl, k)
	}

	partOne(sl, bl)
	partTwo(sl, bl)
}

func partOne(sl []sensor, bl []beacon) {
	targetY := 2000000
	if SAMPLE {
		targetY = 10

	}

	mrs := coverAtY(sl, bl, targetY)
	coverCount := 0
	for _, r := range mrs {
		coverCount += (r.end - r.begin)
	}

	// We need to subtract the number of beacons in this line
	// since they don't configure exactly as cells that a beacon
	// cannot exists
	for _, b := range bl {
		if b.c.y == targetY {
			coverCount--
		}
	}

	fmt.Printf("At Y=%d, %d positions can't have a beacon\n", targetY, coverCount)
}

func partTwo(sl []sensor, bl []beacon) {
	maxCoord := 4000000
	if SAMPLE {
		maxCoord = 20
	}

	for y := 0; y < maxCoord; y++ {
		mrs := coverAtY(sl, bl, y)
		if len(mrs) > 1 {
			fmt.Printf("Found gap at (%d, %d)\n", mrs[0].end, y)
			fmt.Printf("Tunning Frequency: %d\n", tuningFreq(coords{mrs[0].end, y}))
		}
	}
}

func coverAtY(sl []sensor, bl []beacon, y int) []intRange {
	rs := []intRange{}
	for _, s := range sl {
		rs = append(rs, rangeOfSensorAtY(s, y))
	}
	for _, b := range bl {
		if b.c.y == y {
			rs = append(rs, intRange{b.c.x, b.c.x + 1})
		}
	}

	mrs := mergeRangeList(rs)
	return mrs
}

func mergeRangeList(ranges []intRange) []intRange {
	sort.Slice(ranges, func(i, j int) bool {
		return (ranges[i].begin < ranges[j].begin) ||
			(ranges[i].begin == ranges[j].begin && (ranges[i].end <= ranges[j].end))
	})
	mrs := []intRange{}

	nr := ranges[0]
	for _, r := range ranges[1:] {
		if r.begin <= nr.end {
			// If next range starts before our end, then we join them
			nr.end = maxInt(nr.end, r.end)
		} else {
			mrs = append(mrs, nr)
			nr = r
		}
	}
	mrs = append(mrs, nr)

	return mrs
}

func tuningFreq(c coords) int {
	return c.x*4000000 + c.y
}

func rangeOfSensorAtY(s sensor, y int) intRange {
	yDiff := absDiff(s.c.y, y)

	if yDiff > s.md {
		return intRange{s.c.x, s.c.x}
	}

	xDiff := s.md - yDiff
	begin := s.c.x - xDiff
	end := s.c.x + xDiff + 1
	return intRange{begin, end}
}

func distManhatan(a, b coords) int {
	return absDiff(a.x, b.x) + absDiff(a.y, b.y)
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
