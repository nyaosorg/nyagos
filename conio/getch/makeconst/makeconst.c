#include <stdio.h>
#include <windows.h>

#define d(n) printf("const " #n "=%d\n",n)

int main()
{
    printf("package getch\n\n");
    d(KEY_EVENT);
    return 0;
}
