PKGNAME=rmsshkey
PKGCMDDIR=cmd
GOCMD=go
GOBUILD=$(GOCMD) build
SOURCE=main.go
EXECUTABLE=$(PKGNAME)
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
#DEBUILD=debuild
#DEBUILD_ARGS=-us -uc
#DEBUILD_PRESERVE_ENVVARS=GOPATH GOARCH

#ifeq ($(GOARCH),386)
#	DEBUILD_ARGS+=-ai386
#endif

all: clean $(EXECUTABLE)

$(EXECUTABLE):
	$(GOBUILD) -o $@ $(PKGCMDDIR)/$(EXECUTABLE)/$(SOURCE)

clean:
	rm -f $(EXECUTABLE)

install:
	install -d $(DESTDIR)$(BINDIR)
	install -m 0755 -o root -g root $(EXECUTABLE) $(DESTDIR)$(BINDIR)

#deb: clean
#	$(DEBUILD) \
#		$(foreach envvar,$(DEBUILD_PRESERVE_ENVVARS),-e$(envvar)) \
#		$(DEBUILD_ARGS)
