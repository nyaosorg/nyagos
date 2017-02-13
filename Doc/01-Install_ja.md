[English](01-Install_en.md) / Japanese

インストール
------------

バイナリファイルは https://github.com/zetamatta/nyagos/releases よりダウンロード可能です。

    mkdir PATH\TO\INSTALLDIR
    cd PATH\TO\INSTALLDIR
    unzip PATH\TO\DOWNLOADDIR\nyagos-****.zip
    makeicon.cmd

バッチファイル `makeicon.cmd` はデスクトップにアイコンを作成します。

## 簡易カスタマイズ

    notepad %USERPROFILE%\_nyagos

`_nyagos` は DOS のバッチ風の設定ファイルです。
ファイル名が`_`(アンダスコア)で始まる点に注意してください。

例:

    alias "grep=findstr"
    set "GOPATH=%USERPROFILE%\Share\GoSrc"
    suffix "awk=gawk -f"


## Lua によるカスタマイズ

    copy .nyagos "%USERPROFILE%\."
    notepad "%USERPROFILE%\.nyagos"

`%USERPROFILE%\.nyagos` をカスタマイズしてください
