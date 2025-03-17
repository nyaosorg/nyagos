[English](./05-Startup_en.md) / Japanese

## 起動処理

起動時、nyagos.exe は以下のファイルをロード・実行します。

- `(nyagos.exe と同じディレクトリ)\.nyagos` (Luaで記述)
- `(nyagos.exe と同じディレクトリ)\nyagos.d\*.lua` (Luaで記述)
- `(ホームディレクトリ)\.nyagos` (Luaで記述)
- `%APPDATA%\NYAOS_ORG\nyagos.d\*.lua` (luaで記述)

ホームディレクトリとは環境変数 HOME か USERPROFILE の差す先となります。

過去のヒストリ内容を `%APPDATA%\NYAOS_ORG\nyagos.history` から読み出します。
NYAGOS 終了時には、このファイルに再び最後のヒストリ内容が書き出されます。

<!-- set:fenc=utf8: -->
