from mip import Model, INTEGER, xsum, OptimizationStatus
from collections import namedtuple

Group = namedtuple("Group", ["employees", "floors"])

def solve(F, groups):
    allEmployees = 0
    for g in groups:
        allEmployees += g.employees
    G = len(groups)

    print("Solving for {} floors with {} employees in {} groups".format(F, allEmployees, G))

    # A binary search would be much better here, of course, but meh
    for wcs in range(1, allEmployees+1):
        m = Model()
        x = [m.add_var(var_type=INTEGER) for i in range(F*G)]
        for i, g in enumerate(groups):
            m += xsum(x[i*F+f] for f in g.floors) >= g.employees

        for f in range(F):
            m += xsum(x[f+g*F] for g in range(G)) <= wcs

        m.objective = xsum(0 for i in range(F*G))
        m.max_gap = 0.05
        status = m.optimize(max_seconds=3000)
        if status == OptimizationStatus.OPTIMAL:
            print('optimal solution cost {} found'.format(m.objective_value))
        elif status == OptimizationStatus.FEASIBLE:
            print('sol.cost {} found, best possible: {}'.format(m.objective_value, m.objective_bound))
        elif status == OptimizationStatus.NO_SOLUTION_FOUND:
            print('no feasible solution found, lower bound is: {}'.format(m.objective_bound))

        if status == OptimizationStatus.OPTIMAL or status == OptimizationStatus.FEASIBLE:
            print('===================================================================')
            print('solution:')
            for v in m.vars:
                if abs(v.x) > 1e-6: # only printing non-zeros
                    print('{} : {}'.format(v.name, v.x))
            return wcs


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

        print("Case #{}: {}".format(tci+1, solve(F, groups)))