* lua: ole:
    * `variable = OLE.property` が `OLE:_get('property')` のかわりに使えるようになった
    * `OLE.property = value` が `OLE:_set('property',value)` のかわりに使えるようになった
* `nyagos.d/*.ny` のコマンドファイルも読み込むようにした
* #266: `lua_e "nyagos.option.noclobber = true"` でリダイレクトでのファイル上書きを禁止
* #269: `>| FILENAME` で、`nyagos.option.noclobber = true` の時も上書きできるようにした
