package main

import (
		"jake_graphics"
		"fmt"
/*
		"os"
		"code.google.com/p/x-go-binding/xgb"
		"image"
		"image/color"
*/
		)

func main() {
	jg := jake_graphics.NewInstance()
	jg.CreateWindow(400, 400, 100, 100)

	for {
		 fmt.Printf("Flip\n")
		 jg.WaitForEvent()
		 jg.FlipBackBuffer()
	}
}
/*
func main() {
	c, err := xgb.Dial(os.Getenv("DISPLAY"))
	if err != nil {
		fmt.Printf("cannot connect: %v\n", err)
			os.Exit(1)
	}

	win := c.NewId()
	gc := c.NewId()

	c.CreateWindow(0, win, c.DefaultScreen().Root, 150, 150, 200, 200, 0, 0, 0, 0, nil)
	c.ChangeWindowAttributes(win, xgb.CWEventMask, []uint32{xgb.EventMaskExposure | xgb.EventMaskKeyRelease})
	c.CreateGC(gc, win, 0, nil)
	c.MapWindow(win)

	width := 200
	height := 200
	r := image.Rect(0, 0, width, height)
	backbuffer := image.NewNRGBA(r)
	img := backbuffer

	red := color.NRGBA{0xFF, 0, 0, 0x00}

	img.Set(10, 10, red)
	img.Set(20, 20, red)
	img.Set(30, 30, red)
	img.Set(40, 40, red)
	img.Set(50, 50, red)

	var format byte = xgb.ImageFormatZPixmap
	dstX := 0
	dstY := 0
	var leftPad byte = 0
	var depth byte = 24
	var data []byte = backbuffer.Pix

	for {
    reply, err := c.WaitForEvent()
		if err != nil {
		  fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("event %T\n", reply)
		switch reply.(type) {
			case xgb.ExposeEvent:
				c.PutImage(format, win, gc, uint16(width), uint16(height), int16(dstX), int16(dstY), leftPad, depth, data)
			case xgb.KeyReleaseEvent:
				fmt.Printf("key release!\n")
				c.Bell(75)
		}
	}

	c.Close()
}
*/
