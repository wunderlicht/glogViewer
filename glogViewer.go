package main

import (
	"bufio"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"

	//"go/scanner"
	"os"
	//"strings"
)

func main() {
	myApp := app.New()

	w := myApp.NewWindow("glogViewer")
	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	scrollContainer := widget.NewScrollContainer(container)
	w.SetContent(scrollContainer)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			//appendLine(tg, line)
			container.AddObject(canvas.NewText(line, color.White))
			container.Refresh()
			scrollContainer.ScrollToBottom()
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	w.ShowAndRun()
}
