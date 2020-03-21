.POSIX:

VERSION=0.0.1

VPATH=doc
PREFIX?=/usr/local
_INSTDIR=$(DESTDIR)$(PREFIX)
BINDIR?=$(_INSTDIR)/bin
MANDIR?=$(_INSTDIR)/share/man
GO?=go
GOFLAGS?=

GOSRC!=find . -name '*.go'
GOSRC+=go.mod

prunef: $(GOSRC)
	$(GO) build $(GOFLAGS) \
		-o $@

doc:
	scdoc < prunef.1.scd > prunef.1

all: prunef doc

RM?=rm -f

clean:
	$(RM) prunef prunef.1

install: all
	mkdir -m755 -p $(BINDIR)
	mkdir -m755 -p $(MANDIR)
	mkdir -m755 -p $(MANDIR)/man1
	install -m755 prunef $(BINDIR)/prunef
	install -m644 prunef.1 $(MANDIR)/man1/prunef.1

RMDIR_IF_EMPTY:=sh -c '\
if test -d $$0 && ! ls -1qA $$0 | grep -q . ; then \
	rmdir $$0; \
fi'

uninstall:
	$(RM) $(BINDIR)/prunef
	${RMDIR_IF_EMPTY} $(BINDIR)
	$(RM) $(MANDIR)/man1/prunef.1
	${RMDIR_IF_EMPTY} $(MANDIR)/man1
	${RMDIR_IF_EMPTY} $(MANDIR)

test:
	$(GO) test ./...

.DEFAULT_GOAL := all

.PHONY: all clean install uninstall
