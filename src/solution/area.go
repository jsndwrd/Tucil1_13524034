package solution

import (
	"fmt"
	"strings"
	"unicode"
)


type TCell struct {
	TPosition
	color string // Representasi kode unik tiap warna
}

type TArea struct {
	n              int         // NxN
	totalColor     int         // Jumlah warna / representasi area unik
	cells []TCell
	queensLocation []TPosition // Koordinat tiap penempatan ratu
}

func checkAdjacent(first, second TPosition) bool {
	row := first.row-second.row
	if row < 0 {
		row = -row
	}
	col := first.col-second.col
	if col < 0 {
		col = -col
	}
	return (row == 1 && col == 0) || (row == 0 && col == 1)
}

func InputCells(cells string) (area *TArea, err error) {
	cells = strings.TrimSpace(cells)
	rows := strings.Split(cells, "\n")
	// Validasi input
	if len(rows) == 0 || len(rows[0]) == 0 {
		return nil, fmt.Errorf("Input tidak boleh kosong.")
	}
	N := len(rows[0])

	if len(rows) != N {
		return nil, fmt.Errorf("Ukuran harus NxN.")
	}
	for _, row := range rows {
		if len(row) != N {
			return nil, fmt.Errorf("Ukuran harus NxN.")
		}
		for _, char := range row {
			if !unicode.IsLetter(char) { // alphabet only, input dengan queen tidak wajib
				return nil, fmt.Errorf("Tiap sel hanya dapat direpresentasikan oleh alfabet.")
			}
		}
	}

	arrCell := make([]TCell, 0, N*N)
	colorCtr := make(map[string]bool)
	colorPosition := make(map[byte][]TPosition)

	for row := 0; row < len(rows); row++ {
		for col := 0; col < len(rows[row]); col++ {
			region := rows[row][col]
			color := string(region)

			cell := TCell{
				TPosition: TPosition{row, col},
				color:     color,
			}
			arrCell = append(arrCell, cell)
			colorCtr[color] = true
			colorPosition[region] = append(colorPosition[region], TPosition{row, col})
		}
	}
	if len(colorCtr) != N {
		return nil, fmt.Errorf("Jumlah region unik harus sama dengan N.")
	}

	// Region harus tidak terputus sama sekali
	for region, position := range colorPosition {
		if len(position) == 0 {
			continue
		}

		connect := make([]bool, len(position))
		connect[0] = true

		changed := true
		for changed { // while(changed)
			changed = false
			for i := 0; i < len(position); i++ {
				if connect[i] {
					continue
				}

				for j := 0; j < len(position); j++ {
					if !connect[j] {
						continue
					}
					if checkAdjacent(position[i], position[j]) {
						connect[i] = true
						changed = true
						break
					}
				}
			}
		}

		for i := 0; i < len(position); i++ {
			if !connect[i] {
				return nil, fmt.Errorf("Region '%c' tidak boleh terpisah.", region)
			}
		}
	}

	area = &TArea{
		n:              N,
		totalColor:     len(colorCtr),
		cells:          arrCell,
		queensLocation: make([]TPosition, 0, N),
	}

	return area, nil
}