package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TempFilename(t *testing.T) (filename string, cleanup func()) {
	file, err := ioutil.TempFile("", "testing")
	if err != nil {
		t.Fatalf("could not create temporary file: %s", err)
		// unreachable
	}
	name := file.Name()
	file.Close()

	return name, func() {
		err := os.Remove(name)
		if err != nil {
			t.Logf("could not remove temporary file: %s", err)
		}
	}
}

func TestCreateTextReport(t *testing.T) {
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

	filename, cleanup := TempFilename(t)
	defer cleanup()

	err = createTextReport(filename, data)
	if err != nil {
		t.Fatalf("could not write output: %s", err)
	}

	if out, err := ioutil.ReadFile(filename); err != nil {
		t.Fatalf("could not read output file: %s", err)
	} else if expected != string(out) {
		LogNE(t, "output text", expected, string(out))
	}
}

func TestCreateTextReportFail(t *testing.T) {
	data := map[string]FileData{}

	err := createTextReport(".", data)
	if err == nil {
		t.Errorf("unexpected success")
	}
}

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
