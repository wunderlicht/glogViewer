package main

import (
	"bufio"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
	"os"
)

type glogger struct {
	autoScroll      bool
	window          fyne.Window
	scrollContainer widget.ScrollContainer
	container       *fyne.Container
	canvas          fyne.CanvasObject
}

func (g *glogger) handleTypedKey(ke *fyne.KeyEvent) {
	g.scrollContainer.Scrolled(&fyne.ScrollEvent{
		PointEvent: fyne.PointEvent{},
		DeltaX:     0,
		DeltaY:     0,
	})
}

func (g *glogger) handleTypedRune(r rune) {
	switch r {
	case 't':
		g.scrollContainer.ScrollToTop()
		g.autoScroll = false
	case 'b':
		g.scrollContainer.ScrollToBottom()
		g.autoScroll = true
	}
}

func (g *glogger) setup(app fyne.App) {
	g.autoScroll = true
	g.window = app.NewWindow("glogViewer")
	g.container = fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	g.scrollContainer = *widget.NewScrollContainer(g.container)
	g.window.SetContent(&g.scrollContainer)
	g.window.Canvas().SetOnTypedKey(g.handleTypedKey)
	g.window.Canvas().SetOnTypedRune(g.handleTypedRune)
}

func main() {
	myApp := app.New()
	glog := new(glogger)

	glog.setup(myApp)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			glog.container.AddObject(canvas.NewText(line, color.White))
			glog.container.Refresh()
			if glog.autoScroll {
				glog.scrollContainer.ScrollToBottom()
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	glog.window.ShowAndRun()
}
