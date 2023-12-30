#include <assert.h>

extern int f(int a0, int a1);

int
main() {
    int x = f(0, 1);
    assert(x == 0);

    x = f(1, 1);
    assert(x == 1);

    x = f(1, 2);
    assert(x == 0);

    x = f(2, 2);
    assert(x == 1);

    x = f(16, 4);
    assert(x == 4);

    x = f(-8, 2);
    assert(x == -4);

    x = f(8, -2);
    assert(x == -4);

    x = f(-8, -2);
    assert(x == 4);
    
    return 0;
}
