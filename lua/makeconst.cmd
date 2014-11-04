lua %~dp0..\makeconst.lua lua ^
	LUA_TNIL ^
	LUA_TNUMBER ^
	LUA_TBOOLEAN ^
	LUA_TSTRING ^
	LUA_TTABLE ^
	LUA_TFUNCTION ^
	LUA_TUSERDATA ^
	LUA_TTHREAD ^
	LUA_TLIGHTUSERDATA ^
	%~dp0.\include > const.go
go fmt const.go
