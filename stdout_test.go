package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

func TestCreateTextReport(t *testing.T) {
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

			err = createTextReport(filename, data)
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

func TestCreateTextReportFail(t *testing.T) {
	data := map[string]*FileData{}

	err := createTextReport(".", data)
	if err == nil {
		t.Errorf("unexpected success")
	}
}
