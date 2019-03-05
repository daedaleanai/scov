package main

import (
	"strings"
	"testing"
)

const buildTextExpected = ` 59.6%	100.0%	binc.cpp
------	------	
 59.6%	100.0%	Overall
`

func TestBuildText(t *testing.T) {
	data := make(map[string]FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	data = map[string]FileData{
		"binc.cpp": data["binc.cpp"],
	}

	buffer := &strings.Builder{}
	err = buildText(buffer, data)
	if err != nil {
		t.Errorf("could not write output: %s", err)
	}

	if out := buffer.String(); buildTextExpected != out {
		LogNE(t, "output text", buildTextExpected, out)
	}
}
