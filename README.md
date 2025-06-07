[![Build status](https://ci.appveyor.com/api/projects/status/bh7866s6oasvchpj?svg=true)](https://ci.appveyor.com/project/zetamatta/nyagos)
[![GoDoc](https://godoc.org/github.com/nyaosorg/nyagos?status.svg)](https://godoc.org/github.com/nyaosorg/nyagos)
[![Go Report Card](https://goreportcard.com/badge/github.com/nyaosorg/nyagos)](https://goreportcard.com/report/github.com/nyaosorg/nyagos)
[![Github latest Releases](https://img.shields.io/github/downloads/nyaosorg/nyagos/latest/total.svg)](https://github.com/nyaosorg/nyagos/releases/latest)

The Nihongo Yet Another GOing Shell
===================================

**&lt;English&gt;** / [&lt;Japanese&gt;](./README_ja.md)

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

[SKK]: https://ja.wikipedia.org/wiki/SKK
[SKKSetUpEn]: doc/10-SetupSKK_en.md

[Video by @emisjerry](https://www.youtube.com/watch?v=WsfIrBWwAh0)

Install
-------

### Download Binary

Download the latest stable release from:

* https://github.com/nyaosorg/nyagos/releases

### Use "Scoop installer"

```cmd
C:> scoop install nyagos
```

### Use "Chocolatey installer"

```cmd
C:> choco install nyagos
```

### Build from source (snapshot)

If you want to try the latest snapshot version, you can install it via `go install`:

```cmd
C:> go install github.com/nyaosorg/nyagos@latest
```

> ⚠️ Note: This builds a development snapshot, not the latest stable release.

Contents
--------

### Release notes

[4.4.x](doc/release_note_en.md)
/ [4.3.x](doc/history-4.3_en.md)
/ [4.2.x](doc/history-4.2_en.md)
/ [4.1.x](doc/history-4.1_en.md)
/ [4.0.x](doc/history-4.0_en.md)

### Documents

1. [Install](doc/01-Install_en.md)
2. [Option for NYAGOS](doc/02-Options_en.md)
3. [Editor](doc/03-Readline_en.md)
4. [Built-in commands](doc/04-Commands_en.md)
5. [What is done on the Startup](doc/05-Startup_en.md)
6. [Substitution](doc/06-Substitution_en.md)
7. [Lua functions extenteded by NYAGOS](doc/07-LuaFunctions_en.md)
8. [Uninstall](doc/08-Uninstall_en.md)
9. [How To build](doc/09-Build_en.md)
10. [How to setup SKK](doc/10-SetupSKK_en.md) (since v4.4.14)

License
-------

You can use, copy and modify under the New BSD License.

Acknowledgement
---------------

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

Author
------

* [hymkor - HAYAMA Kaoru](https://github.com/hymkor) (a.k.a zetamatta)
