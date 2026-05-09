NAME=$(lastword $(subst /, ,$(abspath .)))
VERSION=$(shell git.exe describe --tags)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"
ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    DEL=del
else
    SET=export
    DEL=rm
endif

ifndef GO
    SUPPORTGO=go1.20.14
    GO:=$(shell $(WHICH) $(SUPPORTGO) 2>$(NUL) || echo go)
endif

all:
	$(GO) fmt
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)

test:
	$(GO) test -v

_package:
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT) && \
	zip -9 $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXE)

package:
	$(SET) "GOOS=linux" && $(SET) "GOARCH=386"   && $(MAKE) _package EXE=
	$(SET) "GOOS=linux" && $(SET) "GOARCH=amd64" && $(MAKE) _package EXE=
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _package EXE=.exe
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _package EXE=.exe

clean:
	$(DEL) *.zip $(NAME) $(NAME).exe
