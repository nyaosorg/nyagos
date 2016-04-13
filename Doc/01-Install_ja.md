インストール
------------

バイナリファイルは https://github.com/zetamatta/nyagos/releases よりダウンロード可能です。

    mkdir PATH\TO\INSTALLDIR
    cd PATH\TO\INSTALLDIR
    unzip PATH\TO\DOWNLOADDIR\nyagos-****.zip
    makeicon.cmd

バッチファイル `makeicon.cmd` はデスクトップにアイコンを作成します。

    copy .nyagos "%USERPROFILE%\."
    notepad "%USERPROFILE%\.nyagos"

`%USERPROFILE%\.nyagos` をカスタマイズしてください
