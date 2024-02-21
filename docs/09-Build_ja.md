[top](../readme_ja.md) &gt; [English](./09-Build_en.md) / Japanese

ビルド方法
----------

Git, [Go 1.20.14~](http://go.dev), GNU Make が必要となります

    git clone https://github.com/nyaosorg/nyagos
    cd nyagos
    make

GNU Make がない場合は

    git clone https://github.com/nyaosorg/nyagos
    cd nyagos

    (Windowsの場合)
    go build

    (Linuxの場合)
    CGO_ENABLED=0 go build

Makefile を使わない場合、起動時のバージョン表記に Git コミットがつきません。

<!-- vim:set fenc=utf8: -->
