#include <stdio.h>
#include "../../include/lua.h"
#include "../../include/lualib.h"
#include "../../include/lauxlib.h"

int main(){
    printf("package lua\n");
    putchar('\n');
    printf("const LUA_REGISTRYINDEX = %d\n",LUA_REGISTRYINDEX);
    putchar('\n');
    printf("const LUA_TBOOLEAN = %d\n",LUA_TBOOLEAN);
    printf("const LUA_TFUNCTION = %d\n",LUA_TFUNCTION);
    printf("const LUA_TLIGHTUSERDATA = %d\n",LUA_TLIGHTUSERDATA);
    printf("const LUA_TNIL = %d\n",LUA_TNIL);
    printf("const LUA_TNUMBER = %d\n",LUA_TNUMBER);
    printf("const LUA_TSTRING = %d\n",LUA_TSTRING);
    printf("const LUA_TTABLE = %d\n",LUA_TTABLE);
    printf("const LUA_TTHREAD = %d\n",LUA_TTHREAD);
    printf("const LUA_TUSERDATA = %d\n",LUA_TUSERDATA);
    putchar('\n');
    printf("const LUA_FILEHANDLE = \"%s\"\n",LUA_FILEHANDLE);
    putchar('\n');
    printf("const LUA_OK = %d\n",LUA_OK);
    printf("const LUA_ERRSYNTAX = %d\n",LUA_ERRSYNTAX );
    printf("const LUA_ERRMEM = %d\n",LUA_ERRMEM);
    printf("const LUA_ERRGCMM = %d\n", LUA_ERRGCMM);
    return 0;
}
