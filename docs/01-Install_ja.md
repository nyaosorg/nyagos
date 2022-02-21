[top](../readme_ja.md) &gt; [English](01-Install_en.md) / Japanese

インストール
------------

バイナリファイルは https://github.com/nyaosorg/nyagos/releases よりダウンロード可能です。

    mkdir PATH\TO\INSTALLDIR
    cd PATH\TO\INSTALLDIR
    unzip PATH\TO\DOWNLOADDIR\nyagos-****.zip
    makeicon.cmd

バッチファイル `makeicon.cmd` はデスクトップにアイコンを作成します。


## カスタマイズ

    copy .nyagos "%USERPROFILE%\."
    notepad "%USERPROFILE%\.nyagos"

`%USERPROFILE%\.nyagos` をカスタマイズしてください

`.nyagos` は Lua で記述する設定ファイルです。
( `.` で始まるファイル名にご注意ください )

### 環境変数の設定

`SET PATH="%PATH%;C:\BIN` に相当する設定をしたい場合、`%USERPROFILE%\.nyagos`
に以下のように記述してください。

    nyagos.env.path = nyagos.env.path .. ";C:\\bin"

### エイリアスの設定

`ls -al` を `lala` と定義するには次のように記述します。

    nyagos.alias.lala = 'ls -al $*'

マクロとして `$1`..`$9`、それらを連結した `$*` が使用できます。

エイリアスとして、Luaで定義した関数も使えます。

    nyagos.alias.lala = function(args)
        nyagos.exec{ "ls","-al", unpack(args) }
    end
