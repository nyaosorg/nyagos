gcc %~dp0.\makeconst\makeconst.c -o %~dp0.\makeconst\makeconst.exe
%~dp0.\makeconst\makeconst.exe > const.go
go fmt const.go
