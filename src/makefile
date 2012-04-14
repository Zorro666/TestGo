GC = go build

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

$2_FMTFILES:=$$($2_SRCFILES:.go=.fmt.tmp)

SRCFILES += $$($2_SRCFILES)
FMTFILES += $$($2_FMTFILES)

TARGETS += $1

$1: $$($2_SRCFILES) 
endef
     
$(foreach project,$(PROJECTS),$(eval $(call PROJECT_template,$(project),$(call upperString,$(project)))))

test:
	@echo PROJECTS=$(PROJECTS)
	@echo TARGETS=$(TARGETS)
	@echo SRCFILES=$(SRCFILES)
	@echo FMTFILES=$(FMTFILES)
	@echo WINDOW_SRCFILES=$(WINDOW_SRCFILES)
	@echo WINDOW_FMTFILES=$(WINDOW_FMTFILES)
	@echo TESTFILE_SRCFILES=$(TESTFILE_SRCFILES)
	@echo TESTFILE_FMTFILES=$(TESTFILE_FMTFILES)
	@echo TESTFILE_DEPENDS=$(TESTFILE_DEPENDS)

%: %.go
	$(GC) $<

.PHONY: all clean nuke format
.SUFFIXES:            # Delete the default suffixes

FORCE:

clean: FORCE

nuke: clean
	rm -f $(TARGETS)

%.fmt.tmp: %.go
	gofmt -tabwidth=4 -w=true $<
	@rm -f $@

format: FORCE $(FMTFILES)
