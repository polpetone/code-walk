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
			time.Sleep(5 * time.Millisecond)
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
var delayStep = DELAY_STEP

func setupUI() (tui.UI, *labelWriter, *labelWriter, *labelWriter) {
	codeBox := tui.NewVBox()

	codeBoxScroll := tui.NewScrollArea(codeBox)
	codeBoxScroll.SetAutoscrollToBottom(true)

	box := tui.NewHBox(codeBoxScroll)
	box.SetBorder(false)

	authorsBox := tui.NewVBox()
	authorsBox.SetBorder(false)

	fileInfoBox := tui.NewVBox()
	fileInfoBox.SetBorder(false)
	fileInfoBox.SetSizePolicy(tui.Maximum, tui.Maximum)

	fileInfoAuthorsContainer := tui.NewVBox(authorsBox)

	container := tui.NewHBox(box, fileInfoAuthorsContainer)

	root := tui.NewVBox(fileInfoBox, container)
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

	ui.SetKeybinding("j", func() {
		jumpFileChannel <-true
		Trace.Println("UI Keybinding Pressed j")
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





