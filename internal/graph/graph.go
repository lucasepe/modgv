// This is a modified version of Modgraphviz created by the Go authors.
// Original Modgraphviz resides in the experimental repository.
// https://github.com/golang/exp/tree/master/cmd/modgraphviz

package graph

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

type Edge struct{ From, To string }
type Graph struct {
	Root        string
	Edges       []Edge
	MvsPicked   []string
	MvsUnpicked []string
}

// convert reads “go mod graph” output from r and returns a graph, recording
// MVS picked and unpicked nodes along the way.
func Convert(r io.Reader) (*Graph, error) {
	scanner := bufio.NewScanner(r)
	var g Graph
	seen := map[string]bool{}
	mvsPicked := map[string]string{} // module name -> module version

	for scanner.Scan() {
		l := scanner.Text()
		if l == "" {
			continue
		}

		parts := strings.Fields(l)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected 2 words in line, but got %d: %s", len(parts), l)
		}

		from := parts[0]
		to := parts[1]
		g.Edges = append(g.Edges, Edge{From: from, To: to})

		for _, node := range []string{from, to} {
			if _, ok := seen[node]; ok {
				// Skip over nodes we've already seen.
				continue
			}
			seen[node] = true

			var m, v string
			if i := strings.IndexByte(node, '@'); i >= 0 {
				m, v = node[:i], node[i+1:]
			} else {
				// Root node doesn't have a version.
				g.Root = node
				continue
			}

			if maxV, ok := mvsPicked[m]; ok {
				if semver.Compare(maxV, v) < 0 {
					// This version is higher - replace it and consign the old
					// max to the unpicked list.
					g.MvsUnpicked = append(g.MvsUnpicked, m+"@"+maxV)
					mvsPicked[m] = v
				} else {
					// Other version is higher - stick this version in the
					// unpicked list.
					g.MvsUnpicked = append(g.MvsUnpicked, node)
				}
			} else {
				mvsPicked[m] = v
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for m, v := range mvsPicked {
		g.MvsPicked = append(g.MvsPicked, m+"@"+v)
	}

	// Make this function deterministic.
	sort.Strings(g.MvsPicked)

	return &g, nil
}
