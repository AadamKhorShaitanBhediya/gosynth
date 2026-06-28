package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BASE_FREQ   = 220
	SAMPLE_RATE = 44100
	BUFFER_SIZE = 4096
	TEXT_SIZE   = 20
)

var (
	sampleTime float32
)

func NoteToFreq(n uint) float32 {
	r := math.Pow(2, float64(1)/12)
	return BASE_FREQ * float32(math.Pow(r, float64(n)))
}

func PollNote(note *uint) {

	if rl.IsKeyDown(rl.KeyOne) {
		*note = 0
	}
	if rl.IsKeyDown(rl.KeyTwo) {
		*note = 1
	}
	if rl.IsKeyDown(rl.KeyThree) {
		*note = 2
	}
	if rl.IsKeyDown(rl.KeyFour) {
		*note = 3
	}
	if rl.IsKeyDown(rl.KeyFive) {
		*note = 4
	}
	if rl.IsKeyDown(rl.KeySix) {
		*note = 5
	}
	if rl.IsKeyDown(rl.KeySeven) {
		*note = 6
	}
	if rl.IsKeyDown(rl.KeyEight) {
		*note = 7
	}
	if rl.IsKeyDown(rl.KeyNine) {
		*note = 8
	}
	if rl.IsKeyDown(rl.KeyZero) {
		*note = 9
	}
	if rl.IsKeyDown(rl.KeyMinus) {
		*note = 10
	}
	if rl.IsKeyDown(rl.KeyEqual) {
		*note = 11
	}
}

func main() {
	var (
		note uint
		A    float32 = 0.3
		f    float32 = BASE_FREQ
	)

	rl.InitWindow(800, 450, "Sound Synthesis")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	//number of samples to keep in memory at a time
	rl.SetAudioStreamBufferSizeDefault(BUFFER_SIZE)
	buffer := [BUFFER_SIZE]float32{}

	//sampleSize = 32 bit floats
	stream := rl.LoadAudioStream(SAMPLE_RATE, 32, 1)
	rl.SetAudioStreamPan(stream, 0.0)
	rl.PlayAudioStream(stream)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		PollNote(&note)

		if rl.IsAudioStreamProcessed(stream) {
			for i := range buffer {
				y := A * float32(math.Sin(2*math.Pi*float64(f*sampleTime)))
				buffer[i] = y
				sampleTime += 1.0 / SAMPLE_RATE
			}
			rl.UpdateAudioStream(stream, buffer[:]) //colon for slice
		}
		f = NoteToFreq(note)
		rl.DrawText(fmt.Sprint("frequency: ", f), 0, TEXT_SIZE*0, TEXT_SIZE, rl.Red)
		rl.DrawText(fmt.Sprint("Note: ", note), 0, TEXT_SIZE*1, TEXT_SIZE, rl.Red)
		rl.EndDrawing()
	}
}
