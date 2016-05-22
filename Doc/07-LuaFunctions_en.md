## Lua functions extenteded by NYAGOS

### `nyagos.alias.NAME = "REPLACED-STRING"`

It defines an alias. These macros are available on "REPLACE-STRING".

* `$1`,`$2`,`$3` ... the number's argument (removed quotations)
* `$\*` ... all arguments (not removed quotations)

### `nyagos.alias.NAME = function(ARGS)...end`

It assigns the function to the command-name `"NAME"`.
`ARGS` is the table:
    {
        [1]=1stArgument,
        [2]=2ndArgument,
        [3]=3rdArgument,
            :
        ["rawargs"]={
            [1]=1stArgument(not quotatations removed),
            [2]=2ndArgument(not quotatations removed),
            [3]=3rdArgument(not quotatations removed),
                :
        }
    }

When an error occures, the function should return
the number(integer) for %ERRORLEVEL% and error-message.
(No 'return' equals 'return 0,nil')

When the return-value is a string(or string-table), nyagos.exe
executes the string(-table) as a new commandline.

Aliases run on the other Lua-instance and can not access variables
assigned on .nyagos but `share[]`. You can use share[] as you like.
Only the member of the table `share[]` are shared on all Lua-instances 
of nyagos.

### `nyagos.env.NAME`

It is linked to the the environment variable, which are able 
to be refered and assigned.

### `errorlevel,errormessage = nyagos.exec("COMMAND")`

It executes "COMMAND" as shell command.
It returns the integer-value for %ERRORLEVEL% and the error-message.
With no error, they are 0 and nil.

### `errorlevel,errormessage = nyagos.rawexec('COMMAND-NAME','ARG-1','ARG-2'...)`

It executes "COMMAND-NAME" with ARGs. COMMAND-NAME is not interpreted as
a built-in command nor an alias. The difference with os.execute is that
the errormessage is written with utf8.

### `OUTPUT = nyagos.eval("COMMAND")`

It executes "COMMAND" and set its standard output into the lua-variable OUTPUT.
When error occures, OUTPUT is set `nil`.

### `OUTPUT,ERR = nyagos.raweval('COMMAND-NAME','ARG-1','ARG-2'...)`

It executes "COMMAND-NAME" with ARGs and returns commands' standard-output.
COMMAND-NAME is not intepreted as a built-in command nor an alias.

### `WD = nyagos.getwd()`

Get current working directory.

### `nyagos.chdir('DIRECTORY')`

Set new current working directory.

### `nyagos.write("TEXT")`

It output "TEXT" to the standard output with UTF8.

### `nyagos.writerr("TEXT")`

It output "TEXT" to the standard error with UTF8.

### `ANSISTRING = nyagos.utoa(UTF8STRING)`

It converts UTF8 string to the current code page multibyte string.

### `UTF8STRING = nyagos.atou(ANSISTRING)`

It converts the current codepage multibyte string to UTF8 string.

### `FILES = nyagos.glob("WILDCARD-PATTERN1","WILDCARD-PATTERN2"...)`

It returns the table which includes files matching the wildcard pattern(s).

### `path = nyagos.pathjoin('path','to','where'...)`

It makes parts of path-string join.

### `nyagos.bindkey("KEYNAME","FUNCNAME")`

KEYNAME are:
        "C_A" "C_B" ... "C_Z" "M_A" "M_B" ... "M_Z"
        "F1" "F2" ..."F24"
        "BACKSPACE" "CTRL" "DEL" "DOWN" "END"
        "ENTER" "ESCAPE" "HOME" "LEFT" "RIGHT" "SHIFT" "UP"
        "C_BREAK" "CAPSLOCK" "PAGEUP", "PAGEDOWN" "PAUSE"

FUNCNAME are:

        "BACKWARD_DELETE_CHAR" "BACKWARD_CHAR" "CLEAR_SCREEN" "DELETE_CHAR"
        "DELETE_OR_ABORT" "ACCEPT_LINE" "KILL_LINE" "UNIX_LINE_DISCARD"
        "FORWARD_CHAR" "BEGINNING_OF_LINE" "PASS" "YANK" "KILL_WHOLE_LINE"
        "END_OF_LINE" "COMPLETE" "PREVIOUS_HISTORY" "NEXT_HISTORY" "INTR"
        "ISEARCH_BACKWARD"

If it succeeded, it returns true only. Failed, it returns nil and error-message.
Cases are ignores and, the character '-' is same as '\_'.

### `nyagos.bindkey("KEYNAME",function(this)...end)`

