- Use Gopher-Lua instead of lua53.dll #300
    - nyagos.exe with lua53.dll can be built with `cd mains ; go build`
    - nyagos.exe with no Lua can be built with `cd ngs ; go build`
- Made `nyagos.option.cleanup_buffer` (default=false). When it is true, clean up console input buffer before readline.
- `set -o OPTION_NAME` and `set +o OPTION_NAME` (=`nyagos.option.OPTION_NAME=` on Lua)
