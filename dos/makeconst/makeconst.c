#include <stdio.h>
#include <windows.h>

#define d(n) printf("const " #n "=%d\n",n)

int main()
{
    printf("package dos\n\n");
    d(FILE_ATTRIBUTE_NORMAL);
    d(FILE_ATTRIBUTE_REPARSE_POINT);
    d(FILE_ATTRIBUTE_HIDDEN);
    d(CP_THREAD_ACP);
    d(MOVEFILE_REPLACE_EXISTING);
    d(MOVEFILE_COPY_ALLOWED);
    d(MOVEFILE_WRITE_THROUGH);
    return 0;
}
