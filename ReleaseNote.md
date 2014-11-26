Latest
======

* source で、ディレクトリ移動も取り込むようにした。
* カーソルの移動量から、Unicode 文字の幅を補正するようにした。

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
