package main

import "os"
import "fmt"

import "image"
import "time"
import "exp/draw"
import "exp/draw/x11"

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
		windowEvent := <-mainWindow.EventChan()
		switch event := windowEvent.(type) {
		case draw.MouseEvent:
			fmt.Printf("Mouse Event Buttons %d\n", event.Buttons)
		case draw.KeyEvent:
			fmt.Printf("Key Event %d\n", event.Key)
			if event.Key == 65307 { // ESC
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
