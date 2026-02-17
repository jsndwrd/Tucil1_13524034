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

// return ukuran board (NxN).
func (a *TArea) N() int {
	return a.n
}

// return huruf region pada (row, col).
func (a *TArea) RegionAt(row, col int) string {
	return a.cells[row*a.n+col].color
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
	
	rowClean := make([]string, 0, len(rows))
	for _, row := range rows {
		fix := strings.TrimSpace(row)
		if fix != "" {
			rowClean = append(rowClean, fix)
		}
	}
	
	// Validasi input
	if len(rowClean) == 0 || len(rowClean[0]) == 0 {
		return nil, fmt.Errorf("Input kosong.")
	}
	N := len(rowClean[0])

	if len(rowClean) != N {
		return nil, fmt.Errorf("Ukuran harus NxN.")
	}
	for _, row := range rowClean {
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

	for row := 0; row < len(rowClean); row++ {
		for col := 0; col < len(rowClean[row]); col++ {
			region := rowClean[row][col]
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

	area = &TArea{
		n:              N,
		totalColor:     len(colorCtr),
		cells:          arrCell,
		queensLocation: make([]TPosition, 0, N),
	}

	return area, nil
}

func ValidRegion(area *TArea) error {
	colorPosition := make(map[byte][]TPosition)
	for _, cell := range area.cells {
		region := byte(cell.color[0])
		colorPosition[region] = append(colorPosition[region], cell.TPosition)
	}

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
				return fmt.Errorf("Region '%c' tidak boleh terpisah.", region)
			}
		}
	}

	return nil
}

// Izin 6 warna cape nyopas hex
var CellColor = map[string]string{
	"A": "#2660A4",
	"B": "#326273",
	"C": "#E39774",
	"D": "#DF2935",
	"E": "#1B998B",
	"F": "#A6CFD5",
	"G": "#2660A4",
	"H": "#326273",
	"I": "#E39774",
	"J": "#DF2935",
	"K": "#1B998B",
	"L": "#A6CFD5",
	"M": "#2660A4",
	"N": "#326273",
	"O": "#E39774",
	"P": "#DF2935",
	"Q": "#1B998B",
	"R": "#A6CFD5",
	"S": "#2660A4",
	"T": "#326273",
	"U": "#E39774",
	"V": "#DF2935",
	"W": "#1B998B",
	"X": "#A6CFD5",
	"Y": "#2660A4",
	"Z": "#326273",
}

func (a *TArea) Color(row, col int) string {
	region := a.RegionAt(row, col)
	if color, ok := CellColor[region]; ok {
		return color
	}
	return ""
}