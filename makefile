# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ${GOROOT}/src/Make.inc

all:window

SRCFILES=window.go file.go
OBJFILES=$(SRCFILES:.go=.8)

window: $(OBJFILES)

window.8: file.8 window.go 

TARG=window

%.8: %.go
	$(GC) $<

%: %.8
	$(LD) -o $@ $<

clean:
	rm -f $(OBJFILES)

nuke: clean
	rm -f $(TARG)
