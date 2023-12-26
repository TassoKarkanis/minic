#include <assert.h>

extern int f(void);

int
main() {
    int x = f();
    assert(x == 1);
    return 0;
}
