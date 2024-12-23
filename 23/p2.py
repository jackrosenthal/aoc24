import networkx


def main():
    g = networkx.Graph()
    with open("input.txt") as f:
        for line in f:
            g.add_edge(*line.strip().split("-"))
    clique, _ = networkx.max_weight_clique(g, weight=None)
    print(",".join(sorted(map(str, clique))))


if __name__ == "__main__":
    main()