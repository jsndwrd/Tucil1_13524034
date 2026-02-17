package display

import (
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"queens/solution"
)

func DisplayUI() {
	a := app.New()
	w := a.NewWindow("Queens")
	w.Resize(fyne.NewSize(1080,720))

	var gridCells []*widget.Select
	var gridBg []*canvas.Rectangle
	var gridN int
	var boardArea *fyne.Container
	var currentArea *solution.TArea // area untuk visualisasi
	var updateGrid func(layout string) (*solution.TArea, error)
	var stopSolve chan struct{} // channel signal

	inputTextFunc := NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.NewInformation("Error", err.Error(), w).Show()
			return
		}
		if closer == nil {
			return
		}
		defer closer.Close()

		b, readErr := io.ReadAll(closer)
		if readErr != nil {
			dialog.NewInformation("Error", readErr.Error(), w).Show()
			return
		}
		content := strings.TrimSpace(string(b))
		if content == "" {
			dialog.NewInformation("Error", "File kosong.", w).Show()
			return
		}

		// txt to grid
		if _, parseErr := updateGrid(content); parseErr != nil {
			dialog.NewInformation("Error", parseErr.Error(), w).Show()
			return
		}
	}, w)
	inputTextButton := widget.NewButton(".txt", inputTextFunc.Show)
	
	nEntry := widget.NewEntry()
	nEntry.SetPlaceHolder("Ukuran board (N)")

	gridArea := container.NewMax()

	showNLetters := func(n int) []string {
		opts := make([]string, 0, n)
		for alp := 'A'; alp <= 'Z' && len(opts) < n; alp++ {
			opts = append(opts, string(alp))
		}
		return opts
	}

	// hex color string ke color.Color
	HexToRgb := func(hex string) color.Color {
		hex = strings.TrimPrefix(hex, "#")
		var row, g, b uint8
		fmt.Sscanf(hex, "%02x%02x%02x", &row, &g, &b)
		return &color.NRGBA{R: row, G: g, B: b, A: 255}
	}

	buildSelectGrid := func(n int, options []string) {
		if n <= 0 {
			return
		}
		cells := make([]*widget.Select, 0, n*n)
		rects := make([]*canvas.Rectangle, 0, n*n)
		grid := container.NewGridWithColumns(n)
		for i := 0; i < n*n; i++ {
			cell := widget.NewSelect(options, nil)
			cell.PlaceHolder = "-"
			
			// background color
			bgRect := canvas.NewRectangle(&color.NRGBA{R: 255, G: 255, B: 255, A: 255})
			bgRect.SetMinSize(fyne.NewSize(50, 30))
			rects = append(rects, bgRect)
		
			j := i
			cell.OnChanged = func(selected string) {
				if selected != "" {
					if hexColor, ok := solution.CellColor[selected]; ok {
						rects[j].FillColor = HexToRgb(hexColor)
					} else {
						rects[j].FillColor = &color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					}
					rects[j].Refresh()
				} else {
					rects[j].FillColor = &color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					rects[j].Refresh()
				}
			}
			
			// Wrap Select dengan background
			cellContainer := container.NewStack(bgRect, cell)
			cells = append(cells, cell)
			grid.Add(cellContainer)
		}

		gridCells = cells
		gridBg = rects
		gridN = n
		gridArea.Objects = []fyne.CanvasObject{grid}
		gridArea.Refresh()
	}

	updateGrid = func(layout string) (*solution.TArea, error) {
		area, err := solution.InputCells(layout)
		if err != nil {
			return nil, err
		}

		currentArea = area

		nEntry.SetText(strconv.Itoa(area.N()))

		opts := layoutToRegion(layout)
		buildSelectGrid(area.N(), opts)

		rawRows := strings.Split(strings.TrimSpace(layout), "\n")
		rows := make([]string, 0, len(rawRows))
		for _, row := range rawRows {
			fixed := strings.TrimSpace(row)
			if fixed != "" {
				rows = append(rows, fixed)
			}
		}
		for row := 0; row < area.N(); row++ {
			for col := 0; col < area.N(); col++ {
				cell := gridCells[row*area.N()+col]
				region := string(rows[row][col])
				cell.SetSelected(region)
				if hexColor, ok := solution.CellColor[region]; ok {
					j := row*area.N() + col
					if j < len(gridBg) {
						gridBg[j].FillColor = HexToRgb(hexColor)
						gridBg[j].Refresh()
					}
				}
			}
		}
		return area, nil
	}

	buildVisualGrid := func(area *solution.TArea, queens []solution.TPosition) *fyne.Container {
		if area == nil {
			area = currentArea
		}
		if area == nil {
			return container.NewMax(widget.NewLabel("No board loaded"))
		}

		n := area.N()
		grid := container.NewGridWithColumns(n)

		queenMap := make(map[string]bool)
		for _, q := range queens {
			key := fmt.Sprintf("%d,%d", q.Row(), q.Col())
			queenMap[key] = true
		}

		for row := 0; row < n; row++ {
			for col := 0; col < n; col++ {
				region := area.RegionAt(row, col)
				key := fmt.Sprintf("%d,%d", row, col)
				hasQueen := queenMap[key]

				// Ambil warna region
				hexColor := area.Color(row, col)
				if hexColor == "" {
					hexColor = "#FFFFFF"
				}
				bgColor := HexToRgb(hexColor)
				
				// Buat background rectangle dengan warna region
				bgRect := canvas.NewRectangle(bgColor)
				bgRect.SetMinSize(fyne.NewSize(50, 50))

				var cellContent *fyne.Container

				if hasQueen {
					queenText := canvas.NewText("👑", &color.NRGBA{R: 0, G: 0, B: 0, A: 255})
					queenText.Alignment = fyne.TextAlignCenter
					queenText.TextSize = 32 // Buat crown lebih besar
					queenText.TextStyle = fyne.TextStyle{Bold: true}
					cellContent = container.NewStack(bgRect, container.NewCenter(queenText))
				} else {
					regionLabel := widget.NewLabel(region)
					regionLabel.Alignment = fyne.TextAlignCenter
					cellContent = container.NewStack(bgRect, regionLabel)
				}

				grid.Add(cellContent)
			}
		}

		return grid
	}

	generateButton := widget.NewButton("Enter", func() {
		textN := nEntry.Text
		if textN == "" {
			dialog.NewInformation("Error", "N kosong.", w).Show()
			return
		}

		n, err := strconv.Atoi(textN)
		if err != nil || n <= 0 {
			dialog.NewInformation("Error", "N harus positif.", w).Show()
			return
		}

		if n > 26 {
			dialog.NewInformation("Error", "Maksimal 26 region (A-Z).", w).Show()
			return
		}

		buildSelectGrid(n, showNLetters(n))
	})

	runSolver := func(
		solver func(*solution.TArea, func([]solution.TPosition) bool) []solution.TPosition,
		area *solution.TArea,
		layout string,
	) {
		// Kill proses lama
		if stopSolve != nil {
			close(stopSolve)
		}
		stopSolve = make(chan struct{})

		go func(area *solution.TArea, layout string) {
			var result []solution.TPosition
			stepCount := 0
			cancelled := false

			ans := solver(area, func(candidate []solution.TPosition) bool {
				select {
				case <-stopSolve:
					cancelled = true
					return false
				default:
				}

				stepCount++

				if stepCount%50 == 0 || stepCount < 20 {
					boardArea.Objects = []fyne.CanvasObject{buildVisualGrid(area, candidate)}
					boardArea.Refresh()
				}
				return true // Continue searching
			})

			if cancelled {
				return
			}

			if ans == nil {
				boardArea.Objects = []fyne.CanvasObject{buildVisualGrid(area, nil)}
				boardArea.Refresh()

				outputPath, err := SaveOutput(area, nil, layout)
				if err == nil {
					dialog.NewInformation("Hasil", fmt.Sprintf("Tidak ada solusi untuk layout ini.\n\nOutput disimpan ke:\n%s", outputPath), w).Show()
				} else {
					dialog.NewInformation("Hasil", fmt.Sprintf("Tidak ada solusi untuk layout ini.\n\nError menyimpan output: %v", err), w).Show()
				}
			} else {
				result = ans
				boardArea.Objects = []fyne.CanvasObject{buildVisualGrid(area, result)}
				boardArea.Refresh()
				outputPath, _ := SaveOutput(area, result, layout)
				dialog.NewInformation("Hasil", fmt.Sprintf("Solusi ditemukan setelah %d langkah!\n\nOutput disimpan ke:\n%s", stepCount, outputPath), w).Show()
			}
		}(area, layout)
	}

	solveButton := widget.NewButton("Solve", func() {
		// Baca layout dari grid
		if gridN <= 0 || len(gridCells) == 0 {
			dialog.NewInformation("Error", "Silakan isi N dan tekan Enter untuk membuat grid, atau load dari .txt.", w).Show()
			return
		}

		var input strings.Builder
		for row := 0; row < gridN; row++ {
			for col := 0; col < gridN; col++ {
				cell := gridCells[row*gridN+col]
				if cell.Selected == "" || cell.Selected == "-" {
					dialog.NewInformation("Error", "Tiap cell harus bernilai.", w).Show()
					return
				}
				input.WriteString(cell.Selected)
			}
			if row < gridN-1 {
				input.WriteString("\n")
			}
		}
		layout := input.String()

		area, err := updateGrid(layout)
		if err != nil {
			dialog.NewInformation("Error", err.Error(), w).Show()
			return
		}

		// Validasi region tidak terputus
		if err := solution.ValidRegion(area); err != nil {
			dialog.NewInformation("Error", err.Error(), w).Show()
			return
		}

		// Grid awal (no queens)
		boardArea.Objects = []fyne.CanvasObject{buildVisualGrid(area, nil)}
		boardArea.Refresh()

		runSolver(solution.TryPosition, area, layout)
	})

	optimizedButton := widget.NewButton("Optimized", func() {
		if gridN <= 0 || len(gridCells) == 0 {
			dialog.NewInformation("Error", "Silakan isi N dan tekan Enter untuk membuat grid, atau load dari .txt.", w).Show()
			return
		}

		var input strings.Builder
		for row := 0; row < gridN; row++ {
			for col := 0; col < gridN; col++ {
				cell := gridCells[row*gridN+col]
				if cell.Selected == "" || cell.Selected == "-" {
					dialog.NewInformation("Error", "Tiap cell harus bernilai.", w).Show()
					return
				}
				input.WriteString(cell.Selected)
			}
			if row < gridN-1 {
				input.WriteString("\n")
			}
		}
		layout := input.String()

		area, err := updateGrid(layout)
		if err != nil {
			dialog.NewInformation("Error", err.Error(), w).Show()
			return
		}

		if err := solution.ValidRegion(area); err != nil {
			dialog.NewInformation("Error", err.Error(), w).Show()
			return
		}

		boardArea.Objects = []fyne.CanvasObject{buildVisualGrid(area, nil)}
		boardArea.Refresh()

		runSolver(solution.TryPositionOptimized, area, layout)
	})

	stopButton := widget.NewButton("Stop", func() {
		if stopSolve != nil {
			close(stopSolve)
			stopSolve = nil
		}
	})

	topBar := container.NewVBox(
		nEntry,
		container.NewBorder(nil, nil, nil, inputTextButton, generateButton),
		container.NewGridWithColumns(2, solveButton, stopButton),
		optimizedButton,
	)

	leftSide := container.NewBorder(nil, nil, nil, nil, gridArea)

	boardPlaceholder := widget.NewLabel("Proses pencarian solusi akan divisualisasikan di sini.")
		boardArea = container.NewMax(boardPlaceholder)

	rightSide := container.NewVBox(
		widget.NewLabel("Grid Process"),
		boardArea,
	)

	leftScroll := container.NewScroll(leftSide)
	rightScroll := container.NewScroll(rightSide)

	split := container.NewHSplit(leftScroll, rightScroll)
	split.SetOffset(0.5) // 50% kiri, 50% kanan

	container := container.NewBorder(topBar, nil, nil, nil, split)

	w.SetContent(container)
	w.ShowAndRun()
}

func layoutToRegion(layout string) []string {
	layout = strings.TrimSpace(layout)
	rawRows := strings.Split(layout, "\n")
	seen := make(map[byte]bool)
	opts := make([]string, 0, 26)
	for _, row := range rawRows {
		fixed := strings.TrimSpace(row)
		if fixed == "" {
			continue
		}
		for i := 0; i < len(fixed); i++ {
			c := fixed[i]
			if !seen[c] {
				seen[c] = true
				opts = append(opts, string(c))
			}
		}
	}
	return opts
}

func NewFileOpen(callback func(fyne.URIReadCloser, error), parent fyne.Window) *dialog.FileDialog {
	dialog := dialog.NewFileOpen(callback, parent)
	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	return dialog
}
