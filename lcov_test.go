package main

import (
	"path/filepath"
	"testing"
)

func LogNE(t *testing.T, field string, expected, got interface{}) {
	t.Errorf("Values for %s not equal, expected %v, got %v", field, expected, got)
}

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

func TestLoadFile(t *testing.T) {
	cases := []struct {
		filename string
		lcov1    Coverage
		fcov1    Coverage
	}{
		{"binc-7.3.0.cpp.gcov", Coverage{90, 151}, Coverage{12, 12}},
	}
	for _, v := range cases {
		t.Run(v.filename, func(t *testing.T) {
			// Reset the nasty globals
			data := make(map[string]FileData)

			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}
			fileData := data["binc.cpp"]
			if lcov := fileData.LineCoverage(); lcov != v.lcov1 {
				LogNE(t, "line coverage", v.lcov1, lcov)
			}
			if fcov := fileData.FuncCoverage(); fcov != v.fcov1 {
				LogNE(t, "function coverage", v.fcov1, fcov)
			}
		})
	}
}
