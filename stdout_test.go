package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestWriteStdoutReport(t *testing.T) {
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

			report := NewTestReport()
			report.CollectStatistics(data)

			buffer := bytes.NewBuffer(nil)
			writeStdoutReport(buffer, report)

			if *update {
				err := ioutil.WriteFile(filepath.Join("./testdata", t.Name()+".golden"), buffer.Bytes(), 0600)
				if err != nil {
					t.Fatalf("could not write golden file: %s", err)
				}
			}

			expected, err := ioutil.ReadFile(filepath.Join("./testdata", t.Name()+".golden"))
			if err != nil {
				t.Fatalf("could not read golden file: %s", err)
			}
			if !bytes.Equal(expected, buffer.Bytes()) {
				t.Errorf("output does not match golden file")
			}
		})
	}
}
