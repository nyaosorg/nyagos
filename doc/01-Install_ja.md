[English](01-Install_en.md) / Japanese

## インストール

### バイナリをダウンロード

* https://github.com/nyaosorg/nyagos/releases

### Scoop installer を使う場合

```
C:> scoop install nyagos
```

### Chocolatey installer を使う場合

```
C:> choco install nyagos
```


### カスタマイズ

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
