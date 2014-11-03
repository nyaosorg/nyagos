lua %~dp0..\makeconst.lua conio ^
	CTRL_CLOSE_EVENT ^
	CTRL_LOGOFF_EVENT ^
	CTRL_SHUTDOWN_EVENT ^
	STD_INPUT_HANDLE ^
	STD_OUTPUT_HANDLE ^
	GENERIC_READ ^
	FILE_SHARE_READ ^
	OPEN_EXISTING ^
	FILE_ATTRIBUTE_NORMAL ^
	KEY_EVENT ^
	C:\MingW\include > const.go
go fmt const.go

