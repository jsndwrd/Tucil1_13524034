package main

import (
	"fmt"
	"queens/solution"
)

func main() {
	fmt.Printf("Queens\n")
	testInput := "AAAB\nAABB\nCCCD\nDDDD"
	area, err := solution.InputCells(testInput)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	ans := solution.TryPosition(area)
	if ans == nil {
		fmt.Printf("Tidak ada solusi.\n")
		return
	}
	fmt.Printf("%v\n",ans)
}
