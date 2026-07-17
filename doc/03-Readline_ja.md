[English](./03-Readline_en.md) / Japanese

## コマンドライン編集

UNIX系シェルに近いキーバインドで、コマンドラインを編集可能です。

* `Backspace`, `Ctrl-H` : カーソル左の一文字を削除 (`"BACKWARD_DELETE_CHAR"`)
* `Enter`, `Ctrl-M`     : 入力終結 (`"ACCEPT_LINE"`)
* `Del`                 : カーソル上の一文字を削除 (`"DELETE_CHAR"`)
* `Home`, `Ctrl-A`      : カーソルを先頭へ移動 (`"BEGINNING_OF_LINE"`)
* `←`, `Ctrl-B`        : カーソルを一文字左へ移動 (`"BACKWARD_CHAR"`)
* `Ctrl-D`              : 0文字の時は NYAGOS を終了、さもなければ `Del` と同じ (`"DELETE_OR_ABORT"`)
* `End`, `Ctrl-E`       : カーソルを末尾へ移動 `("END_OF_LINE")`
* `→`, `Ctrl-F`        : カーソルが行末の時は予測候補確定、さもなければカーソル一文字分移動 (`"FORWARD_CHAR_OR_ACCEPT_PREDICT"`)
* `Ctrl-K`              : カーソル以降の文字を全て削除し、クリップボードへコピー (`"KILL_LINE"`)
* `Ctrl-L`              : 画面をクリアして、入力した内容を再表示 (`"CLEAR_SCREEN"`)
* `Ctrl-U`              : カーソルまでの文字を全て削除し、クリップボードへコピー (`"UNIX_LINE_DISCARD"`)
* `Ctrl-Y`              : クリップボードの内容を貼り付ける (`"YANK"`)
* `↑`, `Ctrl-P`        : ヒストリ：一つ前の入力内容を展開する (`"PREVIOUS_HISTORY"`)
* `↓`, `Ctrl-N`        : ヒストリ：一つ後の入力内容を展開する (`"NEXT_HISTORY"`)
* `Tab`, `Ctrl-I`       : ファイル名・コマンド名補完 (`"COMPLETE"`)
* `Ctrl-C`              : 入力内容を破棄 (`"INTR"`)
* `Ctrl-R`              : インクリメンタルサーチ (`"ISEARCH_BACKWARD"`)
* `Ctrl-W`              : カーソル上の単語を削除する (`"UNIX_WORD_RUBOUT"`)
* `Ctrl-O`              : カーソルで選択したファイル名を挿入する (function in `box.lua`)
* `Ctrl-XR`, `Meta-R`   : カーソルで選択したヒストリを挿入する (function in `box.lua`)
* `Ctrl-XG`, `Meta-G`   : カーソルで選択した Git Revision を挿入する (function in `box.lua`)
* `Ctrl-XH`, `Meta-H`   : カーソルで選択した過去に移動したディレクトリを挿入する (function in `box.lua`)
* `Ctrl-Q`, `Ctrl-V`    : タイプした文字をそのまま挿入する (`"QUOTED_INSERT"`)
* `Ctrl-Right`, `Meta-F`: 次の単語先頭へ移動 (`"FORWARD_WORD"`)
* `Ctrl-Left`, `Meta-B` : 前の単語先頭へ移動 (`"BACKWARD_WORD"`)
* `Ctrl-_`, `Ctrl-Z`    : 直前の変更を取り消す (`"UNDO"`)
* `Meta-O`              : ショートカットのパスをリンク先のファイル名に置換 (function in `box.lua`)
* `Ctrl-T`              : カーソルとその前の文字を入れ替える (`"SWAPCHAR"`)
* (キー未設定) : 二重引用符で囲んでペーストする (`"YANK_WITH_QUOTE"`)
* (キー未設定) : 行全体を削除する (`"KILL_WHOLE_LINE"`)

※ `Meta`は`Alt`+`key`もしくは、`Esc` の後に`key`を押下することを意味します。

### キーのカスタマイズ

キーバインドを変更する場合は .nyagos に次のような文を書きます

```lua
-- Ctrl-D をシェルを終了する機能のない純粋な削除にする
nyagos.key.C_D = "DELETE_CHAR"
```

Escape キーはプリフィックスキーとなるため通常はキーを設定できませんが、singleescape という設定を true にすることで単独の Escape に機能を設定できるようになります。
( そのかわり、一部の端末で稀に上矢印キーが Escape 単品と `[A` に分断されたり、誤動作が発生する場合があります )

```lua
nyagos.option.singleescape = true
nyagos.key.escape = "KILL_WHOLE_LINE"
```

<!-- set:fenc=utf8: -->
