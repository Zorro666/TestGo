package main

import "fmt"

import "time"
import "image/draw"
import "image/color"
import "jake_graphics"

var red = color.RGBA{0xFF, 0, 0, 0xFF}

func render(jg *jake_graphics.Jake_Graphics) {
	var x float32 = 0.0
	var y float32 = 0.0
	var canvas draw.Image = jg.GetBackBuffer()

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

		jg.FlipBackBuffer()
		time.Sleep(1)
	}
}

type Empty interface{}

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

func main() {
	jg := jake_graphics.NewInstance()
	jg.CreateWindow(400, 400, 100, 100)

	go render(jg)

loop:
	for {
		windowEvent := jg.WaitForEvent()
		switch event := windowEvent.(type) {
		  case jake_graphics.MouseButtonEvent:
			  fmt.Printf("Mouse Event Buttons %d %d,%d\n", event.Buttons, event.X, event.Y)
		  case jake_graphics.MouseMoveEvent:
			  fmt.Printf("Mouse Event Motion %d,%d\n", event.X, event.Y)
		  case jake_graphics.KeyEvent:
			  var keyEvent MyKeyEvent
			  keyEvent.drawKeyEvent = event
			  fmt.Println("Key Event", keyEvent)
			  if keyEvent.drawKeyEvent.Key == 0xFF1B { // ESC
				  break loop
			  }
		}
	}
	jg.CloseWindow()
}
