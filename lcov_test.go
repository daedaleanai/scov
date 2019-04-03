package main

import (
	"path/filepath"
	"testing"
)

func TestLoadLCovFile(t *testing.T) {
	cases := []struct {
		filename string
		lcov     Coverage
		fcov     Coverage
		bcov     Coverage
	}{
		{"example-lcov-1.13.info", Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}},
	}
	for _, v := range cases {
		t.Run(v.filename, func(t *testing.T) {
			// Reset the nasty globals
			data := make(map[string]*FileData)
			*srcdir = "./example"

			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}
			fileData := data["example.c"]
			if lcov := fileData.LineCoverage(); lcov != v.lcov {
				LogNE(t, "line coverage", v.lcov, lcov)
			}
			if fcov := fileData.FuncCoverage(); fcov != v.fcov {
				LogNE(t, "function coverage", v.fcov, fcov)
			}
			if bcov := fileData.BranchCoverage(); bcov != v.bcov {
				LogNE(t, "branch coverage", v.bcov, bcov)
			}
		})
	}
}
