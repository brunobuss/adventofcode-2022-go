package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func init() {
	log.SetOutput(io.Discard)
}

type dir struct {
	name    string
	subdirs map[string]*dir
	files   map[string]uint64
	size    uint64
	parent  *dir
}

var RE_CD = regexp.MustCompile(`^\$ cd (.+)$`)
var RE_LS = regexp.MustCompile(`^\$ ls$`)
var RE_DIR = regexp.MustCompile(`^dir (.+)$`)
var RE_FILE = regexp.MustCompile(`^([0-9]+) (.+)$`)

const FS_TOTAL = 70_000_000
const SPACE_NEEDED = 30_000_000

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	root := parse(scanner)
	printTree(root, 0)

	sum := findAllSizesUpTo(root, 100000)
	fmt.Println("Sum up (less than size 100000): ", sum)

	needed := SPACE_NEEDED - (FS_TOTAL - root.size)
	log.Println("Space needed: ", needed)
	toFree := findDirToRemove(root, needed)
	fmt.Printf("Free-up dir `%s`: %d\n", rebuildPath(toFree), toFree.size)
}

func parse(scanner *bufio.Scanner) *dir {
	root := &dir{
		name:    "/",
		subdirs: make(map[string]*dir),
		files:   make(map[string]uint64),
		size:    0,
		parent:  nil,
	}
	curDir := root

	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Parsing: ", line)
		if result := RE_CD.FindAllStringSubmatch(line, 1); result != nil {
			param := result[0][1]
			log.Println("It's a cd command with param: ", param)
			switch param {
			case "/":
				log.Println("Going back to root")
				curDir = root
			case "..":
				log.Println("Going to parent: ", curDir.parent.name)
				curDir = curDir.parent
			default:
				log.Println("Moving to subdir")
				if d, e := curDir.subdirs[param]; e {
					curDir = d
				} else {
					log.Println("CD-ing into unknown subdir: ", param)
				}
			}
		} else if result := RE_LS.FindAllStringSubmatch(line, 1); result != nil {
			// Not necessary to do anything
		} else if result := RE_DIR.FindAllStringSubmatch(line, 1); result != nil {
			param := result[0][1]
			log.Println("It's a dir entry with name: ", param)
			// Only add if there is not already there, to prevent a second listing
			// to erase already computed information
			if _, e := curDir.subdirs[param]; !e {
				curDir.subdirs[param] = &dir{
					name:    param,
					subdirs: make(map[string]*dir),
					files:   make(map[string]uint64),
					size:    0,
					parent:  curDir,
				}
			}
		} else if result := RE_FILE.FindAllStringSubmatch(line, 1); result != nil {
			size, err := strconv.ParseUint(result[0][1], 10, 64)
			if err != nil {
				log.Fatalln("Error converting size to uint: ", result[0][1])
			}
			filename := result[0][2]
			log.Printf("It's a file entry with (%s, %d)\n", filename, size)
			curDir.files[filename] = size
		} else {
			log.Fatalln("Unrecognized line: ", line)
		}
	}
	root.updateSizes()
	return root
}

func (d *dir) updateSizes() {
	var newSize uint64 = 0
	for k := range d.subdirs {
		d.subdirs[k].updateSizes()
		newSize += d.subdirs[k].size
	}
	for _, v := range d.files {
		newSize += v
	}
	d.size = newSize
}

func findAllSizesUpTo(d *dir, limit uint64) uint64 {
	var sum uint64 = 0
	for _, v := range d.subdirs {
		sum += findAllSizesUpTo(v, limit)
	}
	if d.size <= limit {
		sum += d.size
	}
	return sum
}

// findDirToRemove will always return the smallest directory entry
// bigger than `needed` from the `d` and its subdirectories
func findDirToRemove(d *dir, needed uint64) *dir {
	var toRemove *dir = nil
	for _, v := range d.subdirs {
		s := findDirToRemove(v, needed)
		// If we found a candidate from our subdirectories and either
		// we don't have one yet or it's a better fit, we replace our
		// current answer
		if s != nil && (toRemove == nil || s.size < toRemove.size) {
			toRemove = s
		}
	}
	// The current directory can't be smaller than any of the subdirectories
	// so only in case we didn't found any possible candidates in our subdirs
	// we check if this directory is a valid candidate
	if toRemove == nil && d.size >= needed {
		toRemove = d
	}
	return toRemove
}

// -- Debug Utils ---

func printTree(d *dir, depth int) {
	printSpaces(depth * 2)
	log.Printf("- %s (dir, size=%d)\n", d.name, d.size)
	depth++

	for _, v := range d.subdirs {
		printTree(v, depth)
	}

	for k, v := range d.files {
		printSpaces(depth * 2)
		log.Printf("- %s (file, size=%d)\n", k, v)
	}
}

func printSpaces(q int) {
	for i := 0; i < q; i++ {
		log.Print(" ")
	}
}

func rebuildPath(d *dir) string {
	if d.parent == nil {
		return "/"
	}
	return rebuildPath(d.parent) + "/" + d.name
}
