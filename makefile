# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ${GOROOT}/src/Make.inc

TESTFILE_DEPENDS:=file.go
CAT_DEPENDS:=file.go
CAT_ROT13_DEPENDS:=file.go
SORTMAIN_DEPENDS:=sort.go
FILTER_DEPENDS:=bufferOutput.go lj_wav_file.go

PROJECTS:=window\
		testfile\
		cat\
		cat_rot13\
		sortmain\
		sieve\
		sieve1\
		server\
		server1\
		filter\

all: $(PROJECTS)

define upperString
$(shell echo $1 | tr [a-z] [A-Z] )
endef

define PROJECT_template
$2_SRCFILES += $1.go
$2_SRCFILES += $($2_DEPENDS)
$2_DEPEND_OBJS:=$($2_DEPENDS:.go=.8)

$2_OBJFILE:=$1.8
$2_OBJFILES:=$$($2_SRCFILES:.go=.8)
$2_FMTFILES:=$$($2_SRCFILES:.go=.fmt.tmp)

SRCFILES += $$($2_SRCFILES)
OBJFILES += $$($2_OBJFILES)
FMTFILES += $$($2_FMTFILES)

TARGETS += $1

$$($2_OBJFILE): $$($2_DEPEND_OBJS) $1.go
$1: $$($2_OBJFILES) 
endef
     
$(foreach project,$(PROJECTS),$(eval $(call PROJECT_template,$(project),$(call upperString,$(project)))))

test:
	@echo PROJECTS=$(PROJECTS)
	@echo TARGETS=$(TARGETS)
	@echo SRCFILES=$(SRCFILES)
	@echo OBJFILES=$(OBJFILES)
	@echo FMTFILES=$(FMTFILES)
	@echo WINDOW_SRCFILES=$(WINDOW_SRCFILES)
	@echo WINDOW_OBJFILES=$(WINDOW_OBJFILES)
	@echo WINDOW_FMTFILES=$(WINDOW_FMTFILES)
	@echo TESTFILE_SRCFILES=$(TESTFILE_SRCFILES)
	@echo TESTFILE_OBJFILES=$(TESTFILE_OBJFILES)
	@echo TESTFILE_FMTFILES=$(TESTFILE_FMTFILES)
	@echo TESTFILE_DEPENDS=$(TESTFILE_DEPENDS)
	@echo TESTFILE_DEPEND_OBJS=$(TESTFILE_DEPEND_OBJS)
	@echo TESTFILE_OBJFILE=$(TESTFILE_OBJFILE)

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
