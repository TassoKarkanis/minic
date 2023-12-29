#include <assert.h>

extern int f(int a, int b, int c);

int
main() {
    int x = f(1, 2, 3);
    assert(x == 5);

    x = f(2, 3, 4);
    assert(x == 10);
    
    return 0;
}
