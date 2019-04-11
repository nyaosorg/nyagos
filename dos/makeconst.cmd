pushd "%~dp0"
go run importconst_run.go -p dos ^
	CONNECT_UPDATE_PROFILE ^
	S_OK
popd
