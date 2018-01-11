* lua: ole:
    * `variable = OLE.property` が `OLE:_get('property')` のかわりに使えるようになった
    * `OLE.property = value` が `OLE:_set('property',value)` のかわりに使えるようになった
* `nyagos.d/*.ny` のコマンドファイルも読み込むようにした
* リダイレクトによる既存ファイルの上書きを禁止する `nyagos.option.noclobber` をサポート
