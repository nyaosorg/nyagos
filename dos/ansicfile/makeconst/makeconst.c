#include <stdio.h>
#include <fcntl.h>

#define d(n) printf("const " #n "=%d\n", _ ## n)

int main()
{
    printf("package ansicfile\n\n");
    d(O_APPEND);
    d(O_RDONLY);
    d(O_TEXT);
    return 0;
}
