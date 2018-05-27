* Support Windows10's native ESCAPE SEQUENCE processing with --no-go-colorable and --enable-virtual-terminal-processing
* For #304,#312, added options to search for the executable from the current directory
    * --look-curdir-first: do before %PATH% (compatible with CMD.EXE)
    * --look-curdir-last : do after %PATH% (compatible with PowerShell)
    * --look-curdir-never: never (compatible with UNIX Shells)
* nyagos.prompt can now be assigned string literal as prompt template directly.