When the key is pressed, call the function.

`this` is the table which have these members.

* `this.pos` ... cursor position counted with bytes (==1 when beginning of line)
* `this.text` ... all text represented with utf8
* `this:call("FUNCNAME")` ... call function like `this:call("BACKWARD_DELETE_CHAR")`
* `this:insert("TEXT")` ... insert TEXT at the cursor position.
* `this:firstword()` ... get the first word(=command-name) on the command-line.
* `this:lastword()` ... get the last word and its position on the command-line.
* `this:boxprint({...})` ... listing table values like completion-list.
* `this:replacefrom(POS,"TEXT")` ... replace TEXT between POS and cursor.

The return value of function is used as below

* When it is a string, it is inserted into cursor position.
* When it is `true`, accept line as same as Enter is pressed.
* When it is `false`, drop line as same as Ctrl-C is pressed.
* When it is `nil`, it is ignored.

### `nyagos.filter = function(cmdline) ... end`

`nyagos.filter` can modify user input command-line.
If it returns string, NYAGOS.exe replace the command-line-string it.

### `nyagos.argsfilter = function(args) ... end`

`nyagos.argsfilter` is like `nyaos.filter`, but its argument are
not a string but a table as string array which has each command
arguments.

### `length = nyagos.prompt(template)`

`nyagos.prompt` is assigned function which draw prompt.
You can swap the prompt-function as below.

    nyagos.prompt = function(this)
        local title = "NYAGOS - ".. nyagos.getwd():gsub('\\','/')
        return nyagos.default_prompt('$e[40;36;1m'..this..'$e[37;1m',title)
    end

`nyagos.default_prompt` is the default prompt function which can
change the title of the terminal-window with the second parameter.

### `nyagos.gethistory(N)`

Get the n-th command-line history. When N < 0, last (-N)-th history.
With no arguments, get the count of the command-line history.

### `nyagos.access(PATH,MODE)`

Returns the boolean value whether the PATH can be access with MODE.
It equals the access function of the programming language C.

### `nyagos.completion_hook = function(c) ... end`

This is the Hook for completion. It should be assigned a function.
The argument `c` is the table which has these members.

    c.list[1] .. c.list[#c.list] - command/filename completion result
    c.word - original word without double-quotations.
    c.rawword - original word which may has double-quotations.
    c.pos - position word exists.
    c.text - all command-line text.

`nyagos.completion_hook` should return updated list(table) or `nil`.
Returning nil equals to returning c.list with no change.

### `nyagos.on_command_not_found = function(args) ... end`

It is called when the command which user typed is not found.
The command-name and parameters are set to args[0]...args[#args].
If the function returns nil or false, nyagos.exe prints errors of
usual.

Since the function runs the other Lua-instance, accesss to variables
assigned on .nyagos have the same restriction with aliases. 

### `nyagos.getkey()`

It returns three values : typed key's UNICODE,SCANCODE and SHIFT-Status.

### `WIDTH,HEIGHT=nyagos.getviewwidth()`

It returns the width and height of the terminal.

### `STAT = nyagos.stat(FILENAME)`

It returns the file's information.
If the file exists, the table STAT has these members.

    STAT.name
    STAT.isdir (is set true when the file is directory, otherwise false)
    STAT.size  (bytes)
    STAT.mtime.year
    STAT.mtime.month
    STAT.mtime.day
    STAT.mtime.hour
    STAT.mtime.minute
    STAT.mtime.second

If the file does not exist, STAT is nil.

### `nyagos.open(PATH,MODE)`

Same as io.open but PATH must be written in UTF8.

### `nyagos.loadfile(PATH)`

Same as loadfile on root-namespace but PATH must be written in UTF8.

### `nyagos.lines(PATH)`

Same as io.lines but PATH must be written in UTF8.

```
for text in nyagos.lines(PATH) do ... end
```

`text` is bytearray as same as io.lines().

### `OLEOBJECT = nyagos.create_object('SERVERNAME.TYPENAME')`

Create OLEObject. OLEOBJECTs have methods and property.

- Method
    - `OLEOBJECT:METHOD(PARAMETERS)`.
- Property
    - `OLEOBJECT:_set('PROPERTYNAME',value)`
    - `value = OLEOBJECT:_get('PROPERTYNAME')`

### `nyagos.goversion`

Go-version string to build nyagos.exe
(for example, "go1.6")

### `nyagos.goarch`

The string compilation architecture of nyagos.exe.
(for example, "386" or "amd64" )

### `nyagos.exe`

This string variable has the value of the fullpath of nyagos.exe.
