package main

import (
	"fmt"
	"math"
	"slices"
)

func main() {
	fmt.Println("Weclome to advent 2025")

	// vertices := [][]int{...}
	// FindMaxBoundingBox(vertices)
}

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

func FindMaxBoundingBox(vertices [][]int) int {
	// Setup lists and maps
	var xList, yList []int
	yMap := make(map[int][]int)
	xMap := make(map[int][]int)
	for _, v := range vertices {
		x := v[0]
		y := v[1]

		// Initialise arrays within the maps
		if _, ok := xMap[x]; !ok {
			xMap[x] = make([]int, 0)
		}

		if _, ok := yMap[y]; !ok {
			yMap[y] = make([]int, 0)
		}

		xList = append(xList, x)
		yList = append(yList, y)

		xMap[x] = append(xMap[x], y)
		yMap[y] = append(yMap[y], x)
	}

	slices.Sort(xList)
	slices.Sort(yList)
	slices.Compact(xList) // Remove dupes
	slices.Compact(yList) // Remove dupes

	maxArea := 0
	seenCorners := make(map[int]bool)
	for _, v := range vertices {

		area := findBoxForVertex(v, xList, yList, xMap, yMap, seenCorners)
		maxArea = max(area, maxArea)
	}

	return max(maxArea, 1) // Just in case it never finds a box, a straight line of 1 is returned.
}

func findBoxForVertex(vert []int, xList, yList []int, xMap, yMap map[int][]int, seenCorners map[int]bool) int {

	// Try every direction until we find a valid box
	// We should only need to find 1 box as the one we find for this vertex should be the largest valid one.
	// Any other box which contains this vertex must contain at least another vertex and will be found during that vertex's search.
	for dir := DirectionUp; dir <= DirectionLeft; dir++ {

		if area, _ := search(vert, vert[0], vert[1], xList, yList, xMap, yMap, dir, seenCorners, [][]int{vert}); area != 0 {
			return area
		}
	}

	return 0
}

// search recursively returns the area and vertices of the largest rectangle which doesn't exit the bounds of the polygon.
// The rectangle will originate from ogVert and must contain 2 vertices that are diagonal from each other.
func search(ogVert []int, x, y int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, seenCorners map[int]bool, box [][]int) (int, [][]int) {
	dir = dir % 4

	step := len(box)

	// Stop here if this corner has been seen. A seen corner only has one possible box
	if step >= 2 {
		hash := hashCorner([]int{x, y}, box[step-1], box[step-2])
		if seenCorners[hash] {
			return 0, box
		}
	}

	switch step {

	case 1, 2: // Try the farthest valid Xk or Yk

		switch dir {

		case DirectionUp:
		case DirectionDown:

			iterator := slices.All(yList) // iterate forwards (smallest to largest)
			if dir == DirectionUp {
				iterator = slices.Backward(yList) // iterate backwards (largest to smallest)
			}

			for _, yK := range iterator {
				if yK == y {
					break
				}

				if area, bb := moveVertical(ogVert, x, yK, xList, yList, xMap, yMap, dir+1, seenCorners, box); area != 0 {
					return area, bb
				}
			}

		case DirectionLeft:
		case DirectionRight:

			iterator := slices.All(xList) // iterate forwards (smallest to largest)
			if dir == DirectionRight {
				iterator = slices.Backward(xList) // iterate backwards (largest to smallest)
			}

			for _, xK := range iterator {
				if xK == x {
					break
				}

				if area, bb := moveHorizontal(ogVert, xK, y, xList, yList, xMap, yMap, dir+1, seenCorners, box); area != 0 {
					return area, bb
				}
			}

		}

	case 3: // Only consider (xK, yK) where xK==xOg for horizontal or yK==yOg for vertical

		switch dir {
		case DirectionUp:
		case DirectionDown:

			if area, bb := moveVertical(ogVert, x, ogVert[1], xList, yList, xMap, yMap, dir+1, seenCorners, box); area != 0 {
				return area, bb
			}

		case DirectionRight:
		case DirectionLeft:

			if area, bb := moveHorizontal(ogVert, ogVert[0], y, xList, yList, xMap, yMap, dir+1, seenCorners, box); area != 0 {
				return area, bb
			}

		}

	case 4: // Box complete

		// Store the corners that have been seen in this box
		seenCorners[hashCorner(box[0], box[1], box[2])] = true
		seenCorners[hashCorner(box[1], box[2], box[3])] = true
		seenCorners[hashCorner(box[2], box[3], box[0])] = true
		seenCorners[hashCorner(box[3], box[0], box[1])] = true

		// box[0] and box[2] are diagonals and box[1] and box[3] are diagonals
		// Box[0] is already a vertex on the polygon. Only need to check box[2] || (box[1] && box[3])
		if isOnPoligon(box[2], xMap, yMap) || (isOnPoligon(box[1], xMap, yMap) && isOnPoligon(box[3], xMap, yMap)) {
			// Use orVert and box[2] vertices (diagonal) for area calc
			x1, y1 := ogVert[0], ogVert[1]
			x2, y2 := box[2][0], box[2][1]

			length := float64(x2 - x1)
			height := float64(y2 - y1)
			return (int)(math.Abs(length * height)), nil // Return nil for the box cuz we don't care about it anymore.
		}
	}

	return 0, box
}

