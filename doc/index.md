English / [Japanese](./index_ja.md)

## Welcome to the hybrid command-line shell

NYAGOS - Nihongo Yet Another GOing Shell is a versatile command-line shell that blends bash-like command-line editing with seamless integration of Windows file system paths and batch files. It offers extensive customization through the Lua scripting language and supports modern predictive input features.

![demo-animation](./demo.gif)

### Key Features

#### UNIX-Like Shell Behavior
- **Keybindings**
  - By default, keybindings are similar to Bash.
  - Customizable via Lua scripts in `%USERPROFILE%\.nyagos`.
    ```lua
    nyagos.key.c_u = "KILL_WHOLE_LINE"
    ```
  - Lua functions can be bound to keys:
    ```lua
    nyagos.key.escape = function(this) nyagos.exec("start vim.exe") end
    ```
- **History and Aliases**
  - Supports `Ctrl-P` history search and `!-style` command recall.
  - Alias system similar to DOSKEY:
    ```lua
    nyagos.alias["g++"] = "g++.exe -std=gnu++17 $*"
    ```
  - Lua-powered aliases:
    ```lua
    nyagos.alias["lala"] = function(args) nyagos.exec("ls", "-al", unpack(args)) end
    ```
- **Custom Command Completion (Bash-Style Tab Completion)**
  - Allows defining completions for specific commands.
    ```lua
    nyagos.complete_for["go"] = function(args)
        if #args == 2 then
            return { "bug", "doc", "fmt", "install", "run", "version",
                     "build", "env", "generate", "list", "test", "vet",
                     "clean", "fix", "get", "mod", "tool" }
        else
            return nil -- file completion
        end
    end
    ```
- **Predictive Completion (PowerShell 7-Like)**
  - Suggests completions based on command history.
  - Predictions can be accepted using `Ctrl-F` or the right arrow key.

#### Windows Compatibility
- **Seamless Batch File Execution**
  - Runs Windows batch files (`.bat` and `.cmd`) as if executed directly in CMD.exe.
  - Captures environment variable changes and directory switches made within batch files.
- **CMD.EXE-Like Features**
  - Supports Windows path formats (`C:\path\to\file`).
  - Maintains a separate current directory for each drive.
  - Includes built-in equivalents for common DOS commands (`copy`, `move`, etc.).
  - No additional DLLs required, and no registry modifications.

#### Enhanced User Experience
- **Colorized Command-Line Interface**
- **Unicode Support**
  - Full compatibility with Windows Unicode APIs.
  - Supports pasting and editing of Unicode characters.
  - Special Unicode literals: `%U+XXXX%` and `$Uxxxx` for prompts.
- **Built-in `ls` Command**
  - Supports colorized output (`-o` option).
  - Displays hard links, symbolic links, and junction targets.
- **Support [SKK] (Simple Kana Kanji conversion program) - [Setup Guide][SKKSetUpEn]**

### Supported Platforms
- Windows 7, 8.1, 10, 11, Windows Server 2008 or later
- Linux (experimental)

[Video by @emisjerry](https://www.youtube.com/watch?v=WsfIrBWwAh0)

### License
You may use, copy, and modify NYAGOS under the New BSD License.

[SKK]: https://ja.wikipedia.org/wiki/SKK
[SKKSetUpEn]: 10-SetupSKK_en.md

### Acknowledgement

[nocd5](https://github.com/nocd5)
/ [mattn](https://github.com/mattn)
/ [hattya](https://github.com/hattya)
/ [shiena](https://github.com/shiena)
/ [atotto](https://github.com/atotto)
/ [ironsand](https://github.com/ironsand)
/ [kardianos](https://github.com/kardianos)
/ [malys](https://github.com/malys)
/ [pine613](https://github.com/pine613)
/ [NSP-0123456](https://github.com/NSP-0123456)
/ [hokorobi](https://github.com/hokorobi)
/ [amuramatsu](https://github.com/amuramatsu)
/ [spiegel-im-spiegel](https://github.com/spiegel-im-spiegel)
/ [rururutan](https://github.com/rururutan/)
/ [hogewest](https://github.com/hogewest)
/ [cagechi](https://github.com/cagechi)
/ [Matsuyanagi](https://github.com/Matsuyanagi)
/ [Shougo](https://github.com/Shougo)
/ [orthographic-pedant](https://github.com/orthographic-pedant)
/ HABATA Katsuyuki
/ [hisomura](https://github.com/hisomura)
/ [tsuyoshicho](https://github.com/tsuyoshicho)
/ [rane-hs](https://github.com/rane-hs)
/ [hami-jp](https://github.com/hami-jp)
/ [3bch](https://github.com/3bch)
/ [AoiMoe](https://github.com/aoimoe)
/ [DeaR](https://github.com/DeaR)
/ [gracix](https://github.com/gracix)
/ [orz--](https://github.com/orz--)
/ [zkangaroo](https://github.com/zkangaroo)
/ [maskedw](https://github.com/maskedw)
/ [tyochiai](https://github.com/tyochiai)
/ [masamitsu-murase](https://github.com/masamitsu-murase)
/ [hazychill](https://github.com/hazychill)
/ [erw7](https://github.com/erw7)
/ [tignear](https://github.com/tignear)
/ [crile](https://github.com/crile)
/ [fushihara](https://github.com/fushihara)
/ [ChiyosukeF](https://twitter.com/ChiyosukeF)
/ [beepcap](https://twitter.com/beepcap)
/ [tostos5963](https://github.com/tostos5963)
/ [sambatriste](https://github.com/sambatriste)
/ [terepanda](https://github.com/terepanda)
/ [Takmg](https://github.com/Takmg)
/ nu8 <!-- https://github.com/nu8 -->
/ [tomato3713](https://github.com/tomato3713)
/ tGqmJHoJKqgK <!-- https://github.com/tGqmJHoJKqgK -->
/ [juggler999](https://github.com/juggler999)
/ [zztkm](https://github.com/zztkm)
/ [8exBCYJi5ATL](https://github.com/8exBCYJi5ATL)
/ [ousttrue](https://github.com/ousttrue)
/ [kgasawa](https://github.com/kgasawa)
/ [HAYASHI-Masayuki](https://github.com/HAYASHI-Masayuki)
/ [naoyaikeda](https://github.com/naoyaikeda)
/ [emisjerry](https://github.com/emisjerry)

### Author

* [hymkor - HAYAMA Kaoru](https://github.com/hymkor) (a.k.a zetamatta)
