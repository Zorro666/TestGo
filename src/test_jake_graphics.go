package main

import (
		"jake_graphics"
		"fmt"
		)

func main() {
	jg := jake_graphics.NewInstance()
	jg.CreateWindow(300, 400, 100, 100)

	for {
		 fmt.Printf("Flip\n")
//		 jg.FlipBackBuffer()
	}
}
