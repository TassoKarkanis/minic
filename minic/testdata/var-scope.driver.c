#include <assert.h>

extern int f(int a, int b);

int
main() {
    int x = f(1, 2);
    assert(x == 5);

    x = f(2, 3);
    assert(x == 7);
    
    return 0;
}
