English / [Japanese](./03-Readline_ja.md)

## Command Line Editing

You can edit the command line with key bindings similar to UNIX shells.

* `Backspace`, `Ctrl-H` : Delete the character to the left of the cursor
* `Enter`, `Ctrl-M`     : Execute the command line
* `Del`                 : Delete the character under the cursor
* `Home`, `Ctrl-A`      : Move the cursor to the beginning of the line
* `Left`, `Ctrl-B`      : Move the cursor one character to the left
* `Ctrl-D`              : Delete the character under the cursor or exit
* `End`, `Ctrl-E`       : Move the cursor to the end of the line
* `Right`, `Ctrl-F`     : Move the cursor one character to the right
* `Ctrl-K`              : Delete text from the cursor to the end of the line
* `Ctrl-L`              : Refresh the screen
* `Ctrl-U`              : Delete text from the beginning of the line to the cursor
* `Ctrl-Y`              : Paste text from the clipboard
* `Esc`, `Ctrl-[`       : Clear the entire command line
* `Up`, `Ctrl-P`        : Recall the previous command
* `Down`, `Ctrl-N`      : Recall the next command
* `Tab`, `Ctrl-I`       : Complete file or command name
* `Ctrl-C`              : Cancel the current command
* `Ctrl-R`              : Search command history incrementally
* `Ctrl-W`              : Delete the word before the cursor
* `Ctrl-O`              : Insert a filename selected via cursor (requires `box.lua`)
* `Ctrl-XR`, `Alt-R`    : Insert a history entry selected via cursor (requires `box.lua`)
* `Ctrl-XG`, `Alt-G`    : Insert a Git revision selected via cursor (requires `box.lua`)
* `Ctrl-XH`, `Alt-H`    : Insert a `CD`ed directory selected via cursor (requires `box.lua`)
* `Ctrl-Q`, `Ctrl-V`    : Insert the next typed character verbatim
* `Ctrl-Right`, `Alt-F` : Move the cursor forward by one word
* `Ctrl-Left`, `Alt-B`  : Move the cursor backward by one word
* `Ctrl-_, `Ctrl-Z`     : Undo
* `Alt-O`               : Expand a `.lnk` shortcut path to its target path
