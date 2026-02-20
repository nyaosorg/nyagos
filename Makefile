ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=NUL
    DEL=del
    SET=set
    WHICH=where.exe
ifeq ($(shell go env GOOS),windows)
    SYSO=nyagos.syso
else
    SYSO=
endif
else
    NUL=/dev/null
    SET=export
    WHICH=which
    DEL=rm
    SYSO=
endif

ifndef GO
    SUPPORTGO=go1.20.14
    GO:=$(shell $(WHICH) $(SUPPORTGO) 2>$(NUL) || echo go)
endif

NAME:=$(notdir $(CURDIR))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT:=-ldflags "-s -w -X main.version=$(VERSION)"
EXE:=$(shell go env GOEXE)

build:
	$(GO) fmt ./...
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)

define tstlua1
	"./nyagos" --norc -f "$(1)"

endef

define tstlua2
	"$(1)"

endef

test:
	$(GO) test -v ./...
	$(foreach I,$(wildcard test/lua/*.lua),$(call tstlua1,$(I)))
ifeq ($(OS),Windows_NT)
	$(foreach I,$(wildcard test/cmd/*.cmd),$(call tstlua2,$(I)))
endif

clean:
	-$(DEL) nyagos.exe nyagos nyagos.syso 2>$(NUL)

_dist:
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)
	zip -9 "nyagos-$(VERSION)-$(GOOS)-$(GOARCH).zip" \
	    "nyagos$(EXE)" .nyagos \
	    "nyagos.d/catalog/*.lua" \
	    $(FILES)

dist:
	$(SET) "GOPROXY=direct" && $(GO) generate -C Etc
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist "FILES=makeicon.cmd"
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist "FILES=makeicon.cmd"
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist

LATEST=$(GO) run github.com/hymkor/latest-notes@latest -pattern "^\d+\.\d+\.\d+\\?\_\d+$$"
NOTES=doc/CHANGELOG*.md

release:
	$(LATEST) $(NOTES) | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

dry-release:
	$(LATEST) $(NOTES)

bump:
	$(LATEST) -gosrc main -suffix "-goinstall" $(NOTES) > version.go

docs:
	$(MAKE) -C doc all

$(SUPPORTGO):
	go install golang.org/dl/$(SUPPORTGO)@latest
	"$(shell go env GOPATH)/bin/$(SUPPORTGO)" download

.PHONY: build debug test tstlua clean get _dist dist release install docs
