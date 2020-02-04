package main

import (
	"fmt"
	"fyne.io/fyne"
	"os"
	"time"

	"fyne.io/fyne/app"
	"math/rand"
	"strconv"
	// "fyne.io/fyne/canvas"
	// "image/color"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// MineCell ...
type MineCell struct {
	hasMine            bool
	hasFlag            bool
	opened             bool
	neighbourMineCount int
	btn                *widget.Button
}

var buttonGrid [][]MineCell

var width = 16
var height = 16
var minesInit = false
var minesCnt = 40

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// NewCell ...
func NewCell(btn *widget.Button) MineCell {
	cell := MineCell{}
	cell.hasMine = false
	cell.hasFlag = false
	cell.opened = false
	cell.neighbourMineCount = 0
	return cell
}

// initGrid ...
func initGrid() {

	minesInit = false

	buttonGrid = make([][]MineCell, width)
	for i := 0; i < width; i++ {
		buttonGrid[i] = make([]MineCell, height)
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {

			if buttonGrid[i][j].btn != nil {
				p := buttonGrid[i][j].btn
				p.Enable()
				p.SetText("")
				buttonGrid[i][j] = NewCell(nil)
				buttonGrid[i][j].btn = p
				continue
			}

			buttonGrid[i][j] = NewCell(nil)
		}
	}

}

// initMines ...
func initMines(xClick int, yClick int) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < minesCnt; i++ {
		rX := r1.Int31n(int32(width))
		rY := r1.Int31n(int32(height))

		for buttonGrid[rX][rY].hasMine ||
			(abs(xClick-int(rX)) < 2 && abs(yClick-int(rY)) < 2) {
			rX = r1.Int31n(int32(width))
			rY = r1.Int31n(int32(height))
		}

		buttonGrid[rX][rY].hasMine = true

	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {

			if buttonGrid[i][j].hasMine {
				continue
			}

			for dX := -1; dX <= 1; dX++ {

				for dY := -1; dY <= 1; dY++ {

					x := i + dX
					y := j + dY

					if x >= width {
						continue
					}
					if y >= height {
						continue
					}
					if x < 0 {
						continue
					}
					if y < 0 {
						continue
					}

					if buttonGrid[x][y].hasMine {
						buttonGrid[i][j].neighbourMineCount++
					}
				}
			}

		}

	}

}

// checkWin ...
func checkWin() {

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			ch := &buttonGrid[i][j]
			if !ch.hasMine && !ch.opened {
				return
			}
		}
	}

	cnf := dialog.NewConfirm("You won", "You must be very clever, bubence! One more game", func(suc bool) {
		if !suc {
			a.Quit()
			os.Exit(0)
		}
		initGrid()
		resetUI()

	}, w)
	cnf.SetConfirmText("Oh Yes!")
	cnf.SetDismissText("Nah")
	cnf.Show()

}

func propagate(x int, y int) {
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}

	cell := &buttonGrid[x][y]

	if cell.opened || cell.hasMine || cell.hasFlag {
		return
	}

	if cell.neighbourMineCount > 0 {
		cell.btn.SetText(strconv.Itoa(cell.neighbourMineCount))
		cell.opened = true
		cell.btn.Disable()
		return
	}

	if cell.neighbourMineCount == 0 {
		cell.opened = true
		cell.btn.Disable()

		propagate(x, y+1)

		propagate(x, y-1)

		propagate(x+1, y)

		propagate(x-1, y)

		propagate(x+1, y+1)

		propagate(x+1, y-1)

		propagate(x-1, y+1)

		propagate(x-1, y-1)

	}

}

func clickMine(x int, y int) {

	if !minesInit {
		initMines(x, y)
		minesInit = true
	}

	if buttonGrid[x][y].hasMine {

		cnf := dialog.NewConfirm("Game Over", "You clicked on mine. Do you want a new game?", func(suc bool) {
			if !suc {
				a.Quit()
				os.Exit(0)
			}
			initGrid()
			resetUI()

		}, w)
		cnf.SetConfirmText("Oh Yes!")
		cnf.SetDismissText("Nah")
		cnf.Show()

		return
	}

	propagate(x, y)

	checkWin()

}

func clickFlag(x int, y int) {

	if buttonGrid[x][y].hasFlag {
		buttonGrid[x][y].hasFlag = false
		buttonGrid[x][y].btn.SetText("")
		return
	}

	buttonGrid[x][y].hasFlag = true
	buttonGrid[x][y].btn.SetText("P")
}

// restart ...
func restart() {
	initGrid()
	resetUI()
}

func gameScreen(a fyne.App) fyne.CanvasObject {

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.RadioButtonCheckedIcon(), restart),
		widget.NewToolbarSpacer(),
	)

	grid := layout.NewGridLayout(width)

	cont := fyne.NewContainerWithLayout(grid)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			b := widget.NewButton("", func(i int, j int) func() { return func() { clickMine(i, j) } }(i, j))
			b.OnSecondaryTapped = func(i int, j int) func() { return func() { clickFlag(i, j) } }(i, j)
			b.Resize(fyne.NewSize(30, 30))

			// b.SetText(strconv.Itoa(buttonGrid[i][j].neighbourMineCount))

			buttonGrid[i][j].btn = b
			cont.AddObject(b)
		}
	}

	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar, cont)

	return content
}

var a = app.NewWithID("palikar.go.sweeper")
var w = a.NewWindow("GoSweeper")

// resetUI ...
func resetUI() {

	tabs := widget.NewTabContainer(widget.NewTabItemWithIcon("Game", theme.HomeIcon(), gameScreen(a)))
	tabs.SetTabLocation(widget.TabLocationLeading)
	w.SetContent(tabs)

}

func main() {

	fmt.Println("Starting The Game")
	initGrid()

	a.SetIcon(theme.FyneLogo())

	w.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("Game",

			fyne.NewMenuItem("New Game (9x9)",
				func() {
					width = 9
					height = 9
					minesCnt = 10
					initGrid()
					w.Resize(fyne.NewSize(300, 300))
					resetUI()
				}),

			fyne.NewMenuItem("New Game (16x16)",
				func() {
					width = 16
					height = 16
					minesCnt = 40
					initGrid()
					w.Resize(fyne.NewSize(600, 600))
					resetUI()
				}),

			fyne.NewMenuItem("New Game (30x30)",
				func() {
					width = 30
					height = 30
					minesCnt = 120
					initGrid()
					w.Resize(fyne.NewSize(1000, 1000))
					resetUI()
				})),

		fyne.NewMenu("Info",
			fyne.NewMenuItem("About",
				func() {

				}),
			fyne.NewMenuItem("License",
				func() {

				}))))

	resetUI()

	w.Resize(fyne.NewSize(1000, 1000))

	w.ShowAndRun()
}
