import re
from operator import itemgetter

regex = re.compile("[^abcdefghijklmnñopqrstuvwxyzáéíóúü]")

if __name__ == "__main__":
    with open("pg17013.txt", "r") as f:
        contents = f.read()
    contents = contents.lower()
    contents = regex.sub(" ", contents)
    words = contents.split(" ")
    counts = {}
    for w in words:
        if len(w) > 2:
            if w in counts:
                counts[w] += 1
            else:
                counts[w] = 1

    withcounts = list(counts.items())
    withcounts.sort(key=itemgetter(0))
    withcounts.sort(key=itemgetter(1), reverse=True)
    rank = 1
    for w, c in withcounts:
        print("{} {} {}".format(rank, w, c))
        rank += 1


