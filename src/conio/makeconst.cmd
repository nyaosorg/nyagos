setlocal
for %%I in (gcc.exe) do set "MINGW=%%~dp$PATH:I.."
lua %~dp0..\makeconst.lua conio ^
	CTRL_CLOSE_EVENT ^
	CTRL_LOGOFF_EVENT ^
	CTRL_SHUTDOWN_EVENT ^
        CTRL_C_EVENT ^
        ENABLE_ECHO_INPUT ^
        ENABLE_PROCESSED_INPUT ^
	STD_INPUT_HANDLE ^
	STD_OUTPUT_HANDLE ^
	KEY_EVENT ^
	"%MINGW%\include" > const.go
go fmt const.go
endlocal
