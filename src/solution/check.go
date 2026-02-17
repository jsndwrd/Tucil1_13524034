package solution

import "fmt"

type TPosition struct {
	row int
	col int
}

func CheckPosition(area TArea, queensLocation []TPosition) (valid bool) {
	PrintPosition(area, queensLocation)
	if len(queensLocation) == 0 {
		return true
	}

	for i := 0; i < len(queensLocation); i++ {
		p := queensLocation[i]
		
		for j := i + 1; j < len(queensLocation); j++ {
			q := queensLocation[j]

			if p.row == q.row { // Sebaris
				fmt.Printf("Solusi tidak valid!\n\n")

				return false
			}
			if p.col == q.col { // Sekolom
				fmt.Printf("Solusi tidak valid!\n\n")

				return false
			}
			if p.row-p.col == q.row-q.col { // Diagonal
				fmt.Printf("Solusi tidak valid!\n\n")

				return false
			}
			if p.row+p.col == q.row+q.col { // Diagonal
				fmt.Printf("Solusi tidak valid!\n\n")

				return false
			}
		}
	}
	fmt.Printf("Solusi valid!\n\n")
	return true
}

func PrintPosition(area TArea, queensLocation []TPosition) {
	for i := 0; i < area.n; i++ {
		for j := 0; j < area.n; j++ {
			if queensLocation[i].row == i && queensLocation[i].col == j {
				fmt.Printf("# ")
			} else {
				fmt.Printf("%s ", area.cells[i*area.n+j].color)
			}
		}
		fmt.Println()
	}
}