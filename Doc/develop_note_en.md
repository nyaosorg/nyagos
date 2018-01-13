* lua: ole:
    * `variable = OLE.property` is avaliable instead of `OLE:_get('property')`
    * `OLE.property = value` is avaliable instead of `OLE:_set('property',value)`
* Load `nyagos.d/*.ny` as batchlike file
* #266: `lua_e "nyagos.option.noclobber = true"` forbides overwriting existing file by redirect.
* #269: `>| FILENAME` enables to overwrite the file already existing by redirect even if `nyagos.option.noclobber = true`
