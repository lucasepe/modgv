// This is a modified version of Modgraphviz created by the Go authors.
// Original Modgraphviz resides in the experimental repository.
// https://github.com/golang/exp/tree/master/cmd/modgraphviz

package modgv

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/lucasepe/modgv/text"
)

type RenderOptions struct {
	HighlightModules map[string]bool
}

// Render translates “go mod graph” output taken from
// the 'in' reader into Graphviz's DOT language, writing
// to the 'out' writer.
func Render(in io.Reader, out io.Writer, options RenderOptions) error {
	graph, err := convert(in)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "digraph gomodgraph {\n")
	fmt.Fprintf(out, "\tpad=1;\n")
	fmt.Fprintf(out, "\trankdir=TB;\n")
	fmt.Fprintf(out, "\tranksep=\"1.2 equally\";\n")
	fmt.Fprintf(out, "\tsplines=ortho;\n")
	fmt.Fprintf(out, "\tnodesep=\"0.8\";\n")
	fmt.Fprintf(out, "\tnode [shape=plaintext style=\"filled,rounded\" penwidth=2 fontsize=12 fontname=\"monospace\"];\n")

	fmt.Fprintf(out, "\t%q [shape=underline style=\"\" fontsize=14 label=<<b>%s</b>>];\n", graph.root, graph.root)

	coloring(graph, options)

	for _, n := range graph.mvsPicked {
		if options.HighlightModules[n] {
			fmt.Fprintf(out, "\t%q [fillcolor=\"#ff0000\" label=<%s>];\n", n, textToHTML(n, "#fafafa"))
		} else {
			fmt.Fprintf(out, "\t%q [fillcolor=\"#0c5525\" label=<%s>];\n", n, textToHTML(n, "#fafafa"))
		}
	}
	for _, n := range graph.mvsUnpicked {
		if options.HighlightModules[n] {
			fmt.Fprintf(out, "\t%q [fillcolor=\"#ff0000\" label=<%s>];\n", n, textToHTML(n, "#fafafa"))
		} else {
			fmt.Fprintf(out, "\t%q [fillcolor=\"#a3a3a3\" label=<%s>];\n", n, textToHTML(n, "#0e0e0e"))
		}
	}
	out.Write(edgesAsDOT(graph))

	fmt.Fprintf(out, "}\n")

	return nil
}

// coloring highlight all packages that reference package in options.HighlightModules
func coloring(gr *graph, options RenderOptions) {
	if len(options.HighlightModules) == 0 {
		return
	}
	finished := false
	for !finished {
		coloringCnt := 0
		for _, e := range gr.edges {
			if options.HighlightModules[e.to] && !options.HighlightModules[e.from] {
				options.HighlightModules[e.from] = true
				coloringCnt++
			}
		}
		finished = coloringCnt == 0
	}
}

// edgesAsDOT returns the edges in DOT notation.
func edgesAsDOT(gr *graph) []byte {
	var buf bytes.Buffer
	for _, e := range gr.edges {
		fmt.Fprintf(&buf, "\t%q -> %q", e.from, e.to)
		if _, ok := text.Find(gr.mvsUnpicked, e.to); ok {
			fmt.Fprintf(&buf, "[style=dashed]")
		}
		fmt.Fprintf(&buf, ";\n")
	}
	return buf.Bytes()
}

// textToHTML converts the line with module@version
// in a colored HTML table
func textToHTML(line string, color string) string {
	var mod, ver string
	if i := strings.IndexByte(line, '@'); i >= 0 {
		mod, ver = line[:i], line[i+1:]
	}

	u := fmt.Sprintf(`href="https://pkg.go.dev/%s?tab=doc"`, mod)

	var sb strings.Builder
	sb.WriteString(`<table border="0" cellspacing="8" `)
	if mod != "" {
		sb.WriteString(u)
	}
	sb.WriteString(`>`)
	if len(mod) > 0 {
		sb.WriteString(`<tr><td><font color="`)
		sb.WriteString(color)
		sb.WriteString(`"><b>`)
		sb.WriteString(mod)
		sb.WriteString("</b></font></td></tr>")
	}

	if len(ver) > 0 {
		sb.WriteString(`<tr><td><font color="`)
		sb.WriteString(color)
		sb.WriteString(`" point-size="10">`)
		sb.WriteString(ver)
		sb.WriteString("</font></td></tr>")
	}
	sb.WriteString("</table>")

	return sb.String()
}
