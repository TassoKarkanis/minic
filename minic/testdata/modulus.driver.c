#include <assert.h>

extern int f(int a, int b);

int
main() {
    int x = f(0, 1);
    assert(x == 0);

    x = f(1, 2);
    assert(x == 1);

    x = f(7, 8);
    assert(x == 7);
    
    return 0;
}
