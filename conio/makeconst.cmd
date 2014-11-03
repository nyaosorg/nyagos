lua %~dp0..\makeconst.lua conio ^
	CTRL_CLOSE_EVENT ^
	CTRL_LOGOFF_EVENT ^
	CTRL_SHUTDOWN_EVENT ^
	STD_OUTPUT_HANDLE ^
	C:\MingW\include > const.go
go fmt const.go

