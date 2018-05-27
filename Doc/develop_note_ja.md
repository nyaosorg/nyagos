* `--no-go-colorable` と `--enable-virtual-terminal-processing` で、Windows10 ネイティブのエスケープシーケンスをサポート
* #304,#312, カレントディレクトリから実行ファイルを探す時のオプションを追加
    * --look-curdir-first: %PATH% より前に探す(デフォルト:CMD.EXE互換動作)
    * --look-curdir-last : %PATH% より後に探す(PowerShell互換動作)
    * --look-curdir-never: %PATH% だけから実行ファイルを探す(UNIX Shells互換動作)
* nyagos.prompt にプロンプトテンプレートの文字列を直接代入できるようになった。
