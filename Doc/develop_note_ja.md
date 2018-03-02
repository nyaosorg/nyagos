* lua: ole:
    * `variable = OLE.property` が `OLE:_get('property')` のかわりに使えるようになった
    * `OLE.property = value` が `OLE:_set('property',value)` のかわりに使えるようになった
* `nyagos.d/*.ny` のコマンドファイルも読み込むようにした
* #266: `lua_e "nyagos.option.noclobber = true"` でリダイレクトでのファイル上書きを禁止
* #269: `>| FILENAME` もしくは `>! FILENAME` で、`nyagos.option.noclobber = true` の時も上書きできるようにした
* #270: プロンプト表示時にコンソール入力バッファをクリアするようにした
* #228: $ENV[TAB] という補完をネイティブでサポート
* #275: `!str:$` や `!str?:$` といったヒストリ置換が機能しない不具合を修正
* ! で指定されるヒストリが存在しない時「event not found」エラーを出させるようにした
* #285: パイプラインを使っていない GUIプログラムは CMD.EXE 同様終了を待たないようにした (CreateProcess ではなく ShellExecute を使用する)
