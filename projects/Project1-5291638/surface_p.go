package main

import (
	"flag"
	"math"
	"runtime"
	"sync"
)

const (
	width, height = 600, 320
	cells         = 12000
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	workers := flag.Int("g", 4, "Number of goroutines")
	flag.Parse()
	//log.Printf("Number of goroutines: %d", *workers)

	runtime.GOMAXPROCS(*workers)

	var wg sync.WaitGroup

	rows := cells / *workers
	//fmt.Printf("# of go %d", rows)

	for k := 0; k < *workers; k++ {

		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			start := id * rows
			end := start + rows

			for i := start; i < end; i++ {
				for j := 0; j < cells; j++ {
					corner(i+1, j)
					corner(i, j)
					corner(i, j+1)
					corner(i+1, j+1)

				}
			}

		}(k)
	}

	wg.Wait()
}

func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	// Compute surface height z.
	z := f(x, y)
	// Project (x,y,z) isometrically onto 2D 	SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}
func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}
