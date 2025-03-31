[English](./02-Options_en.md) / Japanese

## 起動オプション

### -b "BASE64edCOMMAND"
BASE64形式でエンコードされたコマンドをデコードして実行します。

### -c "COMMAND"
指定したコマンドを実行し、ただちに終了します。

### --clipboard / --no-clipboard
(Lua: `nyagos.option.clipboard = true` / `false`)

コピーバッファのクリップボード連動を有効 / 無効にする。

### --cmd-first "COMMAND"
.nyagos を処理する前に "COMMAND" を実行し、その後シェルを継続する。

### --completion-hidden / --no-completion-hidden
(Lua: `nyagos.option.completion_hidden = true` / `false`)

ファイル名補完に、隠しファイルも含める。
無効にする場合、`--no-completion-hidden` を使う。

### --completion-slash / --no-completion-slash
(Lua: `nyagos.option.completion_slash = true` / `false`)

ファイル名補完で、フォワードスラッシュ（`/`）を使用する / しない。

### -e "SCRIPTCODE"
Luaインタプリタで指定したスクリプトコードを実行し、終了します。

### -f FILE ARG1 ARG2 ...
ファイル名の拡張子が .Lua の場合、ファイル中の Lua コードを実行します。
(引数は配列 arg[] という形で参照できます)

さもなければ、通常コマンドとして読み取って実行します。

### --glob / --no-glob
(Lua: `nyagos.option.glob = true` / `false`)

外部コマンドにおいても、ワイルドカード展開を有効 / 無効にする。

### --glob-slash / --no-glob-slash
(Lua: `nyagos.option.glob_slash = true` / `false` , `set -o glob_slash` / `set +o glob_slash`)

外部コマンド向けのワイルドカード展開で、ディレクトリの区切り文字として `/` を使う/ `\` を使う

### -h , --help
ヘルプを表示します。

### -k "COMMAND"
コマンドを実行してから、通常起動します。

### --look-curdir-first
カレントディレクトリから実行ファイルを %PATH% より前に探します
(デフォルト:CMD.EXE互換動作)

### --look-curdir-last
カレントディレクトリから実行ファイルを %PATH% より後に探します
(PowerShell互換動作)

### --look-curdir-never
%PATH% にカレントディレクトリが含まれない限り、カレントディレクトリから
実行ファイルは探しません(UNIX Shells互換動作)

### --lua-file FILE ARG1 ARG2...
ファイルを Lua スクリプトとして実行します。
引数を arg[] として参照できます。
バッチファイルへの埋め込みができるように、`@` で始まる行は無視されます

### --lua-first "LUACODE"
.nyagos を読み込む前に、引数の LUAコードを実行します

### --noclobber / --no-noclobber
(Lua: `nyagos.option.noclobber=true` / `false`)

リダイレクト時のファイル上書きを禁止 / 許可します。

### --norc
`~\.nyagos`, `(BINDIR)\.nyagos`, `(BINDIR)\nyagos.d\*.lua`, `%APPDATA%\NYAOS_ORG\nyagos.d\*.lua` といった起動スクリプトをロードしないようにします。

### --output-surrogate-pair / --no-output-surrogate-pair
(Lua: `nyagos.option.output_surrogate_pair=true`)

サロゲートペアの文字をそのまま表示する / `<NNNN>` の形式で出力する。

### --predict / --no-predict
(Lua: `nyagos.option.predict = true` / `false`) [default]

一行入力の予測表示を有効化/無効化します

### --read-stdin-as-file / --no-read-stdin-as-file
(Lua: `nyagos.option.read_stdin_as_file=true` / `false`)

編集機能を無効にし、標準入力をファイルとして読み込みます。
コンソール扱いで読み込む場合は `--no-read-stdin-as-file` を指定します。

### --show-version-only
バージョンを表示します(ビルド用です)

### --tilde-expansion / --no-tilde-expansion
(Lua: `nyagos.option.tilde_expansion=true`) [default]

`~` をホームディレクトリに展開する機能を有効にします。無効にする場合は`--no-tilde-expansion` を指定します。

### --usesource / --no-usesource
(Lua: `nyagos.option.usesource=true` / `false`) [default]

バッチファイルから NYAGOS の環境変数を変更できるようにする / 禁止する。

<!-- set:fenc=utf8: -->
