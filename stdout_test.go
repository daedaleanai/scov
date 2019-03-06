package main

import (
	"strings"
	"testing"
)

func TestWriteTextReport(t *testing.T) {
	const expected = ` 59.6%	100.0%	binc.cpp
------	------	
 59.6%	100.0%	Overall
`

	data := make(map[string]FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	data = map[string]FileData{
		"binc.cpp": data["binc.cpp"],
	}

	buffer := &strings.Builder{}
	err = writeTextReport(buffer, data)
	if err != nil {
		t.Errorf("could not write output: %s", err)
	}

	if out := buffer.String(); expected != out {
		LogNE(t, "output text", expected, out)
	}
}
