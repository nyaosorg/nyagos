English / [Japanese](./03-Readline_ja.md)

## Command Line Editing

You can edit the command line with key bindings similar to UNIX shells.

* `Backspace`, `Ctrl-H` : Delete the character to the left of the cursor (`"BACKWARD_DELETE_CHAR"`)
* `Enter`, `Ctrl-M`     : Execute the command line (`"ACCEPT_LINE"`)
* `Del`                 : Delete the character under the cursor (`"DELETE_CHAR"`)
* `Home`, `Ctrl-A`      : Move the cursor to the beginning of the line (`"BEGINNING_OF_LINE"`)
* `Left`, `Ctrl-B`      : Move the cursor one character to the left (`"BACKWARD_CHAR"`)
* `Ctrl-D`              : Delete the character under the cursor or exit (`"DELETE_OR_ABORT"`)
* `End`, `Ctrl-E`       : Move the cursor to the end of the line `("END_OF_LINE")`
* `Right`, `Ctrl-F`     : Move the cursor one character to the right (`"FORWARD_CHAR_OR_ACCEPT_PREDICT"`)
* `Ctrl-K`              : Delete text from the cursor to the end of the line (`"KILL_LINE"`)
* `Ctrl-L`              : Refresh the screen (`"CLEAR_SCREEN"`)
* `Ctrl-U`              : Delete text from the beginning of the line to the cursor (`"UNIX_LINE_DISCARD"`)
* `Ctrl-Y`              : Paste text from the clipboard (`"YANK"`)
* `Up`, `Ctrl-P`        : Recall the previous command (`"PREVIOUS_HISTORY"`)
* `Down`, `Ctrl-N`      : Recall the next command (`"NEXT_HISTORY"`)
* `Tab`, `Ctrl-I`       : Complete file or command name (`"COMPLETE"`)
* `Ctrl-C`              : Cancel the current command (`"INTR"`)
* `Ctrl-R`              : Search command history incrementally (`"ISEARCH_BACKWARD"`)
* `Ctrl-W`              : Delete the word before the cursor (`"UNIX_WORD_RUBOUT"`)
* `Ctrl-O`              : Insert a filename selected via cursor (function in `box.lua`)
* `Ctrl-XR`, `Meta-R`   : Insert a history entry selected via cursor (function in `box.lua`)
* `Ctrl-XG`, `Meta-G`   : Insert a Git revision selected via cursor (function in `box.lua`)
* `Ctrl-XH`, `Meta-H`   : Insert a `CD`ed directory selected via cursor (function in `box.lua`)
* `Ctrl-Q`, `Ctrl-V`    : Insert the next typed character verbatim (`"QUOTED_INSERT"`)
* `Ctrl-Right`, `Meta-F`: Move the cursor forward by one word (`"FORWARD_WORD"`)
* `Ctrl-Left`, `Meta-B` : Move the cursor backward by one word (`"BACKWARD_WORD"`)
* `Ctrl-_`, `Ctrl-Z`    : Undo (`"UNDO"`)
* `Meta-O`              : Expand a `.lnk` shortcut path to its target path (function in `box.lua`)
* `Ctrl-T`              : Swap the character at the cursor with the previous one (`"SWAPCHAR"`)
* (Unassigned)          : Paste the clipboard text surrounded by double quotes (`"YANK_WITH_QUOTE"`)
* (Unassigned)          : Delete the entire line (`"KILL_WHOLE_LINE"`)

`Meta` means either `Alt`+`key` or `Esc` followed by key.

### Customizing key bindings

To change key bindings, add statements like the following to your `.nyagos` file:

```lua
-- Make Ctrl-D perform a simple delete instead of exiting the shell.
nyagos.key.C_D = "DELETE_CHAR"
```

Since the `Esc` key is normally used as a prefix key, it cannot be assigned to a command by default. However, if you enable the `singleescape` option, you can bind a command to a standalone `Esc` key.

**(Note: On some terminals, this may occasionally cause the Up Arrow key to be interpreted as a standalone `Esc` followed by `[A`, resulting in incorrect behavior.)**

```lua
nyagos.option.singleescape = true
nyagos.key.escape = "KILL_WHOLE_LINE"
```

<!-- set:fenc=utf8: -->
