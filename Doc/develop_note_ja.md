* lua: ole:
    * `variable = OLE.property` が `OLE:_get('property')` のかわりに使えるようになった
    * `OLE.property = value` が `OLE:_set('property',value)` のかわりに使えるようになった