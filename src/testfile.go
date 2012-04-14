package main

import (
	"jake_file"
	"fmt"
	"os"
)

func main() {
	hello := []byte("hello, world\n")
	jake_file.Stdout.Write(hello)
	file, err := jake_file.Open("/does/not/exist", 0, 0)
	if file == nil {
		fmt.Printf("can't open file; err=%s\n", err.Error())
		os.Exit(1)
	}
}
