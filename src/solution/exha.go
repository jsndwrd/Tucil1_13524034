package solution

func OneQueen(area *TArea, queens []TPosition) bool {
	if area == nil {
		return false
	}

	regionCount := make(map[string]int, area.totalColor)
	for _, cell := range area.cells {
		regionCount[cell.color] = 0
	}

	// Hitung queen per region
	for _, q := range queens {
		cell := area.cells[q.row*area.n+q.col]
		regionCount[cell.color]++
	}

	for _, c := range regionCount {
		if c != 1 {
			return false
		}
	}

	return true
}

func TryPosition(area *TArea, onStep func([]TPosition) bool) []TPosition {
	n := area.n
	cols := make([]int, n)
	for {
		temp := make([]TPosition, n)
		for row := 0; row < n; row++ {
			temp[row] = TPosition{row, cols[row]}
		}

		// callback update GUI
		if !onStep(temp) {
			return nil
		}

		if CheckPosition(*area, temp) && OneQueen(area, temp) {
			area.queensLocation = temp
			return temp
		}

		var i int
		for i = n - 1; i >= 0; i-- {
			cols[i]++
			if cols[i] < n {
				break
			}
			cols[i] = 0
		}
		if i < 0 {
			break
		}
	}

	return nil
}

func FindPosition(area *TArea) (queensLocation []TPosition) {
	return TryPosition(area, func([]TPosition) bool {
		return true // true maka continue
	})
}