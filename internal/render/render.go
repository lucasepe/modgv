// This is a modified version of Modgraphviz created by the Go authors.
// Original Modgraphviz resides in the experimental repository.
// https://github.com/golang/exp/tree/master/cmd/modgraphviz

package render

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/lucasepe/modgv/internal/graph"
	"github.com/lucasepe/modgv/internal/text"
)

// Render translates “go mod graph” output taken from
// the 'in' reader into Graphviz's DOT language, writing
// to the 'out' writer.
func Render(in io.Reader, out io.Writer) error {
	graph, err := graph.Convert(in)
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

	fmt.Fprintf(out, "\t%q [shape=underline style=\"\" fontsize=14 label=<<b>%s</b>>];\n", graph.Root, graph.Root)

	for _, n := range graph.MvsPicked {
		fmt.Fprintf(out, "\t%q [fillcolor=\"#0c5525\" label=<%s>];\n", n, textToHTML(n, "#fafafa"))
	}
	for _, n := range graph.MvsUnpicked {
		fmt.Fprintf(out, "\t%q [fillcolor=\"#a3a3a3\" label=<%s>];\n", n, textToHTML(n, "#0e0e0e"))
	}
	out.Write(edgesAsDOT(graph))

	fmt.Fprintf(out, "}\n")

	return nil
}

// edgesAsDOT returns the edges in DOT notation.
func edgesAsDOT(gr *graph.Graph) []byte {
	var buf bytes.Buffer
	for _, e := range gr.Edges {
		fmt.Fprintf(&buf, "\t%q -> %q", e.From, e.To)
		if _, ok := text.Find(gr.MvsUnpicked, e.To); ok {
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
