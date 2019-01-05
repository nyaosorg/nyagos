- バッチファイルを呼ぶ時に、`/V:ON` を CMD.EXE に使わないようにした

4.4.0\_beta
===========

- Linux サポート(実験レベル)
- ドライブ毎のカレントディレクトリが子プロセスに継承されなかった問題を修正
- ライブラリ "zetamatta/go-getch" のかわりに "mattn/go-tty" を使うようにした
- msvcrt.dll を直接syscall経由で使わないようにした。
- Linux でも NUL を /dev/null 相当へ
- Lua変数 nyagos.goos を追加
* (#341) Windows10で全角文字の前に文字を挿入すると、不要な空白が入る不具合を修正
    * それに伴い、Windows10 では virtual terminal processing を常に有効に
    * `git.exe push`が無効にしても再び有効にする
* (#339) ワイルドカード `.??*` が `..` にマッチする問題を修正
    * 要 github.com/zetamatta/go-findfile tagged 20181230-2
