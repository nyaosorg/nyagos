setlocal
call :"%1"
endlocal
exit /b

:""
    go build -ldflags "-s -w"
    exit /b
