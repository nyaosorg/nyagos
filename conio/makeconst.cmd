powershell -ExecutionPolicy RemoteSigned -File  "%~dp0..\makeconst.ps1" ^
        conio ^
        CTRL_CLOSE_EVENT ^
        CTRL_LOGOFF_EVENT ^
        CTRL_SHUTDOWN_EVENT ^
        CTRL_C_EVENT ^
        ENABLE_ECHO_INPUT ^
        ENABLE_PROCESSED_INPUT ^
        STD_INPUT_HANDLE ^
        STD_OUTPUT_HANDLE ^
        KEY_EVENT > const.go
go fmt const.go
