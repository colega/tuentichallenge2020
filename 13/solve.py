def f(m, n, h):
    k = h - 2

    s = 0
    sumi = k * (k - 1) // 2
    sumi2 = k * (k - 1) * (2*(k-1) + 1) // 6

    s += sumi * 4 * m
    s += sumi * 4 * n
    s += sumi2 * 16
    s += k * m * n
    s -= k * 2 * m
    s -= k * 2 * n
    s -= 16 * sumi
    s -= 4 * k
    s += (m + 4*k) * (n + 4*k) * 2

    return s


if __name__ == "__main__":
    TC = int(input())
    for tc in range(TC):
        p = int(input())

        i = 0
        j = 1 << 63
        while i < j:
            half = (i+j)//2
            if f(1, 1, half) <= p:
                i = half + 1
            else:
                j = half

        if i < 4:
            print("Case #{}: IMPOSSIBLE".format(tc+1))
            continue

        h = i-1

        i = 0
        j = 1 << 63
        while i < j:
            half = (i+j)//2
            if f(half, half, h) <= p:
                i = half + 1
            else:
                j = half

        m = i-1
        n = m

        if f(m, n+1, h) <= p:
            n += 1

        print("Case #{}: {} {}".format(tc+1, h, f(m, n, h)))
