package drawio

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestUnmarshal(t *testing.T) {
	file, err := os.ReadFile("testdata/TestUnmarshal.drawio")
	if err != nil {
		t.Fatal(err)
	}

	var diagram mxFile
	err = xml.Unmarshal(file, &diagram)
	if err != nil {
		t.Fatal(err)
	}

	g := goldie.New(t)
	diagram.Diagram.raw = true
	g.AssertXml(t, "TestUnmarshal", &diagram)
}

func TestMarshal(t *testing.T) {
	file, err := os.ReadFile("testdata/TestUnmarshal.drawio")
	if err != nil {
		t.Fatal(err)
	}

	var diagram mxFile
	err = xml.Unmarshal(file, &diagram)
	if err != nil {
		t.Fatal(err)
	}

	g := goldie.New(t)
	g.AssertXml(t, "TestMarshal", &diagram)
}
