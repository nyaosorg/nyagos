Latest
======
* 補完リストされる環境変数名を % で囲むようにした。
* ls に -h (help)オプションを追加

BugFix
------
* 存在しないファイルを ls に指定した時、エラーメッセージを出していなかった


NYAGOS 4.0.3\_0
===============

* 環境変数名を補完できるようにした。

BugFix
-------

* `open *.sln` などでワイルドカードがマッチしなかった時、エラーにならなかった
* makeicon.cmd でアイコンがショートカットに紐つかなかった時があった

Trivial
--------
* VBScript の大文字・小文字を修正した(with [vbsfmt](https://github.com/zetamatta/camelfmt))
* license.txt (New BSD License) を用意
* make.cmd sweep で ~ 付きファイルの削除を

NYAGOS 4.0.2\_2
===============

* makeicon.cmd の作成するショートカットの属性を追加
* make resource を実行しなくとも、windres.exe が %PATH% 上にあれば、アイコンリソース(\*.syso)を自動で作成するようにした。

bugfix
------

* デフォルトの .nyagos で定義している nyagos.prompt関数が原因で、画面幅を誤って認識していた問題を修正 (EXE側もエラー処理を追加)


NYAGOS 4.0.2\_1
===============

* ls -1 をサポート
* デスクトップにアイコンを作成するバッチ・VBScript を添付
* 子プロセスで CMD.EXE を起動した時にプロンプトが乱れないよう、初期状態の .nyagos で、%PROMPT からはエスケープシーケンスを削除し、表示時に追加する nyagos.prompt を定義するようにした。
* ビルドする Go を 1.4 にした。

bugfix
------
* 「copy A B」が B が存在する時、実際にコピーしない不具合を修正
* 相対パスでリンクしているジャンクションを、カレントディレクトリが違う時に rmdir で削除できない不具合を修正


NYAGOS 4.0.2\_0
===============

* source で、ディレクトリ移動も取り込むようにした。
* カーソルの移動量から、Unicode 文字の幅を補正するようにした。
* ALT+英字キーに機能をバインドできるようにした。(例: M\_x)
* 2>&1 , 1>&2 などのリダイレクト、|& パイプラインを実装
* echo,rem,which を内蔵コマンド化
* for 文の為に、alias で空白を含まない引数は二重引用符で囲まないようにした
* for 実行中のプロンプトを > だけにした(エイリアス定義変更)

Bugfix
------
* source で、マルチバイト文字列を含む変数を取り込めない不具合を修正

NYAGOS 4.0.1\_0
================

* 内蔵 ls の高速化
* 内蔵版 copy/move/del/erase/mkdir/rmdir[/s]を用意
* ビルドに MinGW を必要としなくなった
* ヒストリを書き換えた時、Ctrl-C 押下時にヒストリ位置をリセットするようにした (#30 & #34 fixed by @nocd5)
* ヒストリをリアルタイムにセーブするようにした
* `__コマンド名__` をコマンド名の別名に自動定義
* F1～F24,PAGEUP,PAGEDOWN 等、サポートキーの追加

Lua
---
* nyagos.access 関数を追加 (pull request #26 by @mattn)
* nyagos.shellexecute 関数を追加(open/su の自前実装可能になった)
* nyagos.prompt でプロンプト表示を横取りできるようにした。

Bugfix
------
* alias + パイプ + & の場合、標準入力から値を受け取れない不具合を修正(#25 reported by @nocd5)
* リダイレクトでファイルを truncate していなかった(#27 reported by @nocd5)
* conio.GetKey の64bit時の不具合を修正 (#32 fixed by @hattya)
