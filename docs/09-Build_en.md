[top](../readme.md) &gt; English / [Japanese](./09-Build_ja.md)

Build
-----

Git, [Go 1.20.8 or later](http://golang.org) and GNU Make are required

    git clone https://github.com/nyaosorg/nyagos
    cd nyagos
    make

If you do not have GNU Make,

    git clone https://github.com/nyaosorg/nyagos
    cd nyagos

    (for Windows)
    go build

    (for Linux)
    CGO_ENABLED=0 go build

When you do not use Makefile, the displayed version do not contain the git 
commit hash.

<!-- vim:set fenc=utf8: -->
