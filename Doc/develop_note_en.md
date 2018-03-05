* lua: ole:
    * `variable = OLE.property` is avaliable instead of `OLE:_get('property')`
    * `OLE.property = value` is avaliable instead of `OLE:_set('property',value)`
* Load `nyagos.d/*.ny` as batchlike file
* #266: `lua_e "nyagos.option.noclobber = true"` forbides overwriting existing file by redirect.
* #269: `>| FILENAME` and `>! FILENAME` enable to overwrite the file already existing by redirect even if `nyagos.option.noclobber = true`
* #270: Console input buffer has been cleaned up when prompt is drawn.
* #228: Completion supports $ENV[TAB]... by native
* #275: Fix: history substitution like `!str:$` , `!?str?:$` did not work.
* The error `event not found` is caused when the event pointed !y does note exists.
* #285: Not wait GUI-process not using pipeline terminating like CMD.EXE (Call them with ShellExecute() instead of CreateProcess() )
* (Replaced `bytes.Buffer` to `strings.Builder` and Go 1.10 is required to build)
* When more than one are to be executed with `open` at once, display error: `open: ambiguous shellexecute`
