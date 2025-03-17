English / [Japanese](./10-SetupSKK_ja.md)

## How to setup SKK

### 1. Download dictionaries

    cd (YOUR-JISYO-DIR)
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.L
    curl -O https://raw.githubusercontent.com/skk-dev/dict/master/SKK-JISYO.emoji

### 2. Describe the paths for the dictionaries on ~/.nyagos

    nyagos.skk{
        user="~/.go-skk-jisyo" , -- user dictionary
        "(YOUR-JISYO-DIR)/SKK-JISYO.L", -- system dictionary (large)
        "(YOUR-JISYO-DIR)/SKK-JISYO.emoji",-- system dictionary (emoji)
        export="GOREADLINESKK",-- export setting to the environment variable GOREADLINESKK
    }

- You should the different path from other SKK applications for user dictionary because they are saved with UTF8.
- `export=` is to share configuration for the child processes using the [go-readline-skk] library via the environment variable.
- Wildcards can be used to specify file names for system dictionaries, but not for user dictionaries.

[go-readline-skk]: https://github.com/nyaosorg/go-readline-skk
