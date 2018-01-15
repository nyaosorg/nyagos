* lua: ole:
    * `variable = OLE.property` is avaliable instead of `OLE:_get('property')`
    * `OLE.property = value` is avaliable instead of `OLE:_set('property',value)`
* Load `nyagos.d/*.ny` as batchlike file
* #266: `lua_e "nyagos.option.noclobber = true"` forbides overwriting existing file by redirect.
* #269: `>| FILENAME` and `>! FILENAME` enable to overwrite the file already existing by redirect even if `nyagos.option.noclobber = true`
* #270: Console input buffer has been cleaned up when prompt is drawn.
