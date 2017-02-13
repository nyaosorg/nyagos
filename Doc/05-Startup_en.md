English / [Japanese](./05-Startup_ja.md)

## What is done on the Startup

On startup, NYAGOS.exe loads and execute below.

- `(the directory NYAGOS is put)\.nyagos` ... written in Lua
- `(the home directory)\.nyagos` ... written in Lua
- `(the home directory)\_nyagos` ... written in script like a batchfile.

The home directory is the one pointed with %HOME% or %USERPROFILE%.
`_nyagos` does not support FOR , BLOCKed-If, yet.

History are recorded on `%APPDATA%\NYAOS_ORG\nyagos.history`
