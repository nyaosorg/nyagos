[English](./09-Build_en.md) / Japanese

ビルド方法
----------

Git, [Go 1.16~](http://golang.org), GNU Make が必要となります

    git clone https://github.com/zetamatta/nyagos
    cd nyagos
    make

GNU Make がない場合は

    git clone https://github.com/zetamatta/nyagos
    cd nyagos

    (Windowsの場合)
    go build

    (Linuxの場合)
    CGO_ENABLED=0 go build

When you do not use Makefile, the displayed version do not contain the git 
commit hash.

<!-- vim:set fenc=utf8: -->
