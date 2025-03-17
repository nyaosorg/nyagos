[English](./10-SetupSKK_en.md) / Japanese

## SKK のセットアップの仕方

### 1. 辞書をダウンロードします

    cd (YOUR-JISYO-DIR)
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.L
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.emoji

### 2. 辞書のパスを .nyagos に記載します

    nyagos.skk{
        user="~/.go-skk-jisyo" , -- ユーザ辞書
        "(YOUR-JISYO-DIR)/SKK-JISYO.L", -- システム辞書(ラージ)
        "(YOUR-JISYO-DIR)/SKK-JISYO.emoji",-- システム辞書(絵文字)
        export="GOREADLINESKK",-- 環境変数 GOREADLINESKK への設定書き出し
    }

- ユーザ辞書は UTF8 で保存するので、他のアプリのものと違うものにしておいた方が安全です。
- `export=` は go-readline-skk ライブラリを使う子プロセスで環境変数経由で設定を共用するために使います。
- システム辞書のファイル名指定にはワイルドカードが使えますが、ユーザ辞書には使えません。

[go-readline-skk]: https://github.com/nyaosorg/go-readline-skk
