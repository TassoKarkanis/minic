#include <assert.h>

extern int f(int a, int b);

int
main() {
    int x = f(0, 1);
    assert(!x);

    x = f(1, 2);
    assert(!x);

    x = f(7, 7);
    assert(x);
    
    x = f(-1, 1);
    assert(!x);
    return 0;
}
