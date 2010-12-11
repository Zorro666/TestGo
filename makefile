# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ${GOROOT}/src/Make.inc

WINDOW_SRCFILES=window.go
WINDOW_OBJFILES=$(WINDOW_SRCFILES:.go=.8)

TESTFILE_SRCFILES=testfile.go file.go
TESTFILE_OBJFILES=$(TESTFILE_SRCFILES:.go=.8)

CAT_SRCFILES=cat.go file.go
CAT_OBJFILES=$(CAT_SRCFILES:.go=.8)

SRCFILES=\
	$(WINDOW_SRCFILES)\
	$(TESTFILE_SRCFILES)\
	$(CAT_SRCFILES)\

OBJFILES=$(SRCFILES:.go=.8)
FMTFILES=$(SRCFILES:.go=.fmt.tmp)

TARGETS = \
	window\
	testfile\
	cat\

all:$(TARGETS)

window: $(WINDOW_OBJFILES)

testfile: $(TESTFILE_OBJFILES)
testfile.8: file.8 testfile.go 

cat: $(CAT_OBJFILES)
cat.8: file.8 cat.go 

%.8: %.go
	$(GC) $<

%: %.8
	$(LD) -o $@ $<

.PHONY: all clean nuke format
.SUFFIXES:            # Delete the default suffixes

FORCE:

clean: FORCE
	rm -f $(OBJFILES)

nuke: clean
	rm -f $(TARGETS)

%.fmt.tmp: %.go
	gofmt -tabwidth=4 -w=true $<
	@rm -f $@

format: FORCE $(FMTFILES)
