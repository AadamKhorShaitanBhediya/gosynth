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
	noteNames = [...]string{
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

	// 13 sequential keys on the QWERTY row
	qwertyKeys = []int32{
		rl.KeyQ, rl.KeyW, rl.KeyE, rl.KeyR, rl.KeyT, rl.KeyY,
		rl.KeyU, rl.KeyI, rl.KeyO, rl.KeyP, rl.KeyLeftBracket,
		rl.KeyRightBracket, rl.KeyBackSlash,
	}
	sampleTime float32
)

func NoteToFreq(n uint) float32 {
	return BASE_FREQ * float32(math.Pow(2, float64(n)/12))
}

func PollNote(note *uint) bool {
	for i, key := range qwertyKeys {
		if rl.IsKeyDown(key) {
			*note = uint(i)
			return true
		}
	}
	return false
}

func main() {
	var (
		note uint
		A    float32 = 0.3
		f    float32 = BASE_FREQ
	)

	rl.InitWindow(800, 450, "Music Synthesis")
	defer rl.CloseWindow()
	rl.SetTargetFPS(24)

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

		if PollNote(&note) {
			f = NoteToFreq(note)
			rl.ResumeAudioStream(stream)
		} else {
			rl.PauseAudioStream(stream)
			sampleTime = 0.0
		}

		if rl.IsAudioStreamProcessed(stream) {
			for i := range buffer {
				y := A * float32(math.Sin(2*math.Pi*float64(f*sampleTime)))
				buffer[i] = float32(math.Floor(float64(y)))
				// buffer[i] = y
				// buffer[i] = y * rand.Float32()
				sampleTime += 1.0 / SAMPLE_RATE
			}
			rl.UpdateAudioStream(stream, buffer[:]) //colon for slice

			for i := 0; i < len(buffer); i += 16 {
				if i <= rl.GetScreenWidth() {
					rl.DrawCircle(int32(i), int32(100*buffer[i]+225), 5, rl.Blue)
				}
			}

		}

		rl.DrawText(fmt.Sprint("frequency: ", f), 0, TEXT_SIZE*0, TEXT_SIZE, rl.Red)
		rl.DrawText(fmt.Sprint("Note: ", noteNames[note]), 0, TEXT_SIZE*1, TEXT_SIZE, rl.Red)
		rl.EndDrawing()
	}
}
