#include <stdio.h>
#include "../../include/lua.h"
#include "../../include/lualib.h"
#include "../../include/lauxlib.h"

#define d(n) printf("const " #n "=%d\n",n)
#define s(n) printf("const " #n "=\"%s\"\n",n)

int main(){
    printf("package lua\n\n");
    d(LUA_REGISTRYINDEX);
    d(LUA_TBOOLEAN);
    d(LUA_TFUNCTION);
    d(LUA_TLIGHTUSERDATA);
    d(LUA_TNIL);
    d(LUA_TNUMBER);
    d(LUA_TSTRING);
    d(LUA_TTABLE);
    d(LUA_TTHREAD);
    d(LUA_TUSERDATA);
    s(LUA_FILEHANDLE);
    d(LUA_OK);
    d(LUA_ERRSYNTAX );
    d(LUA_ERRMEM);
    d(LUA_ERRGCMM);
    return 0;
}
