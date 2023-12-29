#include <assert.h>

extern int f(int a0, int a1);

int
main() {
    int x = f(2, 1);
    assert(x == 1);

    x = f(9, 5);
    assert(x == 4);

    x = f(1, 2);
    assert(x == -1);

    x = f(-1, -2);
    assert(x == 1);
    
    return 0;
}
