PROMPT=$$$$$$S
ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=NUL
    DELTREE=rmdir /s
    SET=set
    D=\\
    TYPE=type
else
    NUL=/dev/null
    D=/
    SET=export
    DELTREE=rm -r
    TYPE=cat
endif

snapshot: fmt
	$(SET) "CGO_ENABLED=0" && go build -ldflags "-s -w -X main.version=$(shell git.exe describe --tags)"

debug:
	$(SET) "CGO_ENABLED=0" && go build -ldflags "-s -w -X main.version=$(shell git.exe describe --tags)" -tags=debug

release: fmt
	cd bin          2>$(NUL) || mkdir bin
	cd bin$(D)386   2>$(NUL) || mkdir bin$(D)386
	cd bin$(D)amd64 2>$(NUL) || mkdir bin$(D)amd64
	$(SET) "GOOS=windows"  && $(SET) "GOARCH=386"   && go build -o bin/386/nyagos.exe   -ldflags "-s -w"
	$(SET) "GOOS=windows"  && $(SET) "GOARCH=amd64" && go build -o bin/amd64/nyagos.exe -ldflags "-s -w"
	$(SET) "CGO_ENABLED=0" && $(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && go build -ldflags "-s -w"

clean:
	$(DELTREE) bin 2>$(NUL)

fmt:
ifeq ($(OS),Windows_NT)
	- for /F "tokens=2" %%I in ('git status -s ^| more.com ^| findstr /R ".M.*\.go$$" ') do \
	     go fmt "%%I"

else
	git status -s | gawk '/^.M.*\.go/{ system("go fmt " $$2) }'
endif

syso:
	pushd $(MAKEDIR)$(D)Etc && go generate & popd

get:
	go get -u
#	go get -u github.com/zetamatta/go-readline-ny@master
	go mod tidy

VERSION=$(shell $(TYPE) Etc$(D)version.txt)

_zip:
	zip -9j "nyagos-$(VERSION)-windows-$(GOARCH).zip" \
	    "bin$(D)$(GOARCH)$(D)nyagos.exe" .nyagos _nyagos makeicon.cmd LICENSE
	zip -9  "nyagos-$(VERSION)-windows-$(GOARCH).zip" \
	    nyagos.d$(D)*.lua nyagos.d$(D)catalog$(D)*.lua

package:
	make _zip GOARCH=386
	make _zip GOARCH=amd64
	tar zcvf "nyagos-$(VERSION)-linux-amd64.tar.gz" -C .. \
	    nyagos/nyagos nyagos/.nyagos nyagos/_nyagos nyagos/nyagos.d

ifeq ($(OS),Windows_NT)
install:
ifeq ($(INSTALLDIR),)
	@echo Please do $(MAKE) INSTALLDIR=...
	@echo or set INSTALLDIR=...
else
	copy /-Y  _nyagos    "$(INSTALLDIR)$(D)."
	xcopy nyagos.d$(D)*  "$(INSTALLDIR)$(D)nyagos.d" /E /I /Y
	copy /-Y  nyagos.exe "$(INSTALLDIR)$(D)." || ( \
	move "$(INSTALLDIR)$(D)nyagos.exe" "$(INSTALLDIR)$(D)nyagos.exe-%RANDOM%" && \
	copy nyagos.exe  "$(INSTALLDIR)$(D)." )
endif

update:
	for /F "skip=1" %%I in ('where nyagos.exe') do $(MAKE) install INSTALLDIR=%%~dpI
endif
