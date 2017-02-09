## 起動処理

1. 起動時に nyagos.exe と同じフォルダの nyagos.lua を読み込みます。nyagos.lua はLua で記述されており、ここから更にホームディレクトリ(%USERPROFILE%)の .nyagos の Lua コードを読み込みます(nyagos拡張は後述)。ユーザカスタマイズは、この .nyagos を編集して行うことができます。
2. 過去のヒストリ内容を `%APPDATA%\NYAOS_ORG\nyagos.history` から読み出します。NYAGOS 終了時には、このファイルに再び最後のヒストリ内容が書き出されます。

<!-- set:fenc=utf8: -->
