beats = {"R": "S", "S": "P", "P": "R"}

if __name__ == "__main__":
    TC = int(input())
    for i in range(TC):
        [a, b] = input().strip().split(' ')
        if a == b:
            winner = "-"
        elif beats[a] == b:
            winner = a
        else:
            winner = b
        print("Case #{}: {}".format(i+1, winner))


