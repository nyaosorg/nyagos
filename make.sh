#!/bin/bash

VER=`gawk '{ gsub(/\r/,"");print $1 }' ./Misc/version.txt`

case "$1" in
    ""|build)
        go build -ldflags "-s -w -X main.version=$VER"
        ;;
    package)
        ARCH=`go env GOARCH`
        ( cd .. && tar -zcvf nyagos/nyagos-$VER-linux-$ARCH.tar.gz \
            nyagos/nyagos \
            nyagos/.nyagos \
            nyagos/_nyagos \
            nyagos/readme.md \
            nyagos/readme_ja.md \
            nyagos/nyagos.d/ \
            nyagos/Doc/*.md )
        ;;
    *)
        echo Usage:
        echo   $0 \"\"
        echo   $0 build
        echo   $0 package
esac
