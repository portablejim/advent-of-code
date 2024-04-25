package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Coord struct {
	x int
	y int
	z int
}

type Brick struct {
	num      int
	posStart Coord
	posEnd   Coord
}

func main() {
	var filename = flag.String("f", "../inputs/d22.sample1.txt", "file to use")
	//var is_verbose = flag.Bool("v", false, "verbose")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f", err)
	}

	pBrick := regexp.MustCompile("([0-9]+),([0-9]+),([0-9]+)~([0-9]+),([0-9]+),([0-9]+)")

	brickList := []Brick{}

	for _, fLine := range strings.Split(strings.ReplaceAll(strings.Trim(string(dat), " \r\n"), "\r\n", "\n"), "\n") {
		lineParts := pBrick.FindStringSubmatch(fLine)
		if len(lineParts) != 7 {
			fmt.Printf("Incorrect line %s %v\n", fLine, lineParts)
			continue
		}

		x1, _ := strconv.ParseInt(lineParts[1], 10, 64)
		y1, _ := strconv.ParseInt(lineParts[2], 10, 64)
		z1, _ := strconv.ParseInt(lineParts[3], 10, 64)
		x2, _ := strconv.ParseInt(lineParts[4], 10, 64)
		y2, _ := strconv.ParseInt(lineParts[5], 10, 64)
		z2, _ := strconv.ParseInt(lineParts[6], 10, 64)

		if z1 > z2 {
			x1, y1, z1, x2, y2, z2 = x2, y2, z2, x1, y1, z1
		}

		brickNum := len(brickList)

		targetBrick := Brick{brickNum, Coord{int(x1), int(y1), int(z1)}, Coord{int(x2), int(y2), int(z2)}}

		brickList = append(brickList, targetBrick)
	}

	stackMap := map[int][]int{}
	stackHeightList := []int{}

	for _, currentBrick := range brickList {
		for curZ := currentBrick.posStart.z; curZ <= currentBrick.posEnd.z; curZ += 1 {
			stackBrickIndexList, hasIndex := stackMap[curZ]
			if !hasIndex {
				stackBrickIndexList = []int{}
				stackHeightList = append(stackHeightList, curZ)
			}

			stackBrickIndexList = append(stackBrickIndexList, currentBrick.num)

			stackMap[curZ] = stackBrickIndexList
		}
	}

	lastLayer := 1

	for currentLayerNum := stackHeightList[1]; currentLayerNum < len(stackHeightList); currentLayerNum += 1 {
		currentLayer := stackHeightList[currentLayerNum]
		currentLayerIndexes := stackMap[currentLayer]
		for _, currentLayerIndex := range currentLayerIndexes {
			currentBrick := brickList[currentLayerIndex]
			brickHeight := currentBrick.posEnd.z - currentBrick.posStart.z

			fmt.Printf("%d, %v\n", lastLayer, brickHeight)
		}
	}

	sort.Ints(stackHeightList)

	fmt.Printf("Layers: %v\n", stackHeightList)
	fmt.Printf("Layer map: %v\n", stackMap)

	total := 0

	//fmt.Printf("numbers list 1: %v\n", numbers_list)
	//fmt.Printf("numbers map: %v\n", numbers_map)
	fmt.Printf("T: %d\n", total)
}
