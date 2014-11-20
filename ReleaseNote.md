NYAGOS 4.0.1\_0
================

* 内蔵 ls の高速化
* 内蔵版 copy/move/del/erase/mkdir/rmdir[/s]を用意
* ビルドに MinGW を必要としなくなった
* ヒストリのサーチポインタを「修正」「Ctrl-C押下」時に初期化するようにした
* ヒストリをリアルタイムにセーブするようにした
* `__コマンド名__` をコマンド名の別名に自動定義
* (エイリアスコマンド) | 

Lua
---
* nyagos.access 関数を追加
* nyagos.shellexecute 関数を追加(open/su の自前実装可能になった)
* nyagos.prompt でプロンプト表示を横取りできるようにした。

Bugfix
------
* リダイレクトでファイルを truncate していなかった
