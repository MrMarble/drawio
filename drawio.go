/*
	Create and edit draw.io diagrams
*/
package drawio

import (
	"encoding/xml"
	"strings"
	"time"
)

type boolean bool

type styleMap map[string]string

type mxFile struct {
	XMLName  xml.Name `xml:"mxfile"`
	Host     string   `xml:"host,attr"`
	Modified string   `xml:"modified,attr"`
	Agent    string   `xml:"agent,attr"`
	Etag     string   `xml:"etag,attr"`
	Version  string   `xml:"version,attr"`
	Type     string   `xml:"type,attr"`
	Diagram  diagram  `xml:"diagram"`
}

type diagram struct {
	raw     bool
	XMLName xml.Name     `xml:"diagram"`
	ID      string       `xml:"id,attr"`
	Name    string       `xml:"name,attr"`
	Model   mxGraphModel `xml:",chardata"`
}

type mxGraphModel struct {
	XMLName    xml.Name `xml:"mxGraphModel"`
	Dx         int      `xml:"dx,attr"`
	Dy         int      `xml:"dy,attr"`
	Grid       *boolean `xml:"grid,attr"`
	GridSize   int      `xml:"gridSize,attr"`
	Guides     *boolean `xml:"guides,attr"`
	Tooltips   *boolean `xml:"tooltips,attr"`
	Connect    *boolean `xml:"connect,attr"`
	Arrows     *boolean `xml:"arrows,attr"`
	Fold       *boolean `xml:"fold,attr"`
	Page       *boolean `xml:"page,attr"`
	PageScale  int      `xml:"pageScale,attr"`
	PageWidth  int      `xml:"pageWidth,attr"`
	PageHeight int      `xml:"pageHeight,attr"`
	Math       int      `xml:"math,attr"`
	Shadow     *boolean `xml:"shadow,attr"`
	Cells      []mxCell `xml:"root>mxCell"`
}

type mxCell struct {
	XMLName  xml.Name    `xml:"mxCell"`
	ID       string      `xml:"id,attr"`
	Parent   string      `xml:"parent,attr,omitempty"`
	Value    string      `xml:"value,attr,omitempty"`
	Style    styleMap    `xml:"style,attr,omitempty"`
	Source   string      `xml:"source,attr,omitempty"`
	Target   string      `xml:"target,attr,omitempty"`
	Vertex   int         `xml:"vertex,attr,omitempty"`
	Edge     int         `xml:"edge,attr,omitempty"`
	Geometry *mxGeometry `xml:"mxGeometry"`
}

type mxGeometry struct {
	XMLName  xml.Name `xml:"mxGeometry"`
	As       string   `xml:"as,attr"`
	Relative *boolean `xml:"relative,attr,omitempty"`
	X        int      `xml:"x,attr,omitempty"`
	Y        int      `xml:"y,attr,omitempty"`
	Width    int      `xml:"width,attr,omitempty"`
	Height   int      `xml:"height,attr,omitempty"`
}

func (d *diagram) UnmarshalXML(decoder *xml.Decoder, s xml.StartElement) error {
	for _, a := range s.Attr {
		if a.Name.Local == "id" {
			d.ID = a.Value
		} else if a.Name.Local == "name" {
			d.Name = a.Value
		}
	}

	// Little hack to extract CDATA
	e := struct {
		XMLName xml.Name
		Content string `xml:",chardata"`
		// FullContent   string `xml:",innerxml"` // for debug purpose, allow to see what's inside some tags
	}{}

	err := decoder.DecodeElement(&e, &s)
	if err != nil {
		return err
	}
	decoded := decode(e.Content)

	var m mxGraphModel
	err = xml.Unmarshal([]byte(decoded), &m)
	if err != nil {
		return err
	}
	d.Model = m
	return nil
}

func (d *diagram) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{{Name: xml.Name{Local: "id"}, Value: d.ID}, {Name: xml.Name{Local: "name"}, Value: d.Name}}
	marshalled, err := xml.Marshal(d.Model)
	if err != nil {
		return err
	}
	encoder.EncodeToken(start)

	if d.raw {
		encoder.EncodeToken(encoder.EncodeElement(d.Model, xml.StartElement{Name: xml.Name{Local: "mxGraphModel"}}))
	} else {
		encoded := encode(string(marshalled))
		encoder.EncodeToken(xml.CharData(encoded))
	}

	return encoder.EncodeToken(xml.EndElement{Name: start.Name})
}

func (b *boolean) UnmarshalXMLAttr(attr xml.Attr) error {
	*b = attr.Value == "1"
	return nil
}

func (b *boolean) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if b != nil && *b {
		return xml.Attr{Name: name, Value: "1"}, nil
	} else {
		return xml.Attr{Name: name, Value: "0"}, nil
	}
}

func (s styleMap) UnmarshalXMLAttr(attr xml.Attr) error {
	if s == nil {
		s = make(styleMap)
	}
	styles := strings.Split(attr.Value, ";")
	for _, property := range styles {
		properties := strings.Split(property, "=")
		if len(properties) > 1 {
			s[properties[0]] = properties[1]
		} else {
			s[properties[0]] = ""
		}
	}
	return nil
}

func (s styleMap) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var css []string
	for key, value := range s {
		if value != "" {
			css = append(css, key+"="+value)
		} else {
			css = append(css, key)
		}
	}
	return xml.Attr{Name: name, Value: strings.Join(css, ";")}, nil
}

func readDiagram(diagram []byte) (*mxFile, error) {
	var mx mxFile
	err := xml.Unmarshal(diagram, &mx)
	if err != nil {
		return nil, err
	}
	return &mx, nil
}

func newDiagram() *mxFile {
	truthy := boolean(true)
	falsy := boolean(false)

	// default width A4
	width := 827
	height := 1169
	return &mxFile{
		Host:     "app.diagrams.net",
		Agent:    "5.0 (X11)",
		Version:  "17.4.2",
		Type:     "device",
		Etag:     randString(20),
		Modified: time.Now().Format("2006-02-01T15:04:05.999Z"),
		Diagram: diagram{
			ID: randString(20),
			Model: mxGraphModel{
				Grid:       &truthy,
				GridSize:   10,
				Page:       &truthy,
				Tooltips:   &truthy,
				Connect:    &truthy,
				Guides:     &truthy,
				Arrows:     &truthy,
				Fold:       &truthy,
				Shadow:     &falsy,
				PageScale:  1,
				PageWidth:  width,
				PageHeight: height,
				Cells:      []mxCell{{ID: "0"}, {ID: "1", Parent: "0"}},
			},
		},
	}
}
