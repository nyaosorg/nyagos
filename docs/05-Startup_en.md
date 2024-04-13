[top](../README.md) &gt; English / [Japanese](./05-Startup_ja.md)

## What is done on the Startup

On startup, NYAGOS.exe loads and execute below.

- `(the directory NYAGOS is put)\.nyagos` (lua-script)
- `(the directory NYAGOS is put)\nyagos.d\*.lua` (lua-script)
- `(the home directory)\.nyagos` (lua-script)
- `%APPDATA%\NYAOS_ORG\nyagos.d\*.lua` (lua-script)

The home directory is the one pointed with %HOME% or %USERPROFILE%.

History are recorded on `%APPDATA%\NYAOS_ORG\nyagos.history`
