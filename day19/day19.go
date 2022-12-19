package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	log.SetOutput(io.Discard)
}

type cost struct {
	ore, clay, obsidian int
}

type blueprint struct {
	oreRobot, clayRobot, obsidianRobot, geodeRobot cost
}

type fleet struct {
	ore, clay, obsidian, geode int
}

type stockpile struct {
	ore, clay, obdisian, geode int
}

type state struct {
	minutes int
	f       fleet
	s       stockpile
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	qls := 0
	t3 := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		b := blueprint{}
		bn := 0
		n, err := fmt.Sscanf(line,
			"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&bn, &b.oreRobot.ore, &b.clayRobot.ore, &b.obsidianRobot.ore, &b.obsidianRobot.clay, &b.geodeRobot.ore, &b.geodeRobot.obsidian,
		)
		if n != 7 {
			log.Fatal(err)
		}

		ss := state{}
		ss.minutes = 24
		ss.f.ore = 1
		memo := make(map[state]int)

		r := backtrack(b, ss, memo)
		log.Printf("24 min - Bp %d, Geodes: %d\n", bn, r)
		log.Println("Memo Size:", len(memo))

		if bn <= 3 {
			ss.minutes = 32
			// We re-use memo since we're still on the same blueprint.
			r32 := backtrack(b, ss, memo)
			log.Printf("32 min - Bp %d, Geodes: %d\n", bn, r32)
			log.Println("Memo Size:", len(memo))
			t3 *= r32
		}

		qls += bn * r
	}

	fmt.Println("[PartOne] Quality Level Sum:", qls)
	fmt.Println("[PartTwo] Top3 Geode Mul:", t3)
}

func backtrack(b blueprint, s state, memo map[state]int) int {
	if s.minutes == 0 {
		return s.s.geode
	}

	// Cap search space if we ever have more stockpile than we could ever use.
	if s.s.obdisian >= b.geodeRobot.obsidian*s.minutes {
		s.s.obdisian = b.geodeRobot.obsidian * s.minutes
	}
	if s.s.clay >= b.obsidianRobot.clay*s.minutes {
		s.s.clay = b.obsidianRobot.clay * s.minutes
	}

	if v, e := memo[s]; e {
		return v
	}

	best := s.s.geode

	// If we can keep building more and more Geode robots each turn, just go for it
	// Our fleet is currently self sufficient for producing geode robots for ever?
	usePureFleet := s.f.ore >= b.geodeRobot.ore && s.f.obsidian >= b.geodeRobot.obsidian
	// If not, with our current stockpile we can support it?
	useFleetAndSp := (s.s.ore >= b.geodeRobot.ore) && ((b.geodeRobot.ore-s.f.ore)*s.minutes <= (s.s.ore - b.geodeRobot.ore)) &&
		(s.s.obdisian >= b.geodeRobot.obsidian) && ((b.geodeRobot.obsidian-s.f.obsidian)*s.minutes <= (s.s.obdisian - b.geodeRobot.obsidian))
	if usePureFleet || useFleetAndSp {
		sum := (s.f.geode + s.f.geode + s.minutes - 1) * s.minutes / 2

		// Storing this state in memory doesn't really save a lot of time
		// and uses *a lot* of memory with the sample input.
		// Uncomment if 24+ GB of mem is available =p
		// memo[s] = nS.s
		return sum + s.s.geode
	}

	// Build Ore Robot next if possible and if needed
	if (s.f.ore < b.geodeRobot.ore) ||
		(s.f.ore < b.obsidianRobot.ore && s.f.obsidian < b.geodeRobot.obsidian) ||
		(s.f.ore < b.clayRobot.ore && s.f.clay < b.obsidianRobot.clay && s.f.obsidian < b.geodeRobot.obsidian) {
		if m := minutesToBuild(b.oreRobot, s.s, s.f); m != -1 && m < s.minutes {
			nS := s
			nS.s = produceForTurns(m+1, nS.f, build(b.oreRobot, nS.s))
			nS.f.ore++
			nS.minutes -= m + 1

			r := backtrack(b, nS, memo)
			if r > best {
				best = r
			}
		}
	}

	// Build Clay Robot next if possible and if needed
	if s.f.clay < b.obsidianRobot.clay && s.f.obsidian < b.geodeRobot.obsidian {
		if m := minutesToBuild(b.clayRobot, s.s, s.f); m != -1 && m < s.minutes {
			nS := s
			nS.s = produceForTurns(m+1, nS.f, build(b.clayRobot, nS.s))
			nS.f.clay++
			nS.minutes -= m + 1

			r := backtrack(b, nS, memo)
			if r > best {
				best = r
			}
		}
	}

	// Build Obsidian Robot next if possible and if needed
	if s.f.obsidian < b.geodeRobot.obsidian {
		if m := minutesToBuild(b.obsidianRobot, s.s, s.f); m != -1 && m < s.minutes {
			nS := s
			nS.s = produceForTurns(m+1, nS.f, build(b.obsidianRobot, nS.s))
			nS.f.obsidian++
			nS.minutes -= m + 1

			r := backtrack(b, nS, memo)
			if r > best {
				best = r
			}
		}
	}

	// Build Geode Robot next if possible
	if m := minutesToBuild(b.geodeRobot, s.s, s.f); m != -1 && m < s.minutes {
		nS := s
		nS.s = produceForTurns(m+1, nS.f, build(b.geodeRobot, nS.s))
		nS.f.geode++
		nS.minutes -= m + 1

		r := backtrack(b, nS, memo)
		if r > best {
			best = r
		}
	}

	memo[s] = best
	return best
}

func minutesToBuild(c cost, s stockpile, f fleet) int {
	max := 0

	oM := minutesToAquireRes(c.ore, f.ore, s.ore)
	if oM == -1 {
		return -1
	} else if oM > max {
		max = oM
	}

	cM := minutesToAquireRes(c.clay, f.clay, s.clay)
	if cM == -1 {
		return -1
	} else if cM > max {
		max = cM
	}

	obM := minutesToAquireRes(c.obsidian, f.obsidian, s.obdisian)
	if obM == -1 {
		return -1
	} else if obM > max {
		max = obM
	}

	return max
}

func minutesToAquireRes(required, prodCycle, start int) int {
	if start >= required {
		return 0
	} else if prodCycle == 0 {
		return -1
	}

	needed := required - start
	if needed%prodCycle == 0 {
		return needed / prodCycle
	}
	return (needed / prodCycle) + 1
}

func build(c cost, s stockpile) stockpile {
	s.ore -= c.ore
	s.clay -= c.clay
	s.obdisian -= c.obsidian
	return s
}

func produceForTurns(t int, f fleet, s stockpile) stockpile {
	s.ore += f.ore * t
	s.clay += f.clay * t
	s.obdisian += f.obsidian * t
	s.geode += f.geode * t
	return s
}
