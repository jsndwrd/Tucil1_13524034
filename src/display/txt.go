package display

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"queens/solution"
)

func SaveOutput(area *solution.TArea, queens []solution.TPosition, inputLayout string) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Gagal mendapatkan directory: %v", err)
	}

	moduleRoot := workDir
	for {
		if _, statErr := os.Stat(filepath.Join(moduleRoot, "go.mod")); statErr == nil {
			break
		}
		parent := filepath.Dir(moduleRoot)
		if parent == moduleRoot {
			moduleRoot = workDir
			break
		}
		moduleRoot = parent
	}

	projectRoot := filepath.Dir(moduleRoot)
	directory := filepath.Join(projectRoot, "test")
	if mkErr := os.MkdirAll(directory, 0o755); mkErr != nil {
		return "", fmt.Errorf("Gagal membuat folder test: %v", mkErr)
	}

	timestamp := time.Now().Format("20260218_130415")
	outputName := filepath.Join(directory, fmt.Sprintf("output_%s.txt", timestamp))
	
	var output strings.Builder
	
	if queens == nil {
		output.WriteString("Tidak ada solusi\n")
	} else {
		n := area.N()
		queenMap := make(map[string]bool)
		for _, q := range queens {
			key := fmt.Sprintf("%d,%d", q.Row(), q.Col())
			queenMap[key] = true
		}
		
		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				key := fmt.Sprintf("%d,%d", row, col)
				if queenMap[key] {
					output.WriteString("#")
				} else {
					output.WriteString(area.RegionAt(row, col))
				}
			}
			if row < n-1 {
				output.WriteString("\n")
			}
		}
	}

	if err := os.WriteFile(outputName, []byte(output.String()), 0o644); err != nil { // save file dan handling error
		return "", fmt.Errorf("Gagal menyimpan output: %v", err)
	}

	return outputName, nil
}
