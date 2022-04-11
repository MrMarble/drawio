package drawio

import (
	"encoding/xml"
	"fmt"
)

type Context struct {
	id         string
	shapeCount int
	width      int
	height     int
	pageName   string
	diagram    *mxFile
}

type DiagramOptions struct {
	GridSize        int
	ShowGrid        bool
	PageView        bool
	ShowBackground  bool
	Shadow          bool
	GridColor       string
	BackgroundColor string

	ConnectionArrows bool
	ConnectionPoints bool
	Guides           bool
}

type StyleOptions struct{}

type Option interface {
	Apply(ctx *Context)
}

// OptionFunc is function that adheres to the Option interface.
type OptionFunc func(ctx *Context)

func (o OptionFunc) Apply(ctx *Context) { o(ctx) } // nolint: revive

func NewContext(width, height int, name string, options ...Option) *Context {
	ctx := Context{width: width, height: height, pageName: name, diagram: newDiagram(), id: randString(20)}
	ctx.diagram.Diagram.Name = name
	for _, option := range options {
		option.Apply(&ctx)
	}
	return &ctx
}

func NewContextForDiagram(diagram []byte) (*Context, error) {
	mx, err := readDiagram(diagram)
	if err != nil {
		return nil, err
	}

	ctx := Context{
		width:    mx.Diagram.Model.PageWidth,
		height:   mx.Diagram.Model.PageHeight,
		pageName: mx.Diagram.Name,
		diagram:  mx,
	}

	return &ctx, nil
}

func ConfigureDiagram(d DiagramOptions) Option {
	return OptionFunc(func(ctx *Context) {
		ctx.diagram.Diagram.Model.GridSize = d.GridSize

		*ctx.diagram.Diagram.Model.Grid = boolean(d.ShowGrid)
		*ctx.diagram.Diagram.Model.Arrows = boolean(d.ConnectionArrows)
		*ctx.diagram.Diagram.Model.Connect = boolean(d.ConnectionPoints)
		*ctx.diagram.Diagram.Model.Page = boolean(d.PageView)
		*ctx.diagram.Diagram.Model.Guides = boolean(d.Guides)
		*ctx.diagram.Diagram.Model.Shadow = boolean(d.Shadow)
	})
}

type EncodeOptions struct {
	Raw    bool
	Indent string
}

func (ctx *Context) Encode(options EncodeOptions) ([]byte, error) {
	var marshalled []byte
	var err error

	if options.Raw {
		ctx.diagram.Diagram.raw = true
	}

	if options.Indent != "" {
		marshalled, err = xml.MarshalIndent(ctx.diagram, "", options.Indent)
	} else {
		marshalled, err = xml.Marshal(ctx.diagram)
	}

	ctx.diagram.Diagram.raw = false
	return marshalled, err
}

func (ctx *Context) Circle(x, y, size int) string {
	id := fmt.Sprintf("%s-%d", ctx.id, ctx.shapeCount)
	ctx.shapeCount++
	ctx.diagram.Diagram.Model.Cells = append(ctx.diagram.Diagram.Model.Cells, mxCell{
		ID:     id,
		Value:  "",
		Vertex: 1,
		Parent: "1",
		Style:  styleMap{"ellipse": "", "aspect": "fixed", "html": "1"},
		Geometry: &mxGeometry{
			As:     "geometry",
			X:      x,
			Y:      y,
			Width:  size,
			Height: size,
		},
	})
	return id
}

func (ctx *Context) Square(x, y, width, height int) string {
	id := fmt.Sprintf("%s-%d", ctx.id, ctx.shapeCount)
	ctx.shapeCount++
	ctx.diagram.Diagram.Model.Cells = append(ctx.diagram.Diagram.Model.Cells, mxCell{
		ID:     id,
		Value:  "",
		Vertex: 1,
		Parent: "1",
		Style:  styleMap{"html": "1", "rounded": "0"},
		Geometry: &mxGeometry{
			As:     "geometry",
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
	})
	return id
}
