package drawio

import (
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestNewContext(t *testing.T) {
	tests := map[string]struct {
		testFile string
		ctx      *Context
	}{
		"Context without options": {testFile: "TestNewContext", ctx: NewContext(10, 10, "Test")},
		"Context with options": {testFile: "TestNewContextWithOptions", ctx: NewContext(10, 10, "Test", ConfigureDiagram(DiagramOptions{
			Shadow:   true,
			ShowGrid: true,
			GridSize: 25,
		}))},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			marshalled, err := test.ctx.Encode(EncodeOptions{Raw: true, Indent: "    "})
			if err != nil {
				t.Fatal(err)
			}

			g := goldie.New(t)

			template := struct {
				Modified string
				Etag     string
				Name     string
				ID       string
			}{
				Modified: test.ctx.diagram.Modified,
				Etag:     test.ctx.diagram.Etag,
				Name:     test.ctx.pageName,
				ID:       test.ctx.diagram.Diagram.ID,
			}
			g.AssertWithTemplate(t, test.testFile, template, marshalled)
		})
	}
}

func TestCircle(t *testing.T) {
	ctx := NewContext(10, 10, "Test")

	id := ctx.Circle(10, 10, 80)

	marshalled, err := ctx.Encode(EncodeOptions{Raw: true})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(marshalled), id) {
		t.Fatal()
	}
}
