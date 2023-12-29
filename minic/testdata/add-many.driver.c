#include <assert.h>

extern int f(int a, int b, int c, int d, int e, int f);

int
main() {
    int x = f(1, 2, 4, 8, 16, 32);
    assert(x == 63);
    return 0;
}
