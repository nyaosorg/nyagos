[English](./11-FAQ_en.md) / Japanese

## FAQ

### Ctrl-D でシェルを終了させないようにする

UNIX系シェルの ignoreeof のような設定はありませんが、キーバインド変更でシェルを終了しない削除機能を Ctrl-D に割り当てることができます。.nyagos に次の一文を書きます。

    nyagos.key.C_D = "DELETE_CHAR"

### Escape キーでコマンドラインをクリアする

Escape キーはプリフィックスキーとなるため通常はキーを設定できませんが、singleescape という設定を true にすることで単独の Escape に機能を設定できるようになります。
( そのかわり、一部の端末で稀に上矢印キーが Escape 単品と `[A` に分断されたり、誤動作が発生する場合があります )

    nyagos.option.singleescape = true
    nyagos.key.escape = "KILL_WHOLE_LINE"
