int f(int a, int b)
{
    int x = a + 1;
    int y;
    {
        int x = b + 1;
        y = x;
    }

    return x + y;
}
