lua %~dp0..\makeconst.lua dos ^
        FILE_ATTRIBUTE_NORMAL ^
	FILE_ATTRIBUTE_REPARSE_POINT ^
	FILE_ATTRIBUTE_HIDDEN ^
	CP_THREAD_ACP ^
	C:\MingW\include > const.go
go fmt const.go
