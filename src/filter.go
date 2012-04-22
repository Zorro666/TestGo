package main

import "log"
import "fmt"

import "math"
import "image"
import "time"
import "image/draw"
import "image/color"
import "jake_graphics"
import "lj_wav_file"

var (
	red = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	green = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	blue = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
	yellow = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
)

var g_debugLevel int = 0
var g_clearImage bool = true
var g_xOrigin float64 = 0.0
var g_xScale float64 = 1.0

var g_yOrigin float64 = -1.5
var g_yScale float64 = 200.0

var g_graphWidth float64 = 0.0
var g_graphHeight float64 = 0.0
var g_graphDisplayMask uint = 0x0

var g_mouseButtonPressed int = 0x0
var g_mouseButtonPressedPosition image.Point
var g_mouseDragging bool = false

var g_viFreq float64 = 1000.0
var g_f0 float64 = 1.0 * 500.0
var g_sampleRate float64 = 48.0*1000.0
var g_Q float64 = 0.8

var g_vbp float64 = 0.0
var g_vhp float64 = 0.0
var g_vlp float64 = 0.0

func plotpixel(canvas draw.Image, x float64, y float64, colour color.Color) {
	var ix int = int((x-g_xOrigin)*g_xScale)
	var iy int = int(g_graphHeight)-int((y-g_yOrigin)*g_yScale)
	canvas.Set(ix, iy, colour)
}

func dragGraph(windowDeltaX int, windowDeltaY int) {
	Log("Old: dx:%d dy:%d origin:%f,%f scale:%f,%f\n", windowDeltaX, windowDeltaY, g_xOrigin, g_yOrigin, g_xScale, g_yScale)

	g_xOrigin -= float64(windowDeltaX)/g_xScale
	g_yOrigin += float64(windowDeltaY)/g_yScale
	g_clearImage = true
	Log("New: origin:%f,%f scale:%f,%f\n", g_xOrigin, g_yOrigin, g_xScale, g_yScale)
}

func zoomGraph(factor float64, zoomWindowCentre image.Point) {
	Log("Old: origin:%f,%f scale:%f,%f\n", g_xOrigin, g_yOrigin, g_xScale, g_yScale)

	var prevXScale float64 = g_xScale
	var prevYScale float64 = g_yScale
	var prevXOrigin float64 = g_xOrigin
	var prevYOrigin float64 = g_yOrigin
	var windowCentreX = float64(zoomWindowCentre.X)
	var windowCentreY = float64(zoomWindowCentre.Y)

	var xCentre float64 = (windowCentreX/prevXScale)+prevXOrigin
	var yCentre float64 = -((windowCentreY-g_graphHeight)/prevYScale)+prevYOrigin

	g_xScale *= factor
	g_yScale *= factor

	g_xOrigin = xCentre - (windowCentreX/g_xScale)
	g_yOrigin = yCentre + ((windowCentreY-g_graphHeight)/g_yScale)

	g_clearImage = true

	Log("New: origin:%f,%f scale:%f,%f\n", g_xOrigin, g_yOrigin, g_xScale, g_yScale)
}

func render(jg *jake_graphics.Jake_Graphics) {

	for {
		if g_clearImage == true {
			var canvas draw.Image = jg.GetBackBuffer()
			renderFrame(canvas)
			jg.FlipBackBuffer()
		}
		time.Sleep(1)
	}
}

func Log(format string, a ...interface{}) {
	if g_debugLevel > 0 {
		log.Printf(format, a...)
	}
	return
}

