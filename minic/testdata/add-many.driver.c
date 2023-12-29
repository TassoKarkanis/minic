#include <assert.h>

extern int f(int a0, int a1, int a2, int a3, int a4, int a5);

int
main() {
    int x = f(1, 2, 4, 8, 16, 32);
    assert(x == 63);
    return 0;
}
