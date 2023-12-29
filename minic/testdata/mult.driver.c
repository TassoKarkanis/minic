#include <assert.h>

extern int f(int a0, int a1);

int
main() {
    int x = f(0, 0);
    assert(x == 0);

    x = f(1, 0);
    assert(x == 0);

    x = f(1, 2);
    assert(x == 2);

    x = f(2, 3);
    assert(x == 6);

    x = f(-2, 3);
    assert(x == -6);

    x = f(-2, -3);
    assert(x == 6);
    
    return 0;
}
