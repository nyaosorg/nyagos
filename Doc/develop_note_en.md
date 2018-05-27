* Support Windows10's native ESCAPE SEQUENCE processing with --no-go-colorable and --enable-virtual-terminal-processing
* Fix: many times C-c typed, C-c and prompt were echoed many times.(by using exec.CommandContext)
* For #304,#312, added options to search for the executable from the current directory
    * --look-curdir-first: do before %PATH% (compatible with CMD.EXE)
    * --look-curdir-last : do after %PATH% (compatible with PowerShell)
    * --look-curdir-never: never (compatible with UNIX Shells)
