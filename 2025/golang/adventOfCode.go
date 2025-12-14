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

func search(ogVert []int, x, y int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, box [][]int) (int, [][]int) {

	step := len(box)
	switch step {
	case 3:

		switch dir {
		case DirectionUp:
		case DirectionDown:
			{

				// Same as Down
				if area, bb := moveVertical(ogVert, x, ogVert[1], xList, yList, xMap, yMap, dir+1, box); area != 0 {
					return area, bb
				}

			}
		case DirectionRight:
		case DirectionLeft:
			{
				// Same as Right
				if area, bb := moveHorizontal(ogVert, ogVert[0], y, xList, yList, xMap, yMap, dir+1, box); area != 0 {
					return area, bb
				}
			}
		}

	case 4:
		// OgVert is only vertex left and must be valid
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

	default: // Steps 1 and 2

		switch dir {
		case DirectionUp:
			{
				// Starting from largest Y to smallest
				for _, yK := range slices.Backward(yList) {
					if yK == y {
						break
					}

					// Same as Down
					if area, bb := moveVertical(ogVert, x, yK, xList, yList, xMap, yMap, dir+1, box); area != 0 {
						return area, bb
					}
				}
			}
		case DirectionRight:
			{
				// Starting from largest X to smallest
				for _, xK := range slices.Backward(xList) {
					if xK == x {
						break
					}

					// Same as Right
					if area, bb := moveHorizontal(ogVert, xK, y, xList, yList, xMap, yMap, dir+1, box); area != 0 {
						return area, bb
					}
				}
			}
		case DirectionDown:
			{
				// Starting from the smallest Y to the largest
				for _, yK := range yList {
					if yK == y {
						break
					}

					// Same as Up
					if area, bb := moveVertical(ogVert, x, yK, xList, yList, xMap, yMap, dir+1, box); area != 0 {
						return area, bb
					}
				}
			}
		case DirectionLeft:
			{
				// Starting from the smallest X to the largest
				for _, xK := range xList {
					if xK == x {
						break
					}

					// Same as Left
					if area, bb := moveHorizontal(ogVert, xK, y, xList, yList, xMap, yMap, dir+1, box); area != 0 {
						return area, bb
					}
				}
			}
		}
	}

	return 0, box
}

func moveVertical(ogVert []int, xI, yK int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, box [][]int) (int, [][]int) {
	var (
		yMin = yList[0]
		yMax = yList[len(yList)-1]
	)

	inBounds := xI >= slices.Min(yMap[yK]) && xI <= slices.Max(yMap[yK])
	if (yK >= yMin && yK <= yMax) && inBounds {
		area, bb := search(ogVert, xI, yK, xList, yList, xMap, yMap, dir, append(box, []int{xI, yK}))
		return area, bb[:len(bb)-1] // Trim of the vert (xI, yK) from the box
	}

	return 0, box
}

func moveHorizontal(ogVert []int, xK, yI int, xList, yList []int, xMap, yMap map[int][]int, dir Direction, box [][]int) (int, [][]int) {
	var (
		xMin = xList[0]
		xMax = xList[len(xList)-1]
	)

	inBounds := yI >= slices.Min(xMap[xK]) && yI <= slices.Max(xMap[xK])
	if (xK >= xMin && xK <= xMax) && inBounds {
		area, bb := search(ogVert, xK, yI, xList, yList, xMap, yMap, dir, append(box, []int{xK, yI}))
		return area, bb[:len(bb)-1] // Trim of the vert (xI, yK) from the box
	}

	return 0, box
}

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
