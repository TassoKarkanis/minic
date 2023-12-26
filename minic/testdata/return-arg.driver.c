#include <assert.h>

extern int f(int);

int
main() {
    int i;
    for (i=0; i<10; i++) {
        int x = f(i);
        assert(x == i);
    }
    return 0;
}
