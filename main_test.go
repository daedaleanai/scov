package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"
)

func LogNE(t *testing.T, field string, expected, got interface{}) {
	t.Errorf("Values for %s not equal, expected %v, got %v", field, expected, got)
}

func TestHandleRequestFlags(t *testing.T) {
	cases := []struct {
		help    bool
		version bool
		ok      bool
	}{
		{false, false, false},
		{true, false, true},
		{false, true, true},
	}

	for i, v := range cases {
		buffer := bytes.NewBuffer(nil)
		ok := handleRequestFlags(buffer, v.help, v.version)

		if ok != (buffer.Len() > 0) {
			t.Errorf("Case %d: mismatch between output and return status", i)
		}
		if ok != v.ok {
			t.Errorf("Case %d: expected %v, got %v", i, v.ok, ok)
		}
	}
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

func TestFilterExternalFileData(t *testing.T) {
	const file1 = "binc.cpp"
	const file2 = "/usr/include/a.h"
	const file3 = "/usr/include/b.h"

	cases := []struct {
		in       []string
		external bool
		expected int
	}{
		{[]string{file1, file2, file3}, false, 1},
		{[]string{file1, file2, file3}, true, 3},
		{[]string{file1}, false, 1},
		{[]string{file1}, true, 1},
	}

	for i, v := range cases {
		name := fmt.Sprintf("Case %d", i)
		t.Run(name, func(t *testing.T) {
			data := make(map[string]*FileData)
			for _, v := range v.in {
				data[v] = NewFileData(v)
			}

			data = filterExternalFileData(data, v.external)
			if out := len(data); out != v.expected {
				LogNE(t, "file count", v.expected, out)
			}
		})
	}
}

func TestFilterExcludedFileData(t *testing.T) {
	const file1 = "binc.cpp"
	const file2 = "/usr/include/a.h"
	const file3 = "/usr/include/b.h"

	cases := []struct {
		in       []string
		filter   string
		ok       bool
		expected int
	}{
		{[]string{file1, file2, file3}, "", true, 3},
		{[]string{file1, file2, file3}, ".+", true, 0},
		{[]string{file1, file2, file3}, "^/", true, 1},
		{[]string{file1, file2, file3}, "\\.h$", true, 1},
		{[]string{file1, file2, file3}, "\\.cpp", true, 2},
		{[]string{file1, file2, file3}, "[]", false, 3},
	}

	for i, v := range cases {
		name := fmt.Sprintf("Case %d", i)
		t.Run(name, func(t *testing.T) {
			data := make(map[string]*FileData)
			for _, v := range v.in {
				data[v] = NewFileData(v)
			}

			out := bytes.NewBuffer(nil)
			data = filterExcludedFileData(out, data, v.filter)
			if ok := out.Len() == 0; ok != v.ok {
				LogNE(t, "ok", v.ok, ok)
			}
			if out := len(data); out != v.expected {
				LogNE(t, "file count", v.expected, out)
			}
		})
	}
}
