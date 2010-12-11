package main

import "os"
import "fmt"

import "image"
import "time"
import "exp/draw"
import "exp/draw/x11"

//import "./file"

var (
	red = image.NewColorImage(image.RGBAColor{0xFF, 0, 0, 0xFF})
)

func render(window draw.Window) {
	var x float = 0.0
	var y float = 0.0
	var canvas draw.Image = window.Screen()

	for {
		x = x + 1
		y = y + x/4
		if x > 700 {
			x = x - 700
		}
		if y > 600 {
			y = y - 600
		}

		var ix int = int(x)
		var iy int = int(y)
		canvas.Set(ix, iy, red)

		window.FlushImage()
		time.Sleep(1)
	}
}

type Empty interface{}

type MyKeyEvent struct {
	drawKeyEvent draw.KeyEvent
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
	var mainWindow draw.Window
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
		case draw.MouseEvent:
			fmt.Printf("Mouse Event Buttons %d\n", event.Buttons)
		case draw.KeyEvent:
			var keyEvent MyKeyEvent
			keyEvent.drawKeyEvent = event
			fmt.Println("Key Event", keyEvent)
			if keyEvent.drawKeyEvent.Key == 65307 { // ESC
				break loop
			}
		case draw.ConfigEvent:
			fmt.Printf("Config Event\n")
		case draw.ErrEvent:
			fmt.Printf("Error Event\n")
			break loop
		}

	}
	error = mainWindow.Close()
}
