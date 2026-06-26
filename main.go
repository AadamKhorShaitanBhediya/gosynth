package main

import (
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

// samples [value][channel, 2 for headphones]
const sr float64 = 44100

var t float64

func MakeNoise(samples [][2]float64) (n int, ok bool) {
	var A float64 = 0.5
	var f float64 = 220
	var w = f * 2 * float64(math.Pi)
	dt := 1 / sr
	for i := range samples {
		value := A * math.Sin(w*t)
		samples[i][0] = value
		samples[i][1] = value
		t += dt
	}
	return len(samples), true
}
func main() {

	sr := beep.SampleRate(sr)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(beep.StreamerFunc(MakeNoise))
	select {}

}
