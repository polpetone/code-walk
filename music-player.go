package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

func play(musicFile string, soundChannel chan bool) error {

	Info.Println("Play music: ", musicFile)

	f, err := os.Open(musicFile)

	if err != nil {
		return err
	}

	s, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(beep.Seq(s))

	twiceAsFast := beep.ResampleRatio(4,2 ,s)
	for {
		select {
		case <-soundChannel:
			Info.Println("Music Player play twice as fast")
			speaker.Play(twiceAsFast)
		default:
		}
	}

	return nil
}
