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

snapshot:
	$(GO) fmt ./...
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)

future:
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT) -tags=orgxwidth

debug:
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT) -tags=debug

define test1
	pushd "$(1)" && $(GO) test && popd

endef

test: tstlua
	$(foreach I,$(wildcard internal/*),$(call test1,$(I)))

define tstlua1
	"./nyagos" --norc -f "$(1)"

endef

define tstlua2
	"$(1)"

endef

tstlua:
	$(foreach I,$(wildcard test/lua/*.lua),$(call tstlua1,$(I)))
ifeq ($(OS),Windows_NT)
	$(foreach I,$(wildcard test/cmd/*.cmd),$(call tstlua2,$(I)))
endif

clean:
	-$(DEL) nyagos.exe nyagos nyagos.syso 2>$(NUL)

get:
	$(GO) get -u
	$(GO) mod tidy

_dist:
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)
	zip -9 "nyagos-$(VERSION)-$(GOOS)-$(GOARCH).zip" \
	    "nyagos$(EXE)" .nyagos LICENSE \
	    "nyagos.d/*.lua" "nyagos.d/catalog/*.lua" \
	    $(FILES)

dist:
	cd Etc && $(GO) generate
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist "FILES=Etc/*.ico Etc/*.png makeicon.cmd"
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist "FILES=Etc/*.ico Etc/*.png makeicon.cmd"
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist

release:
	gh release create -d --notes "" -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

$(SUPPORTGO):
	go install golang.org/dl/$(SUPPORTGO)@latest
	$(SUPPORTGO) download

.PHONY: snapshot debug test tstlua clean get _dist dist release install
