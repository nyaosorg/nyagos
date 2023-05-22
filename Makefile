PROMPT=$$$$$$S
ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=NUL
    DEL=del
    SET=set
ifeq ($(shell go env GOOS),windows)
    SYSO=nyagos.syso
else
    SYSO=
endif
else
    NUL=/dev/null
    SET=export
    DEL=rm
    SYSO=
endif
NAME=$(notdir $(abspath .))
VERSION=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"
EXE=$(shell go env GOEXE)

snapshot:
	go fmt ./...
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

debug:
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT) -tags=debug

define test1
	pushd "$(1)" && go test && popd

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
	go get -u
	go mod tidy

_package:
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	zip -9 "nyagos-$(VERSION)-$(GOOS)-$(GOARCH).zip" \
	    "nyagos$(EXE)" .nyagos _nyagos LICENSE \
	    "nyagos.d/*.lua" "nyagos.d/catalog/*.lua" \
	    $(FILES)

package:
	cd Etc && go generate
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _package "FILES=Etc/*.ico makeicon.cmd"
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _package "FILES=Etc/*.ico makeicon.cmd"
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _package

release:
	gh release create -d --notes "" -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

ifeq ($(OS),Windows_NT)
install:
ifeq ($(INSTALLDIR),)
	@echo Please do $(MAKE) INSTALLDIR=...
	@echo or set INSTALLDIR=...
else
	copy /-Y  _nyagos    "$(INSTALLDIR)\."
	xcopy "nyagos.d\*"   "$(INSTALLDIR)\nyagos.d" /E /I /Y
	copy /-Y  nyagos.exe "$(INSTALLDIR)\." || ( \
	move "$(INSTALLDIR)\nyagos.exe" "$(INSTALLDIR)\nyagos.exe-%RANDOM%" && \
	copy nyagos.exe  "$(INSTALLDIR)\." )
endif

update:
	for /F "skip=1" %%I in ('where nyagos.exe') do $(MAKE) install INSTALLDIR=%%~dpI
endif

.PHONY: snapshot debug test tstlua clean get _package package release install
