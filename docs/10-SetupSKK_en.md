How to setup SKK
================

## 1. Download dictionaries

    cd (YOUR-JISYO-DIR)
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.L
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.emoji

## 2. Describe the paths for the dictionaries on ~/.nyagos

    nyagos.skk{
        user="~/.go-skk-jisyo" , -- user dictionary
        "(YOUR-JISYO-DIR)/SKK-JISYO.L", -- system dictionary (large)
        "(YOUR-JISYO-DIR)/SKK-JISYO.emoji",-- system dictionary (emoji)
    }

You should the different path from other SKK applications for user dictionary because they are saved with UTF8.
