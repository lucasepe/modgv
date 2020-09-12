package modgv

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestMVSPicking(t *testing.T) {
	for _, tc := range []struct {
		name         string
		in           []string
		wantPicked   []string
		wantUnpicked []string
	}{
		{
			name:         "single node",
			in:           []string{"foo@v0.0.1"},
			wantPicked:   []string{"foo@v0.0.1"},
			wantUnpicked: nil,
		},
		{
			name:         "duplicate same node",
			in:           []string{"foo@v0.0.1", "foo@v0.0.1"},
			wantPicked:   []string{"foo@v0.0.1"},
			wantUnpicked: nil,
		},
		{
			name:         "multiple semver - same major",
			in:           []string{"foo@v1.0.0", "foo@v1.3.7", "foo@v1.2.0", "foo@v1.0.1"},
			wantPicked:   []string{"foo@v1.3.7"},
			wantUnpicked: []string{"foo@v1.0.0", "foo@v1.2.0", "foo@v1.0.1"},
		},
		{
			name:         "multiple semver - multiple major",
			in:           []string{"foo@v1.0.0", "foo@v1.3.7", "foo/v2@v2.2.0", "foo/v2@v2.0.1", "foo@v1.1.1"},
			wantPicked:   []string{"foo/v2@v2.2.0", "foo@v1.3.7"},
			wantUnpicked: []string{"foo@v1.0.0", "foo/v2@v2.0.1", "foo@v1.1.1"},
		},
		{
			name:         "semver and pseudo version",
			in:           []string{"foo@v1.0.0", "foo@v1.3.7", "foo/v2@v2.2.0", "foo/v2@v2.0.1", "foo@v1.1.1", "foo@v0.0.0-20190311183353-d8887717615a"},
			wantPicked:   []string{"foo/v2@v2.2.0", "foo@v1.3.7"},
			wantUnpicked: []string{"foo@v1.0.0", "foo/v2@v2.0.1", "foo@v1.1.1", "foo@v0.0.0-20190311183353-d8887717615a"},
		},
		{
			name: "multiple pseudo version",
			in: []string{
				"foo@v0.0.0-20190311183353-d8887717615a",
				"foo@v0.0.0-20190227222117-0694c2d4d067",
				"foo@v0.0.0-20190312151545-0bb0c0a6e846",
			},
			wantPicked: []string{"foo@v0.0.0-20190312151545-0bb0c0a6e846"},
			wantUnpicked: []string{
				"foo@v0.0.0-20190227222117-0694c2d4d067",
				"foo@v0.0.0-20190311183353-d8887717615a",
			},
		},
		{
			name:         "semver and suffix",
			in:           []string{"foo@v1.0.0", "foo@v1.3.8-rc1", "foo@v1.3.7"},
			wantPicked:   []string{"foo@v1.3.8-rc1"},
			wantUnpicked: []string{"foo@v1.0.0", "foo@v1.3.7"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			for _, node := range tc.in {
				fmt.Fprintf(&buf, "A %s\n", node)
			}

			g, err := convert(&buf)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(g.mvsPicked, tc.wantPicked) {
				t.Fatalf("picked: got %v, want %v", g.mvsPicked, tc.wantPicked)
			}
			if !reflect.DeepEqual(g.mvsUnpicked, tc.wantUnpicked) {
				t.Fatalf("unpicked: got %v, want %v", g.mvsUnpicked, tc.wantUnpicked)
			}
		})
	}
}
