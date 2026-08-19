English / [Japanese](./11-FAQ_ja.md)

## FAQ

### Prevent Ctrl-D from exiting the shell

NYAGOS does not provide a setting equivalent to the `ignoreeof` option found in UNIX shells. However, you can remap `Ctrl-D` to the delete-character function so that it no longer exits the shell.

Add the following line to your `.nyagos` file:

    nyagos.key.C_D = "DELETE_CHAR"

### Clear the command line with the Escape key

Normally, the `Escape` key cannot be assigned because it is used as a prefix key. By setting `singleescape` to `true`, however, you can bind a function to a standalone `Escape` key press.

*Note:* On some terminals, this may occasionally cause the Up Arrow key to be split into a standalone `Escape` followed by `[A`, or result in other unexpected behavior.

    nyagos.option.singleescape = true
    nyagos.key.escape = "KILL_WHOLE_LINE"

### Do not propagate environment variable changes from batch files to the calling shell

By default, like `CMD.exe`, when a batch file is executed, changes to environment variables and the current directory made within the batch file are propagated back to nyagos. To suppress this behavior, so that such changes are not reflected in nyagos as they are in shells such as `bash`, add the following line to your `.nyagos` file.

    nyagos.option.usesource = false
