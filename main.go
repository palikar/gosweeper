package main

import (
	"fmt"
	"fyne.io/fyne"

	"fyne.io/fyne/app"
	"math/rand"
	 "strconv"
	// "fyne.io/fyne/canvas"
	// "image/color"
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

var buttonGrid [20][20]MineCell

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

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			buttonGrid[i][j] = NewCell(nil)
		}
	}

	for i := 0; i < 50; i++ {
		rX := rand.Int31n(20)
		rY := rand.Int31n(20)
		buttonGrid[rX][rY].hasMine = true
		}

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {

			if buttonGrid[i][j].hasMine {
				continue
			}
			
			for dX := -1; dX <= 1; dX++ {
				
				for dY := -1; dY <= 1; dY++ {

					x := i + dX
					y := j + dY

					if x >= 20 {
						continue
					}
					if y >= 20 {
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

func propagate(x int, y int) {
	if x < 0 || y < 0  || x >= 20 || y >= 20 {
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
	fmt.Printf("%s:%d:%d\n", "click", x, y)

	if buttonGrid[x][y].hasMine {
		fmt.Printf("%s\n", "Boom!")
		return
	}

	propagate(x,y)



}

func clickFlag(x int, y int) {
	fmt.Printf("%s:%d:%d\n", "flag", x, y)

	buttonGrid[x][y].hasFlag = true
	buttonGrid[x][y].btn.SetText("P")
	buttonGrid[x][y].btn.Disable()
}

// restart ...
func restart() {

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			buttonGrid[i][j].btn.SetText("")
			buttonGrid[i][j].btn.Enable()
		}
	}

}

func gameScreen(a fyne.App) fyne.CanvasObject {

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.RadioButtonCheckedIcon(), restart),
		widget.NewToolbarSpacer(),
	)

	grid := layout.NewGridLayout(20)

	cont := fyne.NewContainerWithLayout(grid)

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			b := widget.NewButton("", func(i int, j int) func() { return func() { clickMine(i, j) } }(i, j))
			b.OnSecondaryTapped = func(i int, j int) func() { return func() { clickFlag(i, j) } }(i, j)
			b.Resize(fyne.NewSize(20, 20))

			// b.SetText(strconv.Itoa(buttonGrid[i][j].neighbourMineCount))
			
			buttonGrid[i][j].btn = b
			cont.AddObject(b)
		}
	}

	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar, cont)

	return content
}

func main() {
	fmt.Println("Starting The Game")
	initGrid()

	a := app.NewWithID("palikar.go.sweeper")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("GoSweeper")

	tabs := widget.NewTabContainer(widget.NewTabItemWithIcon("Game", theme.HomeIcon(), gameScreen(a)))
	tabs.SetTabLocation(widget.TabLocationLeading)

	w.Resize(fyne.NewSize(750, 750))

	w.SetContent(tabs)
	w.ShowAndRun()
}
