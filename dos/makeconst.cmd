powershell -ExecutionPolicy RemoteSigned -File "%~dp0..\makeconst.ps1" ^
        dos ^
        FILE_ATTRIBUTE_NORMAL ^
        FILE_ATTRIBUTE_REPARSE_POINT ^
        FILE_ATTRIBUTE_HIDDEN ^
        CP_THREAD_ACP ^
        MOVEFILE_REPLACE_EXISTING ^
        MOVEFILE_COPY_ALLOWED ^
        MOVEFILE_WRITE_THROUGH > const.go
go fmt const.go
