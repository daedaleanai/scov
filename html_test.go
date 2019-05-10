package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
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
	err := loadFile(data, "./testdata/example-7.4.0.c.gcov")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	data = map[string]*FileData{
		"example.c": data["example.c"],
	}

	name, cleanup := TempDirectory(t)
	defer cleanup()

	*srcdir = "./example"
	err = createHTML(name, data, time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
	if err != nil {
		t.Fatalf("could not write output: %s", err)
	}

	expected, err := ioutil.ReadFile("./testdata/TestCreateHTMLIndex/example-7.4.0.c.gcov.golden")
	if err != nil {
		t.Fatalf("could not read file: %s", err)
	}

	if out, err := ioutil.ReadFile(filepath.Join(name, "index.html")); err != nil {
		t.Fatalf("could not read output file: %s", err)
	} else if !bytes.Equal(expected, out) {
		LogNE(t, "output text", string(expected), string(out))
	}
}

func TestCreateHTMLIndex(t *testing.T) {
	cases := []struct {
		filename string
	}{
		{"example-7.4.0.c.gcov"},
		{"example-7.4.0-branches.c.gcov"},
		{"example-8.3.0.c.gcov"},
		{"example-8.3.0-branches.c.gcov"},
	}

	for _, v := range cases {
		v := v
		t.Run(v.filename, func(t *testing.T) {
			data := make(map[string]*FileData)
			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}

			data = map[string]*FileData{
				"example.c": data["example.c"],
			}

			filename, cleanup := TempFilename(t)
			defer cleanup()

			*srcdir = "./example"
			err = createHTMLIndex(filename, data, time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
			if err != nil {
				t.Fatalf("could not write output: %s", err)
			}
			out, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("could not read the output: %s", err)
			}

			if *update {
				err := ioutil.WriteFile(filepath.Join("./testdata", t.Name()+".golden"), out, 0644)
				if err != nil {
					t.Fatalf("could not write golden file: %s", err)
				}
			}

			expected, err := ioutil.ReadFile(filepath.Join("./testdata", t.Name()+".golden"))
			if err != nil {
				t.Fatalf("could not read golden file: %s", err)
			}
			if !bytes.Equal(expected, out) {
				t.Errorf("output does not match golden file")
			}
		})
	}
}

func TestCreateJS(t *testing.T) {
	name, cleanup := TempFilename(t)
	defer cleanup()

	err := createJS(name)
	if err != nil {
		t.Errorf("could not write output: %s", err)
	}
}

func TestCreateHTMLForSource(t *testing.T) {
	cases := []struct {
		filename string
	}{
		{"example-7.4.0.c.gcov"},
		{"example-7.4.0-branches.c.gcov"},
		{"example-8.3.0.c.gcov"},
		{"example-8.3.0-branches.c.gcov"},
	}

	for _, v := range cases {
		v := v
		t.Run(v.filename, func(t *testing.T) {
			data := make(map[string]*FileData)
			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}

			filename, cleanup := TempFilename(t)
			defer cleanup()

			*srcdir = "./example"
			err = createHTMLForSource(filename, "example.c", data["example.c"], time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC))
			if err != nil {
				t.Fatalf("could not write output: %s", err)
			}
			out, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("could not read the output: %s", err)
			}

			if *update {
				err := ioutil.WriteFile(filepath.Join("./testdata", t.Name()+".golden"), out, 0644)
				if err != nil {
					t.Fatalf("could not write golden file: %s", err)
				}
			}

			expected, err := ioutil.ReadFile(filepath.Join("./testdata", t.Name()+".golden"))
			if err != nil {
				t.Fatalf("could not read golden file: %s", err)
			}
			if !bytes.Equal(expected, out) {
				t.Errorf("output does not match golden file")
			}
		})
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
