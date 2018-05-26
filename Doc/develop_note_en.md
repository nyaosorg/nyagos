* Support Windows10's native ESCAPE SEQUENCE processing with --no-go-colorable and --enable-virtual-terminal-processing
* #304 Call an executable on the current dirctory only if the target executable is not found on directories in %PATH%
* Fix: many times C-c typed, C-c and prompt were echoed many times.(by using exec.CommandContext)
