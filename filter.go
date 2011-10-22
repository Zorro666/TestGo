package main

import "os"
import "fmt"

import "math"
import "image"
import "time"
import "image/draw"
import "exp/gui"
import "exp/gui/x11"

var (
	red = image.NewColorImage(image.RGBAColor{0xFF, 0, 0, 0xFF})
	green = image.NewColorImage(image.RGBAColor{0x00, 0xFF, 0, 0xFF})
	blue = image.NewColorImage(image.RGBAColor{0x00, 0, 0xFF, 0xFF})
	yellow = image.NewColorImage(image.RGBAColor{0xFF, 0xFF, 0x00, 0xFF})
	white = image.NewColorImage(image.RGBAColor{0xFF, 0xFF, 0xFF, 0xFF})
)

var g_xOrigin float64 = 0.0
var g_xScale float64 = 1.0

var g_yOrigin float64 = -1.5
var g_yScale float64 = 200.0

func plotpixel(canvas draw.Image, x float64, y float64, colour image.Color) {
	var ix int = int((x - g_xOrigin)*g_xScale)
	var iy int = 600 - int((y - g_yOrigin)*g_yScale)
	canvas.Set(ix, iy, colour)
}

func render(window gui.Window) {
	var x float64 = 0.0
	var canvas draw.Image = window.Screen()

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
	var vbp float64 = 0.0
	var vhp float64 = 0.0
	var vlp float64 = 0.0
	var w0 float64 = 2.0 * math.Pi * 1.0 * 100.0 / 50000.0
	var Q float64 = 0.8

	for {
		x = x + 1
		var vi float64 = math.Fabs(math.Sin(2.0 * math.Pi * x/100.0))
		var dvbp float64 = w0 * vhp
		var dvlp float64 = w0 * vbp
		vbp -= dvbp
		vlp -= dvlp
		vhp = vbp / Q - vlp - vi

		plotpixel(canvas, x, vi, white)
		plotpixel(canvas, x, vbp, red)
		plotpixel(canvas, x, vlp, green)
		plotpixel(canvas, x, vhp, blue)
		if x > 700 {
			x = x - 700
			window.FlushImage()
			time.Sleep(1)
		}
	}
}

func handleMouseEvent(mouseEvent gui.MouseEvent) {
	fmt.Printf("Mouse Event Buttons 0x%X\n", mouseEvent.Buttons)
}

type Empty interface{}

type MyKeyEvent struct {
	drawKeyEvent gui.KeyEvent
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

func main() {
	var mainWindow gui.Window
	var error os.Error
	mainWindow, error = x11.NewWindow()

	if error != nil {
		fmt.Printf("%s", error.String())
		os.Exit(-1)
	}

	go render(mainWindow)

loop:
	for {
		var windowEvent Empty = <-mainWindow.EventChan()
		switch event := windowEvent.(type) {
		case gui.MouseEvent:
			handleMouseEvent(event)
		case gui.KeyEvent:
			var keyEvent MyKeyEvent
			keyEvent.drawKeyEvent = event
			fmt.Println("Key Event", keyEvent)
			if keyEvent.drawKeyEvent.Key == 65307 { // ESC
				break loop
			}
		case gui.ConfigEvent:
			fmt.Printf("Config Event\n")
		case gui.ErrEvent:
			fmt.Printf("Error Event\n")
			break loop
		}

	}
	error = mainWindow.Close()
}
