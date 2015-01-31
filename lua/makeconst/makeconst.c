#include <stdio.h>
#include "../lua-5.3.0/src/lua.h"
#include "../lua-5.3.0/src/lualib.h"
#include "../lua-5.3.0/src/lauxlib.h"

int main(){
    printf("package lua\n");
    printf("const LUA_REGISTRYINDEX = %d\n",LUA_REGISTRYINDEX);
    printf("const LUA_TBOOLEAN = %d\n",LUA_TBOOLEAN);
    printf("const LUA_TFUNCTION = %d\n",LUA_TFUNCTION);
    printf("const LUA_TLIGHTUSERDATA = %d\n",LUA_TLIGHTUSERDATA);
    printf("const LUA_TNIL = %d\n",LUA_TNIL);
    printf("const LUA_TNUMBER = %d\n",LUA_TNUMBER);
    printf("const LUA_TSTRING = %d\n",LUA_TSTRING);
    printf("const LUA_TTABLE = %d\n",LUA_TTABLE);
    printf("const LUA_TTHREAD = %d\n",LUA_TTHREAD);
    printf("const LUA_TUSERDATA = %d\n",LUA_TUSERDATA);
    return 0;
}
