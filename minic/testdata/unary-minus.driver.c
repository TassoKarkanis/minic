#include <assert.h>

extern int f(int a);

int
main() {
    int x = f(2);
    assert(x == -2);

    x = f(0);
    assert(x == 0);

    x = f(-6);
    assert(x == 6);
    
    return 0;
}
