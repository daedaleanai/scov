package main

import (
	"path/filepath"
	"testing"
)

func TestNewReport(t *testing.T) {
	report := NewReport("deadBEAF")
	if report == nil {
		t.Fatalf("failed to get a new report")
	}
	if report.Title != "deadBEAF" {
		t.Errorf("report constructed with the wrong title")
	}
}

func TestReportCollectStatistics(t *testing.T) {
	cases := []struct {
		filename  string
		lcovTotal Coverage
		lcov      Coverage
		fcov      Coverage
		bcov      Coverage
		fileCount int
		funcCount int
	}{
		// gcc 7.4.0
		{"example-7.4.0.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{}, 1, 1},
		{"example-7.4.0-branches.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, 1, 1},
		{"example-7.4.0-branches", Coverage{18, 22}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, 3, 3},
		// gcc 8.3.0
		{"example-8.3.0.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{}, 1, 1},
		{"example-8.3.0-branches.c.gcov", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, 1, 1},
		{"example-8.3.0-branches", Coverage{18, 22}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, 3, 3},
		// gcc 9.1.0
		{"example-9.1.0.c.gcov.json.gz", Coverage{9, 10}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, 1, 1},
		// gcc with lcov
		{"example-lcov-1.13.info", Coverage{18, 22}, Coverage{9, 10}, Coverage{1, 1}, Coverage{2, 4}, 3, 3},
		// clang 6.0.1
		{"example-llvm-6.0.1.json", Coverage{37, 47}, Coverage{14, 17}, Coverage{1, 1}, Coverage{0, 0}, 3, 3},
		// clang 8.0.1
		{"example-llvm-8.0.1.info", Coverage{57, 67}, Coverage{28, 31}, Coverage{3, 3}, Coverage{0, 0}, 3, 9},
		{"example2-llvm-8.0.1.info", Coverage{57, 67}, Coverage{28, 31}, Coverage{3, 3}, Coverage{0, 0}, 3, 9},
		{"example-llvm-8.0.1.json", Coverage{33, 47}, Coverage{10, 17}, Coverage{1, 1}, Coverage{0, 0}, 3, 3},
	}
	for _, v := range cases {
		t.Run(v.filename, func(t *testing.T) {
			data := make(FileDataSet)

			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}

			report := NewTestReport()
			report.CollectStatistics(data)

			if got := len(report.Files); got != v.fileCount {
				t.Errorf("file count, want %v, got %v", v.fileCount, got)
			}
			if got := len(report.Funcs); got != v.funcCount {
				t.Errorf("func count, want %v, got %v", v.funcCount, got)
			}
		})
	}
}
