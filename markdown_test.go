package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCreateMarkdownReport(t *testing.T) {
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
			err = createMarkdownReport(filename, report)
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

func TestCreateMarkdownReportFail(t *testing.T) {
	report := NewTestReport()
	report.CollectStatistics(map[string]*FileData{})

	err := createMarkdownReport(".", report)
	if err == nil {
		t.Errorf("unexpected success")
	}
}
