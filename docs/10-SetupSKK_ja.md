SKK のセットアップの仕方
========================

## 1. 辞書をダウンロードします

    cd (YOUR-JISYO-DIR)
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.L
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.emoji

## 2. 辞書のパスを .nyagos に記載します

    nyagos.skk{
        user="~/.go-skk-jisyo" , -- ユーザ辞書
        "(YOUR-JISYO-DIR)/SKK-JISYO.L", -- システム辞書(ラージ)
        "(YOUR-JISYO-DIR)/SKK-JISYO.emoji",-- システム辞書(絵文字)
    }

ユーザ辞書は UTF8 で保存するので、他のアプリのものと違うものにしておいた方が安全です。
