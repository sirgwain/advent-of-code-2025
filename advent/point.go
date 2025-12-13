package advent

import (
	"fmt"
	"slices"
)

type Point struct {
	X int
	Y int
}
type PositionDirection struct {
	Point
	Direction Direction
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p Point) addDirection(dir Direction) Point {
	x, y := dir.OffsetMultiplier()
	return Point{p.X + x, p.Y + y}
}

func (p Point) String() string {
	return fmt.Sprintf("%d, %d", p.X, p.Y)
}

func area(p1, p2 Point) int {
	dx := absInt(p1.X-p2.X) + 1
	dy := absInt((p2.Y - p1.Y)) + 1
	return int(dx * dy)
}

// pointInPolygon returns true if p is inside poly.
func pointInPolygon(p Point, poly []Point) bool {
	inside := false

	n := len(poly)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		pi := poly[i]
		pj := poly[j]

		if pointOnSegment(p, pi, pj) {
			return true
		}

		intersects := ((pi.Y > p.Y) != (pj.Y > p.Y)) &&
			(p.X < (pj.X-pi.X)*(p.Y-pi.Y)/(pj.Y-pi.Y)+pi.X)

		if intersects {
			inside = !inside
		}
	}

	return inside
}

func pointOnSegment(p, a, b Point) bool {
	cross := (p.Y-a.Y)*(b.X-a.X) - (p.X-a.X)*(b.Y-a.Y)
	if cross != 0 {
		return false
	}
	dot := (p.X-a.X)*(b.X-a.X) + (p.Y-a.Y)*(b.Y-a.Y)
	if dot < 0 {
		return false
	}
	lengthSq := (b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y)
	return dot <= lengthSq
}

// On a horizontal "row" y-band between y and y+1,
// test if [x1, x2] is completely inside the polygon.
func rowInside(poly []Point, y, x1, x2 int) bool {
	if x1 > x2 {
		x1, x2 = x2, x1
	}

	var xs []int
	n := len(poly)

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		a := poly[j]
		b := poly[i]

		// Only vertical edges matter for a horizontal scan.
		if a.X != b.X {
			continue
		}

		// Normalize so a.Y <= b.Y
		if a.Y > b.Y {
			a, b = b, a
		}

		// We conceptually sample at y+0.5.
		// That lies in the band [a.Y, b.Y) exactly when:
		// a.Y <= y+0.5 < b.Y  <=>  a.Y <= y and y+1 <= b.Y
		if !(a.Y <= y && y+1 <= b.Y) {
			continue
		}

		xs = append(xs, a.X)
	}

	if len(xs) == 0 {
		// This row never enters the polygon.
		return false
	}

	slices.Sort(xs)

	inside := false
	prevX := 0

	// We’ll build the union of inside intervals as we sweep.
	for _, x := range xs {
		if !inside {
			// Entering polygon at x
			prevX = x
			inside = true
		} else {
			// Leaving polygon at x, interior is [prevX, x]
			// For axis-aligned / integer, treat boundary as inside.
			// Check if our whole [x1, x2] fits in this interior run.
			if x1 >= prevX && x2 <= x {
				return true
			}
			inside = false
		}
	}

	// Should end outside for a proper closed polygon.
	return false
}
