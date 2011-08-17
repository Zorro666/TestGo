package main

import "image/draw"

var (
	red = image.NewColorImage(image.RGBAColor{0xFF, 0, 0, 0xFF})
)

func render(window draw.Window) {
	var canvas draw.Image = window.Screen()

	canvas.Set(0, 0, red)
}
