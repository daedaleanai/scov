package main

import (
	"path/filepath"
	"testing"
)

func TestIdentifyFileType(t *testing.T) {
	cases := []struct {
		filename string
		parser   Parser
		ok       bool
	}{
		{"example-7.4.0.c.gcov", ParserGCov, true},
		{"/home/person/example-7.4.0.c.gcov", ParserGCov, true},
		{"example-7.4.0.c.info", ParserLCov, true},
		{"/home/person/example-7.4.0.c.info", ParserLCov, true},
		{"example-7.4.0.json", ParserLLVM, true},
		{"/home/person/example-7.4.0.json", ParserLLVM, true},
		{"example-7.4.0.out", ParserGo, true},
		{"/home/person/example-7.4.0.out", ParserGo, true},
		{"example.7.4.0.c.dummy", 0, false},
		{"/home/person/example.7.4.0.c.dummy", 0, false},
	}

	for _, v := range cases {
		v := v
		t.Run(v.filename, func(t *testing.T) {
			parser, ok := identifyFileType(v.filename)

			if v.parser != parser {
				t.Errorf("parser:  wanted %v, got %v", v.parser, parser)
			}
			if v.ok != ok {
				t.Errorf("parser:  wanted %v, got %v", v.ok, ok)
			}
		})
	}
}

func TestParserLoadFile(t *testing.T) {
	cases := []struct {
		filename  string
		lcovTotal Coverage
		lcov      Coverage
		fcov      Coverage
		bcov      Coverage
		rcov      Coverage
		fileCount int
		funcStart int
	}{
		// gcc 7.4.0
		{"example-7.4.0.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{}, Coverage{}, 1, 28},
		{"example-7.4.0-branches.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, Coverage{}, 1, 28},
		{"example-7.4.0-branches", Coverage{18, 22}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, Coverage{}, 3, 28},
		// gcc 8.3.0
		{"example-8.3.0.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{}, Coverage{}, 1, 28},
		{"example-8.3.0-branches.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, Coverage{}, 1, 28},
		{"example-8.3.0-branches", Coverage{18, 22}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, Coverage{}, 3, 28},
		// gcc 9.1.0
		{"example-9.1.0.c.gcov.json.gz", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, Coverage{}, 1, 28},
		// gcc with lcov
		{"example-lcov-1.13.info", Coverage{18, 22}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, Coverage{}, 3, 28},
		// clang 6.0.1
		{"example-llvm-6.0.1.json", Coverage{37, 47}, Coverage{14, 17}, Coverage{1, 1}, Coverage{}, Coverage{5, 6}, 3, 29},
		// clang 8.0.1
		{"example-llvm-8.0.1.info", Coverage{57, 67}, Coverage{28, 31}, Coverage{3, 3}, Coverage{}, Coverage{}, 3, 29},
		{"example2-llvm-8.0.1.info", Coverage{57, 67}, Coverage{28, 31}, Coverage{3, 3}, Coverage{}, Coverage{}, 3, 29},
		{"example-llvm-8.0.1.json", Coverage{33, 47}, Coverage{10, 17}, Coverage{1, 1}, Coverage{}, Coverage{4, 6}, 3, 29},
	}
	for _, v := range cases {
		t.Run(v.filename, func(t *testing.T) {
			data := make(FileDataSet)

			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}
			data.ConvertRegionToLineData()
			if lcov := data.LineCoverage(); lcov != v.lcovTotal {
				t.Errorf("total line coverage: expected %v, got %v", v.lcovTotal, lcov)
			}
			if got := len(data); got != v.fileCount {
				t.Errorf("total file count: expected %v, got %v", v.fileCount, got)
			}
			data, err = normalizeSourceFilenames(data, "/example/")
			if err != nil {
				t.Fatalf("could not normalize filenames: %s", err)
			}
			if lcov := data.LineCoverage(); lcov != v.lcovTotal {
				t.Errorf("total line coverage: expected %v, got %v", v.lcovTotal, lcov)
			}
			if got := len(data); got != v.fileCount {
				t.Errorf("total file count: expected %v, got %v", v.fileCount, got)
			}

			fileData, ok := data["example.c"]
			if !ok {
				t.Fatalf("missing data for file 'example.c' after loading")
			}
			if lcov := fileData.LineCoverage(); lcov != v.lcov {
				LogNE(t, "line coverage", v.lcov, lcov)
			}
			if fcov := fileData.FuncCoverage(); fcov != v.fcov {
				LogNE(t, "function coverage", v.fcov, fcov)
			}
			if got := fileData.FuncData["main"].StartLine; got != v.funcStart {
				LogNE(t, "location of main", v.funcStart, got)
			}
			if bcov := fileData.BranchCoverage(); bcov != v.bcov {
				LogNE(t, "branch coverage", v.bcov, bcov)
			}
			if rcov := fileData.RegionCoverage(); rcov != v.rcov {
				LogNE(t, "region coverage", v.rcov, rcov)
			}
		})
	}

}
