package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

func start(musicFiles []string, soundChannel chan bool){
	for _, musicFile := range musicFiles {
		Info.Println("Play: ", musicFile)
		play(musicFile, soundChannel)
		Info.Println("Play DONE: ", musicFile)
	}
}

func play(musicFile string, soundChannel chan bool) error {
	ratio := 1.0

	f, err := os.Open(musicFile)

	if err != nil {
		return err
	}

	s, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	Info.Println("Format SampleRate:", format.SampleRate)

	speaker.Play(beep.Seq(s))

	for {
		select {
		case soundSignal := <-soundChannel:
			if soundSignal {
				ratio = ratio * 1.0
				Info.Println("Current Ratio: ", ratio)
				playRatio := beep.ResampleRatio(4, ratio, s)
				Info.Println("Speaker ", s.Position())
				Info.Println("Speaker ", s.Len())
				speaker.Play(playRatio)
			} else {
				ratio = ratio / 1.0
				Info.Println("Current Ratio: ", ratio)
				playRatio := beep.ResampleRatio(4, ratio, s)
				Info.Println("Speaker ", s.Position())
				Info.Println("Speaker ", s.Len())
				speaker.Play(playRatio)
			}
		default:
		}
		time.Sleep(300 * time.Millisecond)
		if s.Position() == s.Len(){
			break
		}
	}

	return nil
}
