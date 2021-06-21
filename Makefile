PROMPT=$$$$$$S
ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=NUL
    DELTREE=rmdir /s
    SET=set
    D=\\
else
    NUL=/dev/null
    D=/
    SET=export
    DELTREE=rm -r
endif

snapshot: fmt
	$(SET) "CGO_ENABLED=0" && go build -ldflags "-s -w -X main.version=$(shell git.exe describe --tags)"

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

package:
	for /F %%V in ('type Etc\version.txt') do \
	for %%D in (%CD%) do \
	for %%A in (386 amd64) do \
	for %%N in (%%~nD-%%V-windows-%%A.zip) do \
	    zip -9j "%%N" "bin\%%A\nyagos.exe" .nyagos _nyagos makeicon.cmd LICENSE & \
	    zip -9  "%%N" nyagos.d\*.lua nyagos.d\catalog\*.lua
	for /F %%V in ('type Etc\version.txt') do \
	for %%D in (%CD%) do \
	for %%A in (amd64) do \
	    tar -zcvf "nyagos-%%V-linux-%%A.tar.gz" -C .. \
		%%~nD/nyagos \
		%%~nD/.nyagos \
		%%~nD/_nyagos \
		%%~nD/nyagos.d

install:
	@if "%INSTALLDIR%" == "" ( \
	    echo Please do $(MAKE) INSTALLDIR=... & \
	    echo or set INSTALLDIR=... & \
	    exit /b 1 \
	)
	-robocopy  nyagos.d    "$(INSTALLDIR)\nyagos.d" /E
	copy /-Y  _nyagos     "$(INSTALLDIR)\."
	copy /-Y  nyagos.exe  "$(INSTALLDIR)\." || ( \
	    move "$(INSTALLDIR)\nyagos.exe" "$(INSTALLDIR)\nyagos.exe-%RANDOM%" & \
	    copy nyagos.exe  "$(INSTALLDIR)\." )

update:
	for /F "skip=1" %%I in ('where nyagos.exe') do $(MAKE) install INSTALLDIR=%%~dpI
