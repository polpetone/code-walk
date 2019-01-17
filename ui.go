package main

import (
	"bytes"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"time"
)

func ui() {
	ui, codeBoxWriter, authorsBoxWriter, fileInfoBoxWriter := setupUI()
	go runUI(ui)
	for {
		select {
		case code := <-codeChannel:
			fmt.Fprint(codeBoxWriter, code)
		case text := <-codeWalkFileInfoChannel:
			authorsBoxWriter.buf.Reset()
			authorsBoxWriter.Label.SetText("")
			fileInfoBoxWriter.buf.Reset()
			fileInfoBoxWriter.Label.SetText("")
			fmt.Fprintln(fileInfoBoxWriter, text.FileName)
			for _, a := range text.Authors {
				fmt.Fprintln(authorsBoxWriter, a)
			}
		default:
		}
	}
}
type labelWriter struct {
	buf bytes.Buffer
	ui tui.UI
	tui.Label
}

func (w *labelWriter) OnKeyEvent(ev tui.KeyEvent) {
	if ev.Key == tui.KeyCtrlL {
		w.buf.Reset()
		w.Label.SetText("")
		go w.ui.Update(func() {})
	}
}

func (w *labelWriter) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	w.Label.SetText(w.buf.String())
	go w.ui.Update(func() {})
	return
}

var delay time.Duration = 200
var delayStep time.Duration = DELAY_STEP

func setupUI() (tui.UI, *labelWriter, *labelWriter, *labelWriter) {

	header := tui.NewHBox()
	header.SetBorder(true)
	header.SetSizePolicy(tui.Expanding, tui.Maximum)

	header.Append(tui.NewHBox(
		tui.NewLabel("CODE WALK"),
		))

	codeBox := tui.NewVBox()
	codeBox.SetSizePolicy(tui.Expanding, tui.Expanding)

	codeBoxScroll := tui.NewScrollArea(codeBox)
	codeBoxScroll.SetAutoscrollToBottom(true)

	box := tui.NewVBox(codeBoxScroll)
	box.SetBorder(true)
	box.SetSizePolicy(tui.Expanding, tui.Maximum)

	authorsBox := tui.NewVBox()
	authorsBox.SetBorder(true)
	authorsBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	fileInfoBox := tui.NewVBox()
	fileInfoBox.SetBorder(true)
	fileInfoBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	fileInfoAuthorsContainer := tui.NewVBox(fileInfoBox, authorsBox)

	container := tui.NewHBox(box, fileInfoAuthorsContainer)
	container.SetSizePolicy(tui.Expanding, tui.Maximum)

	root := tui.NewVBox(header, container)
	ui, err := tui.New(root)
	if err != nil {
		Error.Println(err)
	}
	ui.SetKeybinding("ESC", func()  {ui.Quit()})

	ui.SetKeybinding("c", func() {
		colorChannel <- true
		Trace.Println("UI Keybinding Pressed c")
	})

	ui.SetKeybinding("+", func() {
		delay, delayStep = decreaseDelay(delay, delayStep, delayChannel)
		Trace.Println("UI Keybinding Pressed +")
	})

	ui.SetKeybinding("-", func() {
		delay, delayStep = increaseDelay(delay, delayStep, delayChannel)
		Trace.Println("UI Keybinding Pressed -")
	})

	ui.SetKeybinding("h", func() {
		haltChannel <-true
		Trace.Println("UI Keybinding Pressed h")
	})

	ui.SetKeybinding("g", func() {
		continueChannel <-true
		Trace.Println("UI Keybinding Pressed g")
	})

	codeBoxWriter := &labelWriter{ui : ui}
	fileInfoBoxWriter := &labelWriter{ui : ui}
	authorsBoxWriter := &labelWriter{ui : ui}
	codeBox.Append(codeBoxWriter)
	authorsBox.Append(authorsBoxWriter)
	fileInfoBox.Append(fileInfoBoxWriter)
	return ui, codeBoxWriter, authorsBoxWriter, fileInfoBoxWriter
}

func runUI(ui tui.UI){
	if err := ui.Run(); err != nil {
		Error.Println(err)
	}
}





