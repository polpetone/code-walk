package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

func keyHandler(
	delayChannel chan time.Duration,
	snapShotChannel chan bool,
	haltChannel chan bool,
	continueChannel chan bool,
) {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	var delay time.Duration = 200
	var delayStep time.Duration = DELAY_STEP

	defer Info.Println("KeyHandler stopped")

	var hackerTyperMode = false

mainloop:
	for {
		Trace.Println("HackerTyperMode:", hackerTyperMode)

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			}

			if ev.Ch == 'm' {
				if hackerTyperMode {
					hackerTyperMode = false
				} else {
					hackerTyperMode = true
					haltChannel <- true
				}
			}

			if hackerTyperMode {
				if ev.Ch != 'm' {
					continueChannel <- true
					haltChannel <- true
				}
			} else {
				delay, delayStep = commandMode(ev, delay, delayStep, delayChannel, snapShotChannel, haltChannel, continueChannel)
			}

		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventInterrupt:
			break mainloop
		}

	}
}

func commandMode(ev termbox.Event, delay time.Duration, delayStep time.Duration,
	delayChannel chan time.Duration,
	snapShotChannel chan bool,
	haltChannel chan bool,
	continueChannel chan bool,
) (time.Duration, time.Duration) {
	Info.Println("Pressed: ", ev.Ch)
	if ev.Ch == '+' {
		delay, delayStep = decreaseDelay(delay, delayStep, delayChannel)
	}
	if ev.Ch == '-' {
		delay, delayStep = increaseDelay(delay, delayStep, delayChannel)
	}
	if ev.Ch == 'c' {
		colorChannel <- true
	}
	if ev.Ch == 's' {
		snapShotChannel <- true
	}
	if ev.Ch == 'h' {
		haltChannel <- true
	}
	if ev.Ch == 'g' {
		continueChannel <- true
	}
	return delay, delayStep
}
