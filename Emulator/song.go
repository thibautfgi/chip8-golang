package emulator

import (
	"os"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

func Song(songfile string) error {
	f, err := os.Open("./SoundBank/" + songfile + ".mp3") //ouvre et lis le file sound
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f) //decode en mp3 le file
	if err != nil {
		return err
	}

	c, ready, err := oto.NewContext(d.SampleRate(), 2, 2) // cree un nouvo audio context
	if err != nil {
		return err
	}
	<-ready

	p := c.NewPlayer(d) // crée un player
	defer p.Close()
	p.Play() //joue

	for {
		// time.Sleep(time.Second)
		if !p.IsPlaying() {
			break
		}
	}

	return nil
}
