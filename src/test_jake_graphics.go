package main

import "jake_graphics"
import "fmt"

func main() {
	jg := jake_graphics.NewInstance()
	jg.CreateWindow(400, 400, 100, 100)

	for {
		 fmt.Printf("Flip\n")
		 _ := jg.WaitForEvent()
		 jg.FlipBackBuffer()
	}
}
