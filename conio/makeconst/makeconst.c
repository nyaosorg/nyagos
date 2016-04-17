#include <stdio.h>
#include <windows.h>

#define d(n) printf("const " #n "=%d\n",n)
#define u(n) printf("const " #n "=uint32(0x%08X)\n",n)

int main()
{
    printf("package conio\n\n");
    d(CTRL_CLOSE_EVENT);
    d(CTRL_LOGOFF_EVENT);
    d(CTRL_SHUTDOWN_EVENT);
    d(CTRL_C_EVENT);
    d(ENABLE_ECHO_INPUT);
    d(ENABLE_PROCESSED_INPUT);
    u(STD_INPUT_HANDLE);
    u(STD_OUTPUT_HANDLE);
    d(KEY_EVENT);
    return 0;
}
