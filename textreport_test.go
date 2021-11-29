package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

func TempFilename(t *testing.T) (filename string, closer func()) {
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

			report := NewTestReport()
			report.CollectStatistics(data)

			err = createTextReport(filename, report)
			if err != nil {
				t.Fatalf("could not write output: %s", err)
			}
			out, err := ioutil.ReadFile(filename)
			if err != nil {
				t.Fatalf("could not read the output: %s", err)
			}

			if *update {
				err := ioutil.WriteFile(filepath.Join("./testdata", t.Name()+".golden"), out, 0600)
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
	report := NewTestReport()
	report.CollectStatistics(map[string]*FileData{})

	err := createTextReport(".", report)
	if err == nil {
		t.Errorf("unexpected success")
	}
}
