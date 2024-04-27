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

func tryDestroy(brickList []Brick, stackMap map[int][]int, initialBrickIndex int) int {
    //currentBrick := brickList[initialBrickIndex]
    fallingBrickIndexes := []int{ initialBrickIndex } 
    for currentFallingI := 0; currentFallingI < len(fallingBrickIndexes); currentFallingI += 1 {
        currentBrickIndex := fallingBrickIndexes[currentFallingI]
        currentBrick := brickList[currentBrickIndex]

        for _, aboveBrickIndex := range stackMap[currentBrick.posEnd.z+1] {
            aboveBrick := brickList[aboveBrickIndex]
            if aboveBrick.posStart.z < (currentBrick.posEnd.z + 1) {
                // Can't be resting on the current brick
                continue
            }

            isSupported := false

            // Test all below bricks to see if supported by other bricks
            aboveBrickTest := aboveBrick
            aboveBrickTest.posStart.z -= 1
            aboveBrickTest.posEnd.z -= 1
            candidateSupportingLoop:
            for _, candidateSupportingIndex := range stackMap[currentBrick.posEnd.z] {
                for _,fallingIndex := range fallingBrickIndexes {
                    if candidateSupportingIndex == fallingIndex {
                        // Don't test against a falling brick
                        continue candidateSupportingLoop
                    }
                }
                candidateSupporting := brickList[candidateSupportingIndex]
                for testX := aboveBrickTest.posStart.x; testX <= aboveBrickTest.posEnd.x; testX += 1 {
                    for testY := aboveBrickTest.posStart.y; testY <= aboveBrickTest.posEnd.y; testY += 1 {
                        intersectsX := candidateSupporting.posStart.x <= testX && candidateSupporting.posEnd.x >= testX
                        intersectsY := candidateSupporting.posStart.y <= testY && candidateSupporting.posEnd.y >= testY
                        intersectsZ := candidateSupporting.posStart.z <= aboveBrickTest.posStart.z && candidateSupporting.posEnd.z >= aboveBrickTest.posStart.z
                        if intersectsX && intersectsY && intersectsZ {
                            //fmt.Printf("             %d %d %d intersects | %t %t %t\n", testX, testY, testBrick.posStart.z, intersectsX, intersectsY, intersectsZ)
                            isSupported = true
                            break
                        } else {
                            //fmt.Printf("             %d %d %d not intersect | %t %t %t\n", testX, testY, testBrick.posStart.z, intersectsX, intersectsY, intersectsZ)
                        }
                    }
                    if isSupported {
                        break
                    }
                }
            }
            if !isSupported {
                currentIsFalling := false
                for _,tempFalling := range fallingBrickIndexes {
                    if aboveBrickIndex == tempFalling {
                        currentIsFalling = true
                        break
                    }
                }
                if !currentIsFalling {
                    fallingBrickIndexes = append(fallingBrickIndexes, aboveBrickIndex)
                }
            }
        }
    }
    //fmt.Printf("Falling list %d => %d | %v\n", initialBrickIndex, len(fallingBrickIndexes), fallingBrickIndexes)
    return len(fallingBrickIndexes) - 1
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

	sort.Ints(stackHeightList)

	// Do gravity.
	for currentLayerNum := stackHeightList[0]; currentLayerNum < len(stackHeightList); currentLayerNum += 1 {
		currentLayer := stackHeightList[currentLayerNum]
		currentLayerIndexes := stackMap[currentLayer]
		fmt.Printf("Handling layer %d (%d) | %v | %d\n", currentLayer, currentLayerNum, currentLayerIndexes, lastLayer)
		for _, currentLayerIndex := range currentLayerIndexes {
			currentBrick := brickList[currentLayerIndex]
			//brickHeight := currentBrick.posEnd.z - currentBrick.posStart.z

			fmt.Printf("    Current brick: %v\n", currentBrick)
			//fmt.Printf("%d, %v\n", lastLayer, brickHeight)

			testBrick := currentBrick
			for testBrick.posStart.z > 1 {
				testBrick.posStart.z -= 1
				testBrick.posEnd.z -= 1

				hasIntersection := false
				testLayerIndexes := stackMap[testBrick.posStart.z]
				//fmt.Printf("        Test brick: %v | vs %v | %v\n", testBrick, testLayerIndexes, stackMap)
				for _, comparsionBrickIndex := range testLayerIndexes {
					comparisonBrick := brickList[comparsionBrickIndex]
					//fmt.Printf("            comparison brick: %v\n", comparisonBrick)
					for testX := testBrick.posStart.x; testX <= testBrick.posEnd.x; testX += 1 {
						for testY := testBrick.posStart.y; testY <= testBrick.posEnd.y; testY += 1 {
							intersectsX := comparisonBrick.posStart.x <= testX && comparisonBrick.posEnd.x >= testX
							intersectsY := comparisonBrick.posStart.y <= testY && comparisonBrick.posEnd.y >= testY
							intersectsZ := comparisonBrick.posStart.z <= testBrick.posStart.z && comparisonBrick.posEnd.z >= testBrick.posStart.z
							if intersectsX && intersectsY && intersectsZ {
								//fmt.Printf("             %d %d %d intersects | %t %t %t\n", testX, testY, testBrick.posStart.z, intersectsX, intersectsY, intersectsZ)
								hasIntersection = true
								break
							} else {
								//fmt.Printf("             %d %d %d not intersect | %t %t %t\n", testX, testY, testBrick.posStart.z, intersectsX, intersectsY, intersectsZ)
							}
						}
						if hasIntersection {
							break
						}
					}
					if hasIntersection {
						break
					}
				}

				if hasIntersection {
					testBrick.posStart.z += 1
					testBrick.posEnd.z += 1
					//fmt.Printf("Test brick end with intersection: %v\n", testBrick)
					break
				}
			}
			newBrick := testBrick
			//fmt.Printf("New brick: %v\n", newBrick)
			for oldLayerIndex := currentBrick.posStart.z; oldLayerIndex <= currentBrick.posEnd.z; oldLayerIndex += 1 {
				newStackHeightItems := []int{}
				for _, oldStackHeightItem := range stackMap[oldLayerIndex] {
					if oldStackHeightItem != currentBrick.num {
						newStackHeightItems = append(newStackHeightItems, oldStackHeightItem)
					}
				}
				stackMap[oldLayerIndex] = newStackHeightItems
			}
			for newLayerIndex := newBrick.posStart.z; newLayerIndex <= newBrick.posEnd.z; newLayerIndex += 1 {
				//fmt.Printf("New brick layer: %d\n", newLayerIndex)
				newStackMap, stackMapFound := stackMap[newLayerIndex]
				if !stackMapFound {
					newStackMap = []int{}
					stackHeightList = append(stackHeightList, newLayerIndex)
					sort.Ints(stackHeightList)
				}
				newStackMap = append(newStackMap, newBrick.num)
				stackMap[newLayerIndex] = newStackMap
			}
			brickList[newBrick.num] = newBrick
		}
	}

	sort.Ints(stackHeightList)

	total := 0
    totalFalling := 0

	for _, currentBrick := range brickList {
        countOtherFalling := tryDestroy(brickList, stackMap, currentBrick.num)
        //fmt.Printf("Falling %v => %d\n", currentBrick, countOtherFalling)
        if countOtherFalling == 0 {
            total += 1
        } else {
            totalFalling += countOtherFalling
        }
	}

	fmt.Printf("Layers: %v\n", stackHeightList)
	fmt.Printf("Layer map: %v\n", stackMap)
	fmt.Printf("Bricklist: %v\n", brickList)

	//fmt.Printf("numbers list 1: %v\n", numbers_list)
	//fmt.Printf("numbers map: %v\n", numbers_map)
	fmt.Printf("T: %d, %d\n", total, totalFalling)
}