// from https://stackoverflow.com/questions/3934100/good-hash-function-for-list-of-2-d-positions
func hashCorner(v1, v2, v3 []int) int {
	hash := 17
	tmp := 0
	for _, v := range [][]int{v1, v2, v3} {
		tmp = (hash + hashVert(v))
		hash = (tmp << 5) - tmp
	}

	return hash
}

// from https://stackoverflow.com/questions/3934100/good-hash-function-for-list-of-2-d-positions
func hashVert(v []int) int {
	x, y := v[0], v[1]

	hash := 17
	hash = ((hash + x) << 5) - (hash + x)
	hash = ((hash + y) << 5) - (hash + y)
	return hash
}

// moveVertical checks if the vertex at (Xi, Yk) can be visited vertically (up or down) without escaping the bounds of the polygon. If so, it visits that vertex.
func moveVertical(ogVert []int, xI, yK int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, seenCorners map[int]bool, box [][]int) (int, [][]int) {

	inBounds := yK >= yList[0] && yK <= yList[len(yList)-1]                   // Ymin <= Yk <= Ymax (feels redundant)
	withinPolygon := xI >= slices.Min(yMap[yK]) && xI <= slices.Max(yMap[yK]) // min(x in Yk) <= Xi <= max(x in Yk)
	if inBounds && withinPolygon {
		return search(ogVert, xI, yK, xList, yList, xMap, yMap, dir, seenCorners, append(box, []int{xI, yK}))
	}

	return 0, box
}

// moveHorizontal checks if the vertex at (Xk, Yi) can be visited horizontally (left or right) without escaping the bounds of the polygon. If so, it visits that vertex.
func moveHorizontal(ogVert []int, xK, yI int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, seenCorners map[int]bool, box [][]int) (int, [][]int) {

	inBounds := xK >= xList[0] && xK <= xList[len(xList)-1]                   // Xmin <= Xk <= Xmax (feels redundant)
	withinPolygon := yI >= slices.Min(xMap[xK]) && yI <= slices.Max(xMap[xK]) // min(y in Xk) <= Yi <= max(y in Xk)
	if inBounds && withinPolygon {
		return search(ogVert, xK, yI, xList, yList, xMap, yMap, dir, seenCorners, append(box, []int{xK, yI}))
	}

	return 0, box
}

// isOnPolugon determins if the given vertex is a polgon vertex
func isOnPoligon(vert []int, xMap, yMap map[int][]int) bool {
	x := vert[0]
	y := vert[1]

	if xValues, ok := yMap[y]; ok {
		if slices.Index(xValues, x) != -1 {
			return true
		}
	}

	if yValues, ok := xMap[x]; ok {
		if slices.Index(yValues, y) != -1 {
			return true
		}
	}

	return false
}
