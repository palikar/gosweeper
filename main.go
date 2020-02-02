package main

import (
	"fmt"
	"fyne.io/fyne"

	"fyne.io/fyne/app"
	// "fyne.io/fyne/canvas"
	// "image/color"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func clickMine(x int, y int) {
	fmt.Printf("%s:%d:%d\n", "click", x, y)

}

func gameScreen(a fyne.App) fyne.CanvasObject {

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.RadioButtonCheckedIcon(), func() {}),
		widget.NewToolbarSpacer(),
	)

	grid := layout.NewGridLayout(20)

	cont := fyne.NewContainerWithLayout(grid)

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			b := widget.NewButton("", func(i int, j int) func() { return func() { clickMine(i, j) } }(i, j))
			b.Resize(fyne.NewSize(20, 20))
			cont.AddObject(b)
		}
	}

	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar, cont)

	return content
}

func main() {
	fmt.Println("Starting The Game")
	a := app.NewWithID("palikar.go.sweeper")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("GoSweeper")

	tabs := widget.NewTabContainer(widget.NewTabItemWithIcon("Game", theme.HomeIcon(), gameScreen(a)))
	tabs.SetTabLocation(widget.TabLocationLeading)

	w.Resize(fyne.NewSize(750, 750))

	w.SetContent(tabs)
	w.ShowAndRun()
}
