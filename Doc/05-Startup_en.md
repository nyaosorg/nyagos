[top](../readme.md) &gt; English / [Japanese](./05-Startup_ja.md)

## What is done on the Startup

On startup, NYAGOS.exe loads and execute below.

- `(the directory NYAGOS is put)\.nyagos` (lua-script)
- `(the directory NYAGOS is put)\nyagos.d\*.lua` (lua-script)
- `(the directory NYAGOS is put)\nyagos.d\*.ny` (plain-command-lines)
- `(the home directory)\.nyagos` (lua-script)
- `(the home directory)\_nyagos` (plain-comman-lines)
- `%APPDATA%\NYAOS_ORG\nyagos.d\*.lua` (lua-script)
- `%APPDATA%\NYAOS_ORG\nyagos.d\*.ny` (plain-command-lines)

The home directory is the one pointed with %HOME% or %USERPROFILE%.
`_nyagos` does not support FOR , BLOCKed-If, yet.

History are recorded on `%APPDATA%\NYAOS_ORG\nyagos.history`
