ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=NUL
    DEL=del
    SET=set
    GITDIR=$(or $(GIT_INSTALL_ROOT),$(shell for %%I in (git.exe) do echo %%~dp$$PATH:I..))
    AWK="$(GITDIR)\usr\bin\gawk.exe"
ifeq ($(shell go env GOOS),windows)
    SYSO=nyagos.syso
else
    SYSO=
endif
else
    NUL=/dev/null
    SET=export
    DEL=rm
    AWK=gawk
    SYSO=
endif
NAME=$(notdir $(abspath .))
VERSION=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"
EXE=$(shell go env GOEXE)

snapshot: fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

debug:
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT) -tags=debug

test: tstlua
	$(foreach I,$(wildcard internal/*),pushd "$(I)" &&  go test && popd && ) echo OK

tstlua:
	$(foreach I,$(wildcard t/lua/*.lua),echo $(I) && "./nyagos" --norc -f "$(I)" && ) echo OK
ifeq ($(OS),Windows_NT)
	$(foreach I,$(wildcard t/cmd/*.cmd),echo $(I) && "$(I)" && ) echo OK
endif

clean:
	-$(DEL) nyagos.exe nyagos nyagos.syso 2>$(NUL)

fmt:
	$(foreach I,$(shell git status -s | $(AWK) "/^.M.*\.go$$/{ print $$NF }"),gofmt -w $(I) && ) echo OK

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

.PHONY: snapshot debug test tstlua clean fmt get _package package release
