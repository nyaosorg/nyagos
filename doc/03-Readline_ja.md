[English](./03-Readline_en.md) / Japanese

## 編集機能

UNIX系シェルに近いキーバインドで、コマンドラインを編集可能です。

* BackSpace , Ctrl-H : カーソル左の一文字を削除
* Enter , Ctrl-M     : 入力終結
* Del                : カーソル上の一文字を削除
* Home , Ctrl-A      : カーソルを先頭へ移動
* ← , Ctrl-B        : カーソルを一文字左へ移動
* Ctrl-D             : 0文字の時は NYAGOS を終了、さもなければ Del と同じ
* End , Ctrl-E       : カーソルを末尾へ移動
* → , Ctrl-F        : カーソルを一文字右へ移動
* Ctrl-K             : カーソル以降の文字を全て削除し、クリップボードへコピー
* Ctrl-L             : 画面をクリアして、入力した内容を再表示
* Ctrl-U             : カーソルまでの文字を全て削除し、クリップボードへコピー
* Ctrl-Y             : クリップボードの内容を貼り付ける
* Esc , Ctrl-[       : 入力内容を全て削除する
* ↑ , Ctrl-P        : ヒストリ：一つ前の入力内容を展開する
* ↓ , Ctrl-N        : ヒストリ：一つ後の入力内容を展開する
* TAB , Ctrl-I       : ファイル名・コマンド名補完
* Ctrl-C             : 入力内容を破棄
* Ctrl-R             : インクリメンタルサーチ
* Ctrl-W             : カーソル上の単語を削除する
* Ctrl-O             : カーソルで選択したファイル名を挿入する (by box.lua)
* Ctrl-XR , Alt-R    : カーソルで選択したヒストリを挿入する (by box.lua)
* Ctrl-XG , Alt-G    : カーソルで選択したGit Revisionを挿入する(by box.lua)
* Ctrl-XH , Alt-H    : カーソルで選択した過去に移動したディレクトリを挿入する(by box.lua)
* Ctrl-Q , Ctrl-V    : タイプした文字をそのまま挿入する
* Ctrl-Right , Alt-F : 次の単語先頭へ移動
* Ctrl-Left , Alt-B  : 前の単語先頭へ移動
* Ctrl-`_`, Ctrl-Z   : 直前の変更を取り消す
* Alt-O              : ショートカットのパスをリンク先のファイル名に置換

<!-- set:fenc=utf8: -->
