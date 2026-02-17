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
		for i = n-1; i >= 0; i-- {
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

func TryPositionOptimized(area *TArea, onStep func([]TPosition) bool) []TPosition {
	n := area.n

	cols := make([]int, n)
	for i := 0; i < n; i++ {
		cols[i] = i
	}

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

		i := n-2
		for i >= 0 && cols[i] > cols[i+1] {
			i--
		}
		if i < 0 {
			break
		}

		j := n-1
		for cols[j] < cols[i] {
			j--
		}

		cols[i], cols[j] = cols[j], cols[i]

		for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
			cols[l], cols[r] = cols[r], cols[l]
		}
	}

	return nil
}

func FindPosition(area *TArea) (queensLocation []TPosition) {
	return TryPosition(area, func([]TPosition) bool {
		return true // true maka continue
	})
}