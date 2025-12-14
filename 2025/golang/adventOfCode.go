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

		xList = append(xList, x)
		xMap[x] = append(xMap[x], y)

		yList = append(yList, y)
		yMap[y] = append(yMap[y], x)
	}

	slices.Sort(xList)
	slices.Sort(yList)
	slices.Compact(xList) // Remove dupes
	slices.Compact(yList) // Remove dupes

	result := 0
	for _, v := range vertices {

		area := findBoxForVertex(v, xList, yList, xMap, yMap)
		result = max(area, result)
	}

	return max(result, 1) // Just in case it never finds a box, a straight line of 1 is returned.
}

func findBoxForVertex(vert []int, xList, yList []int, xMap, yMap map[int][]int) int {

	// Try every direction until we find a valid box
	for dir := DirectionUp; dir <= DirectionLeft; dir++ {

		if area, _ := search(vert, vert[0], vert[1], xList, yList, xMap, yMap, dir, [][]int{vert}); area != 0 {
			return area
		}
	}

	return 0
}

// search recursively returns the area and vertices of the largest rectangle which doesn't exit the bounds of the polygon.
// The rectangle will originate from ogVert and must contain 2 vertices that are diagonal from each other.
func search(ogVert []int, x, y int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, box [][]int) (int, [][]int) {
	dir = dir % 4

	step := len(box)
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

				if area, bb := moveVertical(ogVert, x, yK, xList, yList, xMap, yMap, dir+1, box); area != 0 {
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

				if area, bb := moveHorizontal(ogVert, xK, y, xList, yList, xMap, yMap, dir+1, box); area != 0 {
					return area, bb
				}
			}

		}

	case 3: // Only consider (xK, yK) where xK==xOg for horizontal or yK==yOg for vertical

		switch dir {
		case DirectionUp:
		case DirectionDown:

			if area, bb := moveVertical(ogVert, x, ogVert[1], xList, yList, xMap, yMap, dir+1, box); area != 0 {
				return area, bb
			}

		case DirectionRight:
		case DirectionLeft:

			if area, bb := moveHorizontal(ogVert, ogVert[0], y, xList, yList, xMap, yMap, dir+1, box); area != 0 {
				return area, bb
			}

		}

	case 4: // OgVert is only vertex left and must be valid
		box = append(box, ogVert)

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

// moveVertical checks if the vertex at (Xi, Yk) can be visited vertically (up or down) without escaping the bounds of the polygon. If so, it visits that vertex.
func moveVertical(ogVert []int, xI, yK int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, box [][]int) (int, [][]int) {
	var (
		yMin = yList[0]
		yMax = yList[len(yList)-1]
	)

	inBounds := yK >= yMin && yK <= yMax
	withinPolygon := xI >= slices.Min(yMap[yK]) && xI <= slices.Max(yMap[yK])
	if inBounds && withinPolygon {
		return search(ogVert, xI, yK, xList, yList, xMap, yMap, dir, append(box, []int{xI, yK}))
	}

	return 0, box
}

// moveHorizontal checks if the vertex at (Xk, Yi) can be visited horizontally (left or right) without escaping the bounds of the polygon. If so, it visits that vertex.
func moveHorizontal(ogVert []int, xK, yI int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, box [][]int) (int, [][]int) {
	var (
		xMin = xList[0]
		xMax = xList[len(xList)-1]
	)

	inBounds := xK >= xMin && xK <= xMax
	withinPolygon := yI >= slices.Min(xMap[xK]) && yI <= slices.Max(xMap[xK])
	if inBounds && withinPolygon {
		return search(ogVert, xK, yI, xList, yList, xMap, yMap, dir, append(box, []int{xK, yI}))
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
