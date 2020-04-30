from mip import Model, INTEGER, xsum, OptimizationStatus
from collections import namedtuple

Group = namedtuple("Group", ["employees", "floors"])

def describe(F, groups):
    allEmployees = 0
    for g in groups:
        allEmployees += g.employees
    G = len(groups)

    return "{} floors with {} employees in {} groups".format(F, allEmployees, G)

if __name__ == "__main__":
    testCases = int(input())
    print("{} test cases".format(testCases))
    for tci in range(testCases):
        F, G = [int(s) for s in input().strip().split(' ')]
        m = Model()
        groups = []
        for gi in range(G):
            E, N = [int(s) for s in input().strip().split(' ')]
            floors = sorted([int(s) for s in input().strip().split(' ')])
            group = Group(E, floors)
            groups.append(group)

        print("Case #{}: {}".format(tci+1, describe(F, groups)))