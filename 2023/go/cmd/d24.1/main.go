package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Hailstone struct {
    stoneNum int
    px int64
    py int64
    pz int64
    vx int64
    vy int64
    vz int64
}

type Interaction struct {
    hs1 int
    hs2 int
}

func interactionPoint(h1 Hailstone, h2 Hailstone) (float64, float64, bool) {

    dx := float64(h2.px) - float64(h1.px)
    dy := float64(h2.py) - float64(h1.py)
    det := float64(h2.vx) * float64(h1.vy) - float64(h2.vy) * float64(h1.vx)

    if det == 0 {
        return -1, -1, false
    }

    u := (dy * float64(h2.vx) - dx * float64(h2.vy)) / det
    v := (dy * float64(h1.vx) - dx * float64(h1.vy)) / det

    if u < 0 || v < 0 {
        return -1, -1, false
    }

    // Another point.
    h1NewX := float64(h1.px) + float64(h1.vx)
    h1NewY := float64(h1.py) + float64(h1.vy)
    h2NewX := float64(h2.px) + float64(h2.vx)
    h2NewY := float64(h2.py) + float64(h2.vy)

    // Slopes of lines
    m1 := (h1NewY - float64(h1.py)) / (h1NewX - float64(h1.px))
    m2 := (h2NewY - float64(h2.py)) / (h2NewX - float64(h2.px))

    // X intercept
    b1 := float64(h1.py) - m1 * float64(h1.px)
    // Y Intercept
    b2 := float64(h2.py) - m2 * float64(h2.px)

    if m1 == m2 {
        return -1, -1, false
    }

    px := (b2 - b1) / (m1 - m2)
    py := m1 * px + b1


    return float64(px), float64(py), true
}

func countIntersections(hailstoneList []Hailstone, rangeStart int64, rangeEnd int64) int {
    interactions := []Interaction{}
    for i := 0; i < len(hailstoneList); i += 1 {
        h1 := hailstoneList[i]
        for j := i+1; j < len(hailstoneList); j += 1 {
            h2 := hailstoneList[j]

            interX, interY, hasInter := interactionPoint(h1, h2)
            if hasInter && interX >= float64(rangeStart) && interX <= float64(rangeEnd) && interY >= float64(rangeStart) && interY <= float64(rangeEnd) {
                interactions = append(interactions, Interaction{i, j})
                //fmt.Printf("X inter %v %v @ %f\n", h1, h2, interX)
            }
        }
    }

    for _,inter := range interactions {
        fmt.Printf("I: %v\n", inter)
    }

    return len(interactions)
}

func main() {
	var filename = flag.String("f", "../inputs/d24.sample1.txt", "file to use")
	//var is_verbose = flag.Bool("v", false, "verbose")
    var rangeStart = flag.Int64("rs", 7, "Range start")
    var rangeEnd = flag.Int64("re", 27, "Range end")
	flag.Parse()
	dat, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("unable to read file: %f", err)
	}

	pBrick := regexp.MustCompile("(-?[0-9]+), (-?[0-9]+), (-?[0-9]+) @ +(-?[0-9]+), +(-?[0-9]+), +(-?[0-9]+)")

	hailstoneList := []Hailstone{}

	for _, fLine := range strings.Split(strings.ReplaceAll(strings.Trim(string(dat), " \r\n"), "\r\n", "\n"), "\n") {
		lineParts := pBrick.FindStringSubmatch(fLine)
		if len(lineParts) != 7 {
			fmt.Printf("Incorrect line %s %v\n", fLine, lineParts)
			continue
		}
        //fmt.Printf("Correct line %s %v\n", fLine, lineParts)

		px, _ := strconv.ParseInt(lineParts[1], 10, 64)
		py, _ := strconv.ParseInt(lineParts[2], 10, 64)
		pz, _ := strconv.ParseInt(lineParts[3], 10, 64)
		vx, _ := strconv.ParseInt(lineParts[4], 10, 64)
		vy, _ := strconv.ParseInt(lineParts[5], 10, 64)
		vz, _ := strconv.ParseInt(lineParts[6], 10, 64)

		hailstoneNum := len(hailstoneList)

        targetHailstone := Hailstone{hailstoneNum, px, py, pz, vx, vy, vz}

		hailstoneList = append(hailstoneList, targetHailstone)
	}

    total := countIntersections(hailstoneList, *rangeStart, *rangeEnd)

	//fmt.Printf("Bricklist: %v\n", hailstoneList)

	//fmt.Printf("numbers list 1: %v\n", numbers_list)
	//fmt.Printf("numbers map: %v\n", numbers_map)
	fmt.Printf("T: %d\n", total)
}
