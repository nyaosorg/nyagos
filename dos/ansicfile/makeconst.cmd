pushd "%~dp0"
gcc makeconst\makeconst.c && a > const.go && go fmt const.go
if exist a.exe del a.exe
popd
