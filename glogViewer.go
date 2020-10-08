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
	"strings"
)

type glogger struct {
	autoScroll      bool
	window          fyne.Window
	scrollContainer *widget.ScrollContainer
	container       *fyne.Container
}

var (
	red    = color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff}
	yellow = color.RGBA{R: 0xff, G: 0xff, B: 0, A: 0xff}
)

func (g *glogger) handleTypedKey(ke *fyne.KeyEvent) {
	delta := g.container.Objects[0].Size().Height
	switch ke.Name {
	case fyne.KeyUp:
		g.scrollContainer.Scrolled(&fyne.ScrollEvent{DeltaY: delta})
	case fyne.KeyDown:
		g.scrollContainer.Scrolled(&fyne.ScrollEvent{DeltaY: -delta})
	default:
		return //no key handled, leave autoScroll untouched
	}
	g.autoScroll = false
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
	g.window.Resize(fyne.Size{Width: 1000, Height: 600})
	g.container = fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	g.scrollContainer = widget.NewScrollContainer(g.container)
	g.window.SetContent(g.scrollContainer)
	g.window.Canvas().SetOnTypedKey(g.handleTypedKey)
	g.window.Canvas().SetOnTypedRune(g.handleTypedRune)
}

func (g *glogger) addLine(line string, severity string) {
	switch severity {
	case "ERROR":
		g.container.AddObject(canvas.NewText(line, yellow))
	case "FATAL":
		g.container.AddObject(canvas.NewText(line, red))
	default:
		g.container.AddObject(canvas.NewText(line, color.White))
	}
	g.container.Refresh()
}

func (g *glogger) splitLine(line string) (source string, severity string, msg string) {
	parts := strings.Split(line, "|")
	source = parts[0]
	severity = parts[1]
	msg = parts[2]
	return
}

func main() {
	myApp := app.New()
	glog := new(glogger)

	glog.setup(myApp)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			source, severity, msg := glog.splitLine(line)
			glog.addLine(line, severity)
			if severity == "FATAL" {
				myApp.SendNotification(&fyne.Notification{
					Title:   source + " " + severity,
					Content: msg,
				})
			}

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
