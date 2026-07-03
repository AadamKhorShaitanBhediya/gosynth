package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BASE_FREQ   = 220
	SAMPLE_RATE = 44100
	BUFFER_SIZE = 1024
	TEXT_SIZE   = 20
)

const (
	WaveSine = iota
	WaveSquare
	WaveTriangle
	WaveSaw
)

var (
	noteNames = [...]string{
		"None",
		"A",
		"A#",
		"B",
		"C",
		"C#",
		"D",
		"D#",
		"E",
		"F",
		"F#",
		"G",
		"G#",
		"A^",
	}

	//13 sequential keys on the QWERTY row
	qwertyKeys = []int32{
		rl.KeyQ, rl.KeyW, rl.KeyE, rl.KeyR, rl.KeyT, rl.KeyY,
		rl.KeyU, rl.KeyI, rl.KeyO, rl.KeyP, rl.KeyLeftBracket,
		rl.KeyRightBracket, rl.KeyBackSlash,
	}
	waveNames = []string{
		"Sine",
		"Square",
		"Triangle",
		"Saw",
	}
	sampleTime float32
)

type Envelope struct {
	startAmp    float32
	attackTime  float32
	decayTime   float32
	sustainAmp  float32
	releaseTime float32

	triggerOnTime  float32
	triggerOffTime float32
	keyPressed     bool
}

func (e *Envelope) GetAmp(time float32) float32 {
	var amp float32 = 0.0

	if e.keyPressed {
		lifeTime := time - e.triggerOnTime
		// attack
		if lifeTime > 0 && lifeTime <= e.attackTime {
			amp = (e.startAmp / e.attackTime) * lifeTime
		}
		// decay
		if lifeTime > e.attackTime && lifeTime <= (e.decayTime+e.attackTime) {
			amp = ((e.sustainAmp-e.startAmp)/e.decayTime)*(lifeTime-e.attackTime) + e.startAmp
		}
		// sustain
		if lifeTime > (e.decayTime + e.attackTime) {
			amp = e.sustainAmp
		}
	} else {
		// release
		lifeTime := time - e.triggerOffTime
		amp = ((0-e.sustainAmp)/e.releaseTime)*lifeTime + e.sustainAmp
	}
	if amp <= 0.0001 {
		amp = 0
	}

	return amp
}

func (e *Envelope) Trigger(timeOn float32) {
	e.triggerOnTime = timeOn
	e.keyPressed = true
}

func (e *Envelope) Untrigger(timeOff float32) {
	e.triggerOffTime = timeOff
	e.keyPressed = false
}

func DefaultEnvelope() Envelope {
	return Envelope{

		startAmp:       1.0,
		attackTime:     0.3,
		decayTime:      0.3,
		sustainAmp:     0.7,
		releaseTime:    0.5,
		triggerOnTime:  0.0,
		triggerOffTime: 0.0,
		keyPressed:     false,
	}
}

func NoteToFreq(n int32) float32 {
	return BASE_FREQ * float32(math.Pow(2, float64(n)/12))
}

func MakeWaves(amp float32, freq float32, waveType uint, sampleTime *float32, sampleSize uint32, sampleRate uint32) []float32 {
	buffer := make([]float32, sampleSize)
	for i := range buffer {
		*sampleTime += 1.0 / float32(sampleRate)
		y := float32(math.Sin(2 * math.Pi * float64(freq*(*sampleTime))))
		switch waveType {
		case WaveSine:
			y = y + 0
		case WaveSquare:
			if y > 0 {
				y = 1
			} else {
				y = -1
			}
		case WaveTriangle:
			y = float32(math.Asin(float64(y))) / (math.Pi / 2)
		case WaveSaw:
			// pi*freq = omega
			y = float32(math.Atan(math.Tan(math.Pi*float64(freq*(*sampleTime))))) / (math.Pi / 2)
		}
		buffer[i] = y * amp
	}
	return buffer
}

func main() {
	var (
		note    int32
		A       float32 = 0.3
		f       float32 = BASE_FREQ
		lastkey int32   = -1 //invalid note key
	)
	env := DefaultEnvelope()

	rl.InitWindow(800, 450, "Music Synthesis")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	// number of samples to keep in memory at a time
	rl.SetAudioStreamBufferSizeDefault(BUFFER_SIZE)

	// sampleSize -> 32 bit floats
	stream := rl.LoadAudioStream(SAMPLE_RATE, 32, 1)
	rl.SetAudioStreamPan(stream, 0.0)
	rl.PlayAudioStream(stream)

	waveType := WaveSine

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for notenum, key := range qwertyKeys {
			if rl.IsKeyPressed(key) {
				lastkey = key
				env.Trigger(float32(rl.GetTime()))
				f = NoteToFreq(int32(notenum))
				note = int32(notenum + 1)
				break
			}
		}

		if lastkey != -1 && rl.IsKeyReleased(lastkey) {
			lastkey = -1
			note = 0
			env.Untrigger(float32(rl.GetTime()))
		}

		A = env.GetAmp(float32(rl.GetTime()))

		if rl.IsKeyPressed(rl.KeySpace) {
			waveType = (waveType + 1) % 4
		}

		visualBuffer := []float32{}
		if rl.IsAudioStreamProcessed(stream) {
			buffer := MakeWaves(A, f, uint(waveType), &sampleTime, BUFFER_SIZE, SAMPLE_RATE)
			rl.UpdateAudioStream(stream, buffer[:]) // colon for slice
			visualBuffer = buffer[:800]
		}
		for i := 0; i < len(visualBuffer); i += 2 {
			rl.DrawCircle(int32(i), int32(100*visualBuffer[i]+float32(rl.GetScreenHeight()/2)), 3, rl.Blue)
		}

		rl.DrawText(fmt.Sprintf("Frequency: %f Hz", f), 0, TEXT_SIZE*0, TEXT_SIZE, rl.Red)
		rl.DrawText(fmt.Sprintf("Note: %v", noteNames[note]), 0, TEXT_SIZE*1, TEXT_SIZE, rl.Red)
		rl.DrawText(fmt.Sprintf("Amplitude: %v", A), 0, TEXT_SIZE*2, TEXT_SIZE, rl.Red)
		rl.DrawText(fmt.Sprintf("Wave Type (Press Spacebar): %v", waveNames[waveType]), 0, TEXT_SIZE*3, TEXT_SIZE, rl.Red)
		rl.EndDrawing()
	}
}