func renderFrame(canvas draw.Image) {
	var x float64 = 0.0

	if g_clearImage == true {
		var width int = canvas.Bounds().Max.X - canvas.Bounds().Min.X
		var height int = canvas.Bounds().Max.Y - canvas.Bounds().Min.Y
		g_graphWidth = float64(width)
		g_graphHeight = float64(height)

		draw.Draw(canvas, canvas.Bounds(), image.Transparent, image.ZP, draw.Src)
		log.Println("Clear Image", canvas.Bounds())
		log.Printf("viFreq %f f0 %f sampleRate %f Q %f vbp=red vlp=green vhp=blue\n", g_viFreq, g_f0, g_sampleRate, g_Q)
		g_clearImage = false
	}
	// delta_t is converted to seconds given a 1MHz clock by dividing
	// with 1 000 000. This is done in two operations to avoid integer
	// multiplication overflow.

	// Calculate filter outputs.
	// Vhp = Vbp / Q - Vlp - Vi;
	// dVbp = -w0 * Vhp*dt;
	// dVlp = -w0 * Vbp*dt;

	// w0 = cutoff frequency = 2*PI*Freq_in_Hz
	// Q = resonance = 0.707 - 1.707
	// Vi = input volume
	//sound_sample w0_delta_t = w0_ceil_dt*delta_t_flt >> 6;

	//sound_sample dVbp = (w0_delta_t*Vhp >> 14);
	//sound_sample dVlp = (w0_delta_t*Vbp >> 14);
	//Vbp -= dVbp;
	//Vlp -= dVlp;
	//Vhp = (Vbp*_1024_div_Q >> 10) - Vlp - Vi;

	var w0 float64 = 2.0 * math.Pi * g_f0 / g_sampleRate

	var vi_K0 float64 = 2.0 * math.Pi * g_viFreq / g_sampleRate
	var vi_K1 float64 = vi_K0 * 1.5
	var vi_K2 float64 = vi_K0 * 2.0

	var vi_phi0 float64 = 0.0
	var vi_phi1 float64 = 300.0
	var vi_phi2 float64 = 500.0

	var vbp float64 = g_vbp
	var vhp float64 = g_vhp
	var vlp float64 = g_vlp

	for {
		x = x + 1
		var vi0 float64 = math.Abs(math.Sin(vi_K0*(x+vi_phi0)))
		var vi1 float64 = math.Abs(math.Sin(vi_K1*(x+vi_phi1)))
		var vi2 float64 = math.Abs(math.Sin(vi_K2*(x+vi_phi2)))

		var vi float64 = (vi0 + vi1 + vi2)*0.333333333333333

		var dvbp float64 = w0 * vhp
		var dvlp float64 = w0 * vbp
		vbp -= dvbp
		vlp -= dvlp
		vhp = vbp / g_Q - vlp - vi

		if g_graphDisplayMask & 0x1 == 0x1 {
			plotpixel(canvas, x, vi, white)
		}
		if g_graphDisplayMask & 0x2 == 0x2 {
			plotpixel(canvas, x, vbp, red)
		}
		if g_graphDisplayMask & 0x4 == 0x4 {
			plotpixel(canvas, x, vlp, green)
		}
		if g_graphDisplayMask & 0x8 == 0x8 {
			plotpixel(canvas, x, vhp, blue)
		}
		if x > g_graphWidth {
			g_vbp = vbp
			g_vhp = vhp
			g_vlp = vlp
			return
		}
	}
}

func handleMouseButtonEvent(mouseEvent jake_graphics.MouseButtonEvent) {
	Log("Mouse Event Buttons 0x%X Position:%d,%d\n", mouseEvent.Buttons, mouseEvent.X, mouseEvent.Y)

	if mouseEvent.X == 0 && mouseEvent.Y == 0 {
		return
	}

	if mouseEvent.Buttons == 0x0 {
		var factor float64 = 1.0
		// Mouse Button Release
		if g_mouseButtonPressed == 0x1 {
			factor = 2.0
		}
		if g_mouseButtonPressed == 0x4 {
			factor = 0.5
		}

		var deltaX int = mouseEvent.X - g_mouseButtonPressedPosition.X
		var deltaY int = mouseEvent.Y - g_mouseButtonPressedPosition.Y
		if (g_mouseDragging == false) && (deltaX == 0) && (deltaY == 0) && g_mouseButtonPressed != 0x0 {
			zoomGraph(factor, g_mouseButtonPressedPosition)
		}
		g_mouseButtonPressed = 0x0
		g_mouseDragging = false
		g_mouseButtonPressedPosition.X = mouseEvent.X
		g_mouseButtonPressedPosition.Y = mouseEvent.Y
	} else {
		if g_mouseButtonPressed == 0x0 {
			// Store the mouse button location
			g_mouseButtonPressedPosition.X = mouseEvent.X
			g_mouseButtonPressedPosition.Y = mouseEvent.Y
			g_mouseDragging = false
		}
		if (mouseEvent.Buttons == g_mouseButtonPressed) {
			var deltaX int = mouseEvent.X - g_mouseButtonPressedPosition.X
			var deltaY int = mouseEvent.Y - g_mouseButtonPressedPosition.Y
			if (deltaX != 0) || (deltaY != 0) {
				// Dragging
				g_mouseDragging = true
				dragGraph(deltaX, deltaY)
				g_mouseButtonPressedPosition.X = mouseEvent.X
				g_mouseButtonPressedPosition.Y = mouseEvent.Y
			}
		}
		// Store the mouse button pressed buttons
		g_mouseButtonPressed = mouseEvent.Buttons
	}
}

type MyKeyEvent struct {
	drawKeyEvent jake_graphics.KeyEvent
}

func (keyEvent MyKeyEvent) String() string {
	key := keyEvent.drawKeyEvent.Key
	isPressed := "Press"
	if key < 0 {
		isPressed = "Release"
		key = -key
	}

	keyString := "'UNKNOWN'"
	if ' ' <= key && key <= 'z' {
		keyString = fmt.Sprintf("'%c'", key)
	}

	return fmt.Sprintf("%s %s %v 0x%X", isPressed, keyString, key, key)
}

