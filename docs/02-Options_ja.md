[top](../readme_ja.md) &gt; [English](./02-Options_en.md) / Japanese

## 起動オプション

### --cleanup-buffer (lua: `nyagos.option.cleanup_buffer=true`)
プロンプト表示のタイミングで、キーバッファをクリアします

### --cmd-first "COMMAND"
.nyagos を処理する前に "COMMAND" を実行し、終了後、シェルを継続します。

### --completion-hidden (lua: `nyagos.option.completion_hidden=true`)
ファイル名補完に、隠しファイルも含めます

### --completion-slash (lua: `nyagos.option.completion_slash=true`)
ファイル名補完で、スラッシュを使います。

### --glob (lua: `nyagos.option.glob=true`)
外部コマンドにおいても、ワイルドカード展開を有効にします。

### --help
ヘルプを表示します。

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

### --no-cleanup-buffer (lua: `nyagos.option.cleanup_buffer=false`) [default]
プロンプト表示時にキーバッファをクリアさせません。

### --no-completion-hidden (lua: `nyagos.option.completion_hidden=false`) [default]
ファイル名補完に隠しファイルを含ませません。

### --no-completion-slash (lua: `nyagos.option.completion_slash=false`) [default]
ファイル名補完でスラッシュを使いません(バックスラッシュを使います)

### --no-glob (lua: `nyagos.option.glob=false`) [default]
外部コマンドで、ワイルドカード展開をしません。

### --no-output-surrogate-pair (lua: `nyagos.option.output_surrogate_pair=false`) [default]
サロゲートペアな文字を `<NNNN>` と表記します。

### --no-noclobber (lua: `nyagos.option.noclobber=false`) [default]
リダイレクトでの上書きを許可します。

### --no-read-stdin-as-file (lua: `nyagos.option.read_stdin_as_file=false`) [default]
標準入力からコンソール扱いでコマンドを読み込みます。
(編集機能が有効になります)

### --no-tilde-expansion (lua: `nyagos.option.tilde_expansion=false`)
~ の置換を無効にする

### --no-usesource (lua: `nyagos.option.usesource=false`)
バッチファイルに、NYAGOS側の環境変数の変更させるのを禁止します。

### --noclobber (lua: `nyagos.option.noclobber=true`)
リダイレクトでの上書きを禁止します。

### --norc
`~\.nyagos` , `~\_nyagos` and `(BINDIR)\nyagos.d\*` といった起動スクリプトをロードしないようにします。

### --output-surrogate-pair (lua: `nyagos.option.output_surrogate_pair=true`)
サロゲートペアな文字をそのまま表示します

### --read-stdin-as-file (lua: `nyagos.option.read_stdin_as_file=true`)
標準入力からファイル扱いでコマンドを読み込みます。
(編集機能が無効になります)

### --show-version-only
バージョンを表示します(ビルド用です)

### --tilde-expansion (lua: `nyagos.option.tilde_expansion=true`) [default]
~ 置換を有効にします

### --usesource (lua: `nyagos.option.usesource=true`) [default]
バッチファイルに、NYAGOS側の環境変数の変更させるのを許可します。

### -b "BASE64edCOMMAND"
BASE64形式でエンコードされたコマンドをデコードして実行します。

### -c "COMMAND"
コマンドを実行して、ただちに終了します。

### -e "SCRIPTCODE"
Luaインタプリタでスクリプトコードを実行後、終了します。

### -f FILE ARG1 ARG2 ...
ファイル名の拡張子が .Lua の場合、ファイル中の Lua コードを実行します。
(引数は配列 arg[] という形で参照できます)

さもなければ、通常コマンドとして読み取って実行します。

### -h
オプションのヘルプを表示します。

### -k "COMMAND"
コマンドを実行してから、通常起動します。

<!-- set:fenc=utf8: -->
