package main

import (
	"jake_file"
	"flag"
	"fmt"
	"os"
)

func cat(f *jake_file.Jake_File) {
	const NBUF = 512
	var buf [NBUF]byte
	for {
		switch nr, er := f.Read(buf[:]); true {
		case nr < 0:
			fmt.Fprintf(os.Stderr, "cat: error reading from %s: %s\n", f.String(), er.Error())
			os.Exit(1)
		case nr == 0: //EOF
			return
		case nr > 0:
			if nw, ew := jake_file.Stdout.Write(buf[0:nr]); nw != nr {
				fmt.Fprintf(os.Stderr, "cat: error writing from %s: %s\n", f.String(), ew.Error())
			}
		}
	}
}

func main() {
	flag.Parse() // Scans the arg list and sets up flags
	if flag.NArg() == 0 {
		cat(jake_file.Stdin)
	}
	for i := 0; i < flag.NArg(); i++ {
		f, err := jake_file.Open(flag.Arg(i), 0, 0)
		if f == nil {
			fmt.Fprintf(os.Stderr, "cat: can't open %s: error %s\n", flag.Arg(i), err)
			os.Exit(1)
		}
		cat(f)
		f.Close()
	}
}