func clampInt(value int, min int, max int) (result int) {
	result = value
	if result < min {
		result = min
	}
	if result > max {
		result = max
	}
	return result
}

func clampFloat(value float64, min float64, max float64) (result float64) {
	result = value
	if result < min {
		result = min
	}
	if result > max {
		result = max
	}
	return result
}

func handleKeyEvent(keyEvent jake_graphics.KeyEvent) {
	var myKeyEvent MyKeyEvent
	myKeyEvent.drawKeyEvent = keyEvent
	log.Println("Key Event", myKeyEvent)
	// Key Release (key -ve)
	if keyEvent.Key < 0 {
		var key uint = uint(-keyEvent.Key)
		var numberKey uint = uint(key) -0x31
		if numberKey < 5 {
			g_graphDisplayMask ^= (1 << numberKey)
			log.Printf("g_graphDisplayMask 0x%X\n", g_graphDisplayMask)
			g_clearImage = true
		}
		var delta_f0 float64 = 10.0
		var delta_Q float64 = 0.05
		var delta_sampleRate float64 = 500.0
		var delta_viFreq float64 = 100.0

		if key == 'r' {
			g_vbp = 0.0
			g_vhp = 0.0
			g_vlp = 0.0
			g_clearImage = true
		}
		if key == 'w' {
			saveWavFile()
		}
		if key == 'd' {
			g_debugLevel += 1
		}
		if key == 'D' {
			g_debugLevel -= 1
		}
		g_debugLevel = clampInt(g_debugLevel, 0, 5)

		if key == 'f' {
			g_f0 += delta_f0
			g_clearImage = true
		}
		if key == 'F' {
			g_f0 -= delta_f0
			g_clearImage = true
		}
		g_f0 = clampFloat(g_f0, 10.0, 30000.0)

		if key == 'q' {
			g_Q += delta_Q
			g_clearImage = true
		}
		if key == 'Q' {
			g_Q -= delta_Q
			g_clearImage = true
		}
		g_Q = clampFloat(g_Q, 0.5, 1.5)

		if key == 's' {
			g_sampleRate += delta_sampleRate
			g_clearImage = true
		}
		if key == 'S' {
			g_sampleRate -= delta_sampleRate
			g_clearImage = true
		}
		g_sampleRate = clampFloat(g_sampleRate, 10.0, 200000.0)

		if key == 'v' {
			g_viFreq += delta_viFreq
			g_clearImage = true
		}
		if key == 'V' {
			g_viFreq -= delta_viFreq
			g_clearImage = true
		}
		g_viFreq = clampFloat(g_viFreq, 100.0, 50000.0)
	}
}

type Empty interface{}

func saveWavFile() {
	result := lj_wav_file.LJ_WAV_ERROR
	filename := "test.wav"
	format := lj_wav_file.LJ_WAV_FORMAT_PCM
	var numChannels uint32 = 1
	var sampleRate uint32 = 44000
	var numBytesPerChannel uint32 = 2
	wavFile := lj_wav_file.LJ_WAV_create(filename, format, numChannels, sampleRate, numBytesPerChannel)
	if wavFile == nil {
		log.Fatalf("ERROR creating wav file:'%s'", filename)
	}
	var sample uint16 = 1000
	result = lj_wav_file.LJ_WAV_FILE_writeChannel(wavFile, sample)
	result = lj_wav_file.LJ_WAV_FILE_writeChannel(wavFile, sample)
	result = lj_wav_file.LJ_WAV_FILE_writeChannel(wavFile, sample)
	result = lj_wav_file.LJ_WAV_FILE_writeChannel(wavFile, sample)
	result = lj_wav_file.LJ_WAV_FILE_writeChannel(wavFile, sample)
	result = lj_wav_file.LJ_WAV_FILE_writeChannel(wavFile, sample)

	result = lj_wav_file.LJ_WAV_close(wavFile)
	if result != lj_wav_file.LJ_WAV_OK {
		log.Fatalf("ERROR during LJ_WAV_close:'%s'", filename)
	}
	log.Printf("Saved WAV file '%s'\n", filename)
}

func main() {
	log.SetFlags(log.Ltime|log.Lmicroseconds)

	jg := jake_graphics.NewInstance()
	jg.CreateWindow(400, 400, 100, 100)

	go render(jg)

loop:
	for {
		windowEvent := jg.WaitForEvent()
		switch event := windowEvent.(type) {
		case jake_graphics.MouseButtonEvent:
			handleMouseButtonEvent(event)
		case jake_graphics.KeyEvent:
			handleKeyEvent(event)
			if event.Key == 9 { // ESC
				break loop
			}
		}

	}
	jg.CloseWindow()
}
