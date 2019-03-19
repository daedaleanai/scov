package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TempDirectory(t *testing.T) (filename string, cleanup func()) {
	name, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatalf("could not create temporary file: %s", err)
		// unreachable
	}

	return name, func() {
		err := os.RemoveAll(name)
		if err != nil {
			t.Logf("could not remove temporary directory: %s", err)
		}
	}
}

func TestCreateHTML(t *testing.T) {
	data := make(map[string]*FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	data = map[string]*FileData{
		"binc.cpp": data["binc.cpp"],
	}

	name, cleanup := TempDirectory(t)
	defer cleanup()

	*srcdir = "./testdata"
	err = createHTML(name, data, time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
	if err != nil {
		t.Fatalf("could not write output: %s", err)
	}

	expected, err := ioutil.ReadFile("./testdata/TestWriteHTMLIndex.golden")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}
	if out, err := ioutil.ReadFile(filepath.Join(name, "index.html")); err != nil {
		t.Fatalf("could not read output file: %s", err)
	} else if string(expected) != string(out) {
		LogNE(t, "output text", string(expected), string(out))
	}

	expected, err = ioutil.ReadFile("./testdata/TestWriteHTMLForSource.golden")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}
	if out, err := ioutil.ReadFile(filepath.Join(name, "binc.cpp.html")); err != nil {
		t.Fatalf("could not read output file: %s", err)
	} else if string(expected) != string(out) {
		LogNE(t, "output text", string(expected), string(out))
	}
}

func TestCreateHTMLIndex(t *testing.T) {
	data := make(map[string]*FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	data = map[string]*FileData{
		"binc.cpp": data["binc.cpp"],
	}

	expected, err := ioutil.ReadFile("./testdata/TestWriteHTMLIndex.golden")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	name, cleanup := TempFilename(t)
	defer cleanup()

	*srcdir = "./testdata"
	err = createHTMLIndex(name, data, time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
	if err != nil {
		t.Fatalf("could not write output: %s", err)
	}

	if out, err := ioutil.ReadFile(name); err != nil {
		t.Fatalf("could not read output file: %s", err)
	} else if string(expected) != string(out) {
		LogNE(t, "output text", string(expected), string(out))
	}
}

func TestWriteHTMLIndex(t *testing.T) {
	data := make(map[string]*FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	data = map[string]*FileData{
		"binc.cpp": data["binc.cpp"],
	}

	expected, err := ioutil.ReadFile("./testdata/TestWriteHTMLIndex.golden")
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
	data := make(map[string]*FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	expected, err := ioutil.ReadFile("./testdata/TestWriteHTMLForSource.golden")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	name, cleanup := TempFilename(t)
	defer cleanup()

	*srcdir = "./testdata"
	err = createHTMLForSource(name, "binc.cpp", data["binc.cpp"], time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
	if err != nil {
		t.Fatalf("could not write output: %s", err)
	}

	if out, err := ioutil.ReadFile(name); err != nil {
		t.Fatalf("could not read output file: %s", err)
	} else if string(expected) != string(out) {
		LogNE(t, "output text", string(expected), string(out))
	}
}

func TestWriteHTMLForSource(t *testing.T) {
	data := make(map[string]*FileData)
	err := loadFile(data, "./testdata/binc-7.3.0.cpp.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	expected, err := ioutil.ReadFile("./testdata/TestWriteHTMLForSource.golden")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	buffer := &strings.Builder{}
	*srcdir = "./testdata"
	err = writeHTMLForSource(buffer, "binc.cpp", data["binc.cpp"], time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
	if err != nil {
		t.Errorf("could not write output: %s", err)
	}

	if out := buffer.String(); len(expected) != len(out) {
		LogNE(t, "length of output", len(expected), len(out))
	} else if string(expected) != out {
		LogNE(t, "output", string(expected), out)
	}
}

func TestWriteBranchDescription(t *testing.T) {
	cases := []struct {
		data     []BranchStatus
		withData bool
		out      string
	}{
		{nil, false, ""},
		{nil, true, `<td></td>`},
		{[]BranchStatus{}, true, `<td></td>`},
		{[]BranchStatus{BranchNotExec}, true, `<td>[ NE ]</td>`},
		{[]BranchStatus{BranchTaken, BranchNotTaken}, true, `<td>[ + - ]</td>`},
	}

	for i, v := range cases {
		s := bytes.NewBuffer(nil)
		w := bufio.NewWriter(s)

		writeBranchDescription(w, v.withData, v.data)
		err := w.Flush()
		if err != nil {
			t.Errorf("failed to write, %s", err)
		}
		if s.String() != v.out {
			t.Errorf("Case %d: expected %s, got %s", i, v.out, s.String())
		}
	}
}
