package main

import (
	"path/filepath"
	"testing"
)

func TestParseFunctionRecord(t *testing.T) {
	cases := []struct {
		value    string
		name     string
		hitCount uint64
	}{
		{"function:222,6,main", "main", 6},
		{"function:222,223,6,main", "main", 6},
		{"function:222,0,main", "main", 0},
		{"function:222,223,0,main", "main", 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "function" {
				LogNE(t, "record type", "function", rt)
			}
			name, hc, err := parseFunctionRecord(value)
			if name != v.name {
				LogNE(t, "function name", v.name, name)
			}
			if hc != v.hitCount {
				LogNE(t, "hit count", v.hitCount, hc)
			}
			if err != nil {
				LogNE(t, "error", nil, err)
			}
		})
	}
}

func TestParseLCountRecord(t *testing.T) {
	cases := []struct {
		value    string
		lineNo   int
		hitCount uint64
	}{
		{"lcount:32,18", 32, 18},
		{"lcount:32,0", 32, 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "lcount" {
				LogNE(t, "record type", "lcount", rt)
			}
			lineNo, hc, err := parseLCountRecord(value)
			if lineNo != v.lineNo {
				LogNE(t, "function name", v.lineNo, lineNo)
			}
			if hc != v.hitCount {
				LogNE(t, "hit count", v.hitCount, hc)
			}
			if err != nil {
				LogNE(t, "error", nil, err)
			}
		})
	}
}

func TestLoadGCovFile(t *testing.T) {
	cases := []struct {
		filename string
		lcov     Coverage
		fcov     Coverage
		bcov     Coverage
	}{
		{"binc-7.3.0.cpp.gcov", Coverage{90, 151}, Coverage{12, 12}, Coverage{}},
		{"binc-8.2.0.cpp.gcov", Coverage{58, 119}, Coverage{10, 10}, Coverage{}},
		{"binc-8.2.0-branches.cpp.gcov", Coverage{59, 122}, Coverage{10, 10}, Coverage{17, 60}},
	}
	for _, v := range cases {
		t.Run(v.filename, func(t *testing.T) {
			// Reset the nasty globals
			data := make(map[string]*FileData)

			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}
			fileData := data["binc.cpp"]
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
