.POSIX:

VERSION=0.0.1

VPATH=doc
PREFIX?=/usr/local
_INSTDIR=$(DESTDIR)$(PREFIX)
BINDIR?=$(_INSTDIR)/bin
GO?=go
GOFLAGS?=

GOSRC!=find . -name '*.go'
GOSRC+=go.mod

prunef: $(GOSRC)
	$(GO) build $(GOFLAGS) \
		-o $@


all: prunef

RM?=rm -f

clean:
	$(RM) prunef

install: all
	mkdir -m755 -p $(BINDIR)
	install -m755 prunef $(BINDIR)/prunef

RMDIR_IF_EMPTY:=sh -c '\
if test -d $$0 && ! ls -1qA $$0 | grep -q . ; then \
	rmdir $$0; \
fi'

uninstall:
	$(RM) $(BINDIR)/prunef
	${RMDIR_IF_EMPTY} $(BINDIR)

.DEFAULT_GOAL := all

.PHONY: all clean install uninstall
