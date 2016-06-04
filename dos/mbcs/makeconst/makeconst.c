#include <stdio.h>
#include <windows.h>

#define d(n) printf("const " #n "=%d\n",n)

int main()
{
    printf("package mbcs\n\n");
    d(CP_THREAD_ACP);
    return 0;
}
