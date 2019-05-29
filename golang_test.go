package main

import (
	"path/filepath"
	"testing"
)

func TestParserLoadGoFile(t *testing.T) {
	cases := []struct {
		filename  string
		lcovTotal Coverage
		lcov      Coverage
	}{
		// go 1.10.4
		{"scov-1.10.4.out", Coverage{796, 932}, Coverage{28, 28}},
	}
	for _, v := range cases {
		t.Run(v.filename, func(t *testing.T) {
			data := make(FileDataSet)

			err := loadFile(data, filepath.Join("./testdata", v.filename))
			if err != nil {
				t.Fatalf("could not read file: %s", err)
			}
			if lcov := data.LineCoverage(); lcov != v.lcovTotal {
				t.Errorf("total line coverage: expected %v, got %v", v.lcovTotal, lcov)
			}

			fileData, ok := data["gitlab.com/stone.code/scov/gcovjs.go"]
			if !ok {
				t.Fatalf("missing data for file 'gcovjs.go' after loading")
			}
			if lcov := fileData.LineCoverage(); lcov != v.lcov {
				LogNE(t, "line coverage", v.lcov, lcov)
			}
		})
	}
}

func TestParseGoRecord(t *testing.T) {
	cases := []struct {
		in         string
		filename   string
		start      Position
		end        Position
		statements int
		hitCount   uint64
	}{
		{"gitlab.com/stone.code/scov/gcovjs.go:39.60,43.16 3 1", "gitlab.com/stone.code/scov/gcovjs.go", Position{39, 60}, Position{43, 16}, 3, 1},
		{"gitlab.com/stone.code/scov/gcovjs.go 39.60,43.16 3 1", "", Position{}, Position{}, 0, 0},
		{"gcovjs.go:39.60.43.16 3 1", "", Position{}, Position{}, 0, 0},
		{"gcovjs.go:39.60,43.16,3,1", "", Position{}, Position{}, 0, 0},
		{"gcovjs.go:39.60,43.16 3,1", "", Position{}, Position{}, 0, 0},
		{"gcovjs.go:39.,43.16 3 1", "", Position{}, Position{}, 0, 0},
		{"gcovjs.go:39.60,43. 3 1", "", Position{}, Position{}, 0, 0},
		{"gcovjs.go:39.60,43.16 # 1", "", Position{}, Position{}, 0, 0},
		{"gcovjs.go:39.60,43.16 3 #", "", Position{}, Position{}, 0, 0},
	}

	for _, v := range cases {
		t.Run(v.in, func(t *testing.T) {
			filename, start, end, statements, hitCount, err := parseGoRecord(v.in)
			if v.filename == "" {
				if err == nil {
					t.Errorf("missing error on malformed input")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if v.filename != filename {
					t.Errorf("expected %v, got %v", v.filename, filename)
				}
				if v.start != start {
					t.Errorf("expected %v, got %v", v.start, start)
				}
				if v.end != end {
					t.Errorf("expected %v, got %v", v.end, end)
				}
				if v.statements != statements {
					t.Errorf("expected %v, got %v", v.statements, statements)
				}
				if v.hitCount != hitCount {
					t.Errorf("expected %v, got %v", v.hitCount, hitCount)
				}
			}
		})
	}

}

func TestParseGoPosition(t *testing.T) {
	cases := []struct {
		in  string
		out Position
	}{
		{"12.34", Position{12, 34}},
		{"1234", Position{}},
		{"1a.34", Position{}},
		{"12.3a4", Position{}},
		{"12.34.", Position{}},
	}

	for _, v := range cases {
		t.Run(v.in, func(t *testing.T) {
			pos, err := parseGoPosition(v.in)
			if v.out.IsZero() {
				if err == nil {
					t.Errorf("missing error on malformed input")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if pos != v.out {
					t.Errorf("expected %v, got %v", v.out, pos)
				}
			}
		})
	}
}
