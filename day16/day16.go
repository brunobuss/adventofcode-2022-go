package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

type memoState struct {
	mins       int
	pos        int
	valveState int64
}

func init() {
	log.SetOutput(io.Discard)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	vAtoI := make(map[string]int)
	tunS := make([][]string, 0)
	flow := make([]int, 0)
	allowed := make([]int, 0) // The non-zero valves

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var o, dst string
		var f int

		line = strings.ReplaceAll(line, ", ", ",")

		var lineFmt string
		if strings.Contains(line, "valves") {
			lineFmt = "Valve %s has flow rate=%d; tunnels lead to valves %s"
		} else {
			lineFmt = "Valve %s has flow rate=%d; tunnel leads to valve %s"
		}

		n, err := fmt.Sscanf(line, lineFmt, &o, &f, &dst)
		if n != 3 {
			log.Fatalln(err)
		}

		vAtoI[o] = len(flow)
		if f > 0 {
			allowed = append(allowed, len(flow))
		}
		flow = append(flow, f)
		tunS = append(tunS, strings.Split(dst, ","))
	}

	// Convert the edge list from string to ints to make it easier to query them
	tunI := make([][]int, 0)
	for _, tuns := range tunS {
		tI := []int{}
		for _, t := range tuns {
			tI = append(tI, vAtoI[t])
		}
		tunI = append(tunI, tI)
	}

	log.Println("Flow:", flow)
	log.Println("TunS:", tunS)
	log.Println("TunI:", tunI)

	valveState := make([]bool, len(flow))
	memo := make(map[memoState]int)
	fw := floydWarshall(tunI, flow)
	r := partOneBackTrackWithMemo(30, vAtoI["AA"], valveState, memo, tunI, flow, fw, allowed)
	fmt.Println("PartOne - Best flow starting from AA with 30 minutes and one agent:", r)

	r = partTwoBacktrackCombination(vAtoI["AA"], valveState, tunI, flow, fw, allowed)
	fmt.Println("PartTwo - Best flow starting from AA with 26 minutes and two agents:", r)
}

func floydWarshall(tuns [][]int, flow []int) [][]int {
	fw := make([][]int, len(flow))
	// Init matrix
	for i := 0; i < len(flow); i++ {
		fw[i] = make([]int, len(flow))
		for j := 0; j < len(flow); j++ {
			fw[i][j] = math.MaxInt16
		}
		fw[i][i] = 0
	}
	// Mark tunels with distance 1
	for i, ts := range tuns {
		for _, j := range ts {
			fw[i][j] = 1
		}
	}
	// Compute Floyd-Warshall
	for k := 0; k < len(flow); k++ {
		for i := 0; i < len(flow); i++ {
			for j := 0; j < len(flow); j++ {
				if fw[i][k]+fw[k][j] < fw[i][j] {
					fw[i][j] = fw[i][k] + fw[k][j]
				}
			}
		}
	}
	return fw
}

func partTwoBacktrackCombination(
	startPos int,
	valveState []bool,
	tuns [][]int,
	flow []int,
	fw [][]int,
	allowedValves []int,
) int {
	best := 0

	// Since the number of allowedValves are not very big, lets
	// attempt all the possible partitions of the valves between
	// two backtrack attemps and sum the results.
	e := 1 << (len(allowedValves) - 1)
	for s := 0; s < e; s++ {
		allowPlayer := []int{}
		allowElephant := []int{}

		for i, v := range allowedValves {
			if (s>>i)&1 == 1 {
				allowPlayer = append(allowPlayer, v)
			} else {
				allowElephant = append(allowElephant, v)
			}
		}

		memo := make(map[memoState]int)
		pB := partOneBackTrackWithMemo(26, startPos, valveState, memo, tuns, flow, fw, allowPlayer)
		memo = make(map[memoState]int)
		pE := partOneBackTrackWithMemo(26, startPos, valveState, memo, tuns, flow, fw, allowElephant)

		if pB+pE > best {
			best = pB + pE
		}
	}

	return best
}

func partOneBackTrackWithMemo(mLeft int,
	pos int,
	valveState []bool,
	memo map[memoState]int,
	tuns [][]int,
	flow []int,
	fw [][]int,
	allowedValves []int,
) int {
	state := memoState{
		mins:       mLeft,
		pos:        pos,
		valveState: valveStateToInt64(valveState),
	}

	if mLeft <= 0 {
		return 0
	}

	if m, e := memo[state]; e {
		return m
	}

	// Flow of already opened valves
	flowPerMinute := 0
	for i, s := range valveState {
		if s {
			flowPerMinute += flow[i]
		}
	}

	best := mLeft * flowPerMinute // Do nothing, just stay here

	// Go open a valve, but only try the ones we got time for
	for _, i := range allowedValves {
		if !valveState[i] && fw[pos][i]+1 < mLeft {
			// d = min distance + 1 turn to open the valve
			d := fw[pos][i] + 1
			valveState[i] = true
			r := partOneBackTrackWithMemo(mLeft-d, i, valveState, memo, tuns, flow, fw, allowedValves)
			if r+flowPerMinute*d > best {
				best = r + flowPerMinute*d
			}
			valveState[i] = false
		}
	}

	memo[state] = best
	return best
}

func valveStateToInt64(valveState []bool) int64 {
	var i64 int64
	for _, v := range valveState {
		i64 <<= 1
		if v {
			i64 += 1
		}
	}
	return i64
}
