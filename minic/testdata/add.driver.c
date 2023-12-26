#include <assert.h>

extern int f(void);

int
main() {
    int x = f();
    assert(x == 2);
    return 0;
}
