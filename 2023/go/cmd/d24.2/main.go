package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
    gonum "github.com/shabbyrobe/go-num"
)

type Coord struct {
    x int64
    y int64
    z int64
}

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

func subtractHailstones(ha Hailstone, hb Hailstone) Hailstone {
    //return Hailstone{ha.stoneNum, ha.px - hb.px, ha.py - hb.py, ha.pz - hb.pz, ha.vx - hb.vx, ha.vy - hb.vy, ha.vz - hb.vz}
    return Hailstone{ha.stoneNum, ha.px - hb.px, ha.py - hb.py, ha.pz - hb.pz, ha.vx, ha.vy, ha.vz}
}

func findCollisionTime(hs1, hs2, hs3 Hailstone) int64 {
    /*
    x1, y1, z1, vx1, vy1, vz1 := float64(hs1.px), float64(hs1.py), float64(hs1.pz), float64(hs1.vx), float64(hs1.py), float64(hs1.pz)
    x2, y2, z2, vx2, vy2, vz2 := float64(hs2.px), float64(hs2.py), float64(hs2.pz), float64(hs2.vx), float64(hs2.py), float64(hs2.pz)
    x3, y3, z3, vx3, vy3, vz3 := float64(hs3.px), float64(hs3.py), float64(hs3.pz), float64(hs3.vx), float64(hs3.py), float64(hs3.pz)
    */
    fmt.Printf("%v\n", hs1)
    fmt.Printf("%v\n", hs2)
    fmt.Printf("%v\n", hs3)
    x1, y1, z1, vx1, vy1, vz1 := gonum.I128From64(hs1.px), gonum.I128From64(hs1.py), gonum.I128From64(hs1.pz), gonum.I128From64(hs1.vx), gonum.I128From64(hs1.vy), gonum.I128From64(hs1.vz)
    x2, y2, z2, vx2, vy2, vz2 := gonum.I128From64(hs2.px), gonum.I128From64(hs2.py), gonum.I128From64(hs2.pz), gonum.I128From64(hs2.vx), gonum.I128From64(hs2.vy), gonum.I128From64(hs2.vz)
    x3, y3, z3, vx3, vy3, vz3 := gonum.I128From64(hs3.px), gonum.I128From64(hs3.py), gonum.I128From64(hs3.pz), gonum.I128From64(hs3.vx), gonum.I128From64(hs3.vy), gonum.I128From64(hs3.vz)

    fmt.Printf("%d * %d, + %d * %d + %d * %d", y1, (z2.Sub(z3)),  y2, (z3.Sub(z1)), y3, (z1.Sub(z2)))
    yz := y1.Mul(z2.Sub(z3)).Add(y2.Mul(z3.Sub(z1))).Add(y3.Mul(z1.Sub(z2)))
    xz := x1.Mul(z3.Sub(z2)).Add(x2.Mul(z1.Sub(z3))).Add(x3.Mul(z2.Sub(z1)))
    xy := x1.Mul(y2.Sub(y3)).Add(x2.Mul(y3.Sub(y1))).Add(x3.Mul(y1.Sub(y2)))
    vxvy := vx1.Mul(vy2.Sub(vy3)).Add(vx2.Mul(vy3.Sub(vy1))).Add(vx3.Mul(vy1.Sub(vy2)))
    vxvz := vx1.Mul(vz3.Sub(vz2)).Add(vx2.Mul(vz1.Sub(vz3))).Add(vx3.Mul(vz2.Sub(vz1)))
    vyvz := vy1.Mul(vz2.Sub(vz3)).Add(vy2.Mul(vz3.Sub(vz1))).Add(vy3.Mul(vz1.Sub(vz2)))

    n := (vx2.Sub(vx3)).Mul(yz).Add((vy2.Sub(vy3)).Mul(xz)).Add((vz2.Sub(vz3)).Mul(xy))
    d := (z2.Sub(z3)).Mul(vxvy).Add((y2.Sub(y3)).Mul(vxvz)).Add((x2.Sub(x3)).Mul(vyvz))

    fmt.Printf("n d: %d %d %d\n", n, d, n.AsInt64()/d.AsInt64())

    return int64(n.AsInt64()/d.AsInt64())
}

func findRockPos(hailstoneList []Hailstone) (int64, int64, int64, bool) {
    if len(hailstoneList) < 4 {
        return 0, 0, 0, false
    }

    hs1 := subtractHailstones(hailstoneList[0], hailstoneList[0])
    hs2 := subtractHailstones(hailstoneList[1], hailstoneList[0])
    hs3 := subtractHailstones(hailstoneList[2], hailstoneList[0])

    fmt.Printf("%d\n", math.MaxInt64)

    hs1 = hailstoneList[0]
    hs2 = hailstoneList[1]
    hs3 = hailstoneList[2]

    t1 := findCollisionTime(hs1, hs2, hs3)
    t2 := findCollisionTime(hs2, hs1, hs3)
    fmt.Printf("t1 t2: %d %d\n", t1, t2)

    c1x := hs1.px + t1*hs1.vx
    c1y := hs1.py + t1*hs1.vy
    c1z := hs1.pz + t1*hs1.vz
    fmt.Printf("c1: %d %d %d\n", c1x, c1y, c1z)
    c2x := hs2.px + t2*hs2.vx
    c2y := hs2.py + t2*hs2.vy
    c2z := hs2.pz + t2*hs2.vz
    fmt.Printf("c2: %d %d %d\n", c2x, c2y, c2z)

    vx := (c2x - c1x) / (t2 - t1)
    vy := (c2y - c1y) / (t2 - t1)
    vz := (c2z - c1z) / (t2 - t1)
    fmt.Printf("v: %d %d %d\n", vx, vy, vz)

    px := hs1.px + hs1.vx*t1 - vx*t1
    py := hs1.py + hs1.vy*t1 - vy*t1
    pz := hs1.pz + hs1.vz*t1 - vz*t1
    fmt.Printf("p: %d %d %d\n", px, py, pz)

    //return px+hailstoneList[0].px, py+hailstoneList[0].py, pz+hailstoneList[0].pz, true
    return px, py, pz, true
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

    totalP2 := int64(0);
    px, py, pz, rockFound := findRockPos(hailstoneList)
    if rockFound {
        totalP2 = px + py + pz
    }

	//fmt.Printf("Bricklist: %v\n", hailstoneList)

	//fmt.Printf("numbers list 1: %v\n", numbers_list)
	//fmt.Printf("numbers map: %v\n", numbers_map)
	fmt.Printf("T: %d %d\n", total, totalP2)
}
