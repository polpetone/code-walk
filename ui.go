package main

import (
	"bytes"
	"fmt"
	"github.com/marcusolsson/tui-go"
)

var codeWalkFileInfoChannel = make(chan CodeWalkFileInfo)

func ui() {
	ui, w , w2:= setupUI()
	go runUI(ui)
	for {
		select {
		case code := <-codeChannel:
			fmt.Fprint(w, code)
		case text := <-codeWalkFileInfoChannel:
			w2.buf.Reset()
			w2.Label.SetText("")
			for _, a := range text.Authors {
				fmt.Fprintln(w2, a)
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

func setupUI() (tui.UI, *labelWriter, *labelWriter) {

	header := tui.NewHBox()
	header.SetBorder(true)
	header.SetSizePolicy(tui.Expanding, tui.Maximum)

	header.Append(tui.NewHBox(
		tui.NewLabel("CODE WALK"),
		))


	fileInfoBox := tui.NewVBox()
	fileInfoBox.SetBorder(true)
	fileInfoBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	authorsBox := tui.NewVBox()
	authorsBox.SetBorder(true)
	authorsBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	container := tui.NewHBox(fileInfoBox, authorsBox)
	container.SetSizePolicy(tui.Expanding, tui.Maximum)

	root := tui.NewVBox(header, container)
	ui, err := tui.New(root)
	if err != nil {
		Error.Println(err)
	}
	ui.SetKeybinding("ESC", func()  {ui.Quit()})

	fileInfoBoxWriter := &labelWriter{ui : ui}
	authorsBoxWriter := &labelWriter{ui : ui}
	fileInfoBox.Append(fileInfoBoxWriter)
	authorsBox.Append(authorsBoxWriter)
	return ui, fileInfoBoxWriter, authorsBoxWriter
}

func runUI(ui tui.UI){
	if err := ui.Run(); err != nil {
		Error.Println(err)
	}
}
