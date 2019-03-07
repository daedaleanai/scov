package main

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestWriteHTMLIndex(t *testing.T) {
	data := make(map[string]FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	data = map[string]FileData{
		"binc.cpp": data["binc.cpp"],
	}

	expected, err := ioutil.ReadFile("./testdata/binc-7.3.0.cpp.gcov.index.html")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	buffer := &strings.Builder{}
	err = writeHTMLIndex(buffer, data, time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
	if err != nil {
		t.Errorf("could not write output: %s", err)
	}

	if out := buffer.String(); len(expected) != len(out) {
		LogNE(t, "length of output", len(expected), len(out))
	} else if string(expected) != out {
		LogNE(t, "output", string(expected), out)
	}

}

func TestCreateHTMLForSource(t *testing.T) {
	data := make(map[string]FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	expected, err := ioutil.ReadFile("./testdata/binc-7.3.0.cpp.gcov.binc.cpp.html")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	buffer := &strings.Builder{}
	*srcdir = "./testdata"
	err = writeHTMLForSource(buffer, "binc.cpp", data["binc.cpp"])
	if err != nil {
		t.Errorf("could not write output: %s", err)
	}

	if out := buffer.String(); len(expected) != len(out) {
		LogNE(t, "length of output", len(expected), len(out))
		ioutil.WriteFile("./tmp.html", []byte(out), 0600)
	} else if string(expected) != out {
		LogNE(t, "output", string(expected), out)
	}

}
