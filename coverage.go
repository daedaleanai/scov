package main

import (
	"fmt"
	"strconv"
)

// Coverage holds information about the coverage of some element over some scope.
// Typically elements are either individual lines, or functions.  Tyipcally, the
// scope is a source file or executable.
type Coverage struct {
	Hits  int // Count of lines or functions that were executed.
	Total int // Count of lines or functions in the scope.
}

// P returns the percentage of lines or functions that were executed.
func (c Coverage) P() float32 {
	return float32(c.Hits) * 100 / float32(c.Total)
}

// Q returns the percentage of lines or functions that were not executed.
func (c Coverage) Q() float32 {
	return 100 - float32(c.Hits)*100/float32(c.Total)
}

// Add combines the coverage data from different scopes.
func (c Coverage) Add(delta Coverage) Coverage {
	c.Hits += delta.Hits
	c.Total += delta.Total
	return c
}

// Rating returns the rating (low, medium, or high) for this coverage.
func (c Coverage) Rating() CoverageRating {
	if c.P() >= 90 {
		return HighCoverage
	}
	if c.P() >= 75 {
		return MediumCoverage
	}
	return LowCoverage
}

// String returns a human readable string representing the coverage.
func (c Coverage) String() string {
	return strconv.FormatInt(int64(c.Hits), 10) + "/" +
		strconv.FormatInt(int64(c.Total), 10)
}

// Format implements fmt.Formatter.
func (c Coverage) Format(f fmt.State, v rune) {
	if v == 'f' {
		width, _ := f.Width()
		prec, _ := f.Precision()

		if c.Total == 0 {
			fmt.Fprintf(f, "%*s", width, "--")
		} else {
			fmt.Fprintf(f, "%*.*f", width, prec, c.P())
		}
	} else if v == 's' {
		fmt.Fprintf(f, "%d/%d", c.Hits, c.Total)
	} else {
		panic("unsupported verb")
	}
}

// Valid returns true if data was collected.  In otherwords, unless the
// coverage is valid, any calculations may lead to divide by zero.
func (c Coverage) Valid() bool {
	return c.Total > 0
}

// CoverageRating is a classification of the coverage into low, medium or high.
type CoverageRating uint8

// These constants provide a rough classification for the amount of coverage in a scope.
const (
	LowCoverage CoverageRating = iota
	MediumCoverage
	HighCoverage
)

// String returns a string representation of the rating.
func (cr CoverageRating) String() string {
	if cr == LowCoverage {
		return "low"
	}
	if cr == MediumCoverage {
		return "medium"
	}
	return "high"
}

// FuncData represents data about a function.
type FuncData struct {
	StartLine int
	HitCount  uint64
}

// BranchStatus indicates whether a branch was taken, not taken, or if the
// conditional was never executed.
type BranchStatus uint8

// The following constants represent the possible status for branches.
const (
	BranchTaken BranchStatus = iota
	BranchNotTaken
	BranchNotExec
)

// Region defines a region of code in a source file.
type Region struct {
	StartLine int
	StartByte int
	EndLine   int
	EndByte   int
}

// FileData maintains coverage statistics for a single file.
type FileData struct {
	Filename   string
	LineData   map[int]uint64
	FuncData   map[string]FuncData
	BranchData map[int][]BranchStatus
	RegionData map[Region]uint64
}

// NewFileData initializes a new FileData.
func NewFileData(filename string) *FileData {
	return &FileData{
		Filename:   filename,
		LineData:   make(map[int]uint64),
		FuncData:   make(map[string]FuncData),
		BranchData: make(map[int][]BranchStatus),
		RegionData: make(map[Region]uint64),
	}
}

// LineCoverage calculates line coverage for the file.
func (file *FileData) LineCoverage() Coverage {
	a, b := 0, 0

	for _, v := range file.LineData {
		if v != 0 {
			a++
		}
		b++
	}
	return Coverage{a, b}
}

// FuncCoverage calculates function coverage for the file.
func (file *FileData) FuncCoverage() Coverage {
	a, b := 0, 0

	for _, v := range file.FuncData {
		if v.HitCount != 0 {
			a++
		}
		b++
	}
	return Coverage{a, b}
}

// BranchCoverage calculates branch coverage for the file.
func (file *FileData) BranchCoverage() Coverage {
	a, b := 0, 0

	for _, v := range file.BranchData {
		for _, v := range v {
			if v == BranchTaken {
				a++
			}
			b++
		}
	}
	return Coverage{a, b}
}

// RegionCoverage calculates region coverage for the file.
func (file *FileData) RegionCoverage() Coverage {
	a, b := 0, 0

	for _, v := range file.RegionData {
		if v != 0 {
			a++
		}
		b++
	}
	return Coverage{a, b}
}

// ConvertRegionToLineData will use hitcounts from region data to infer hit
// counts for line data.
//
// This method cannot be called if the line data already exists.
func (file *FileData) ConvertRegionToLineData() {
	if len(file.LineData) != 0 {
		panic("can not convert region data to line data if line data already present")
	}

	for k, hitCount := range file.RegionData {
		for i := k.StartLine; i <= k.EndLine; i++ {
			file.AppendLineCountData(i, hitCount)
		}
	}
}

// AppendLineCountData appends hit count data for a line.
func (file *FileData) AppendLineCountData(lineNo int, hitCount uint64) {
	file.LineData[lineNo] += hitCount
}

// AppendFunctionData appends hit count data for a function.
func (file *FileData) AppendFunctionData(funcName string, funcStart int, hitCount uint64) {
	if v, ok := file.FuncData[funcName]; ok {
		v.HitCount += hitCount
		file.FuncData[funcName] = v
	} else {
		file.FuncData[funcName] = FuncData{
			StartLine: funcStart,
			HitCount:  hitCount,
		}
	}
}

// AppendBranchData appends hit count data for a branch.
func (file *FileData) AppendBranchData(lineNo int, status BranchStatus) {
	tmp := file.BranchData[lineNo]
	tmp = append(tmp, status)
	file.BranchData[lineNo] = tmp
}

// AppendRegionData appends hit count data for a region.
func (file *FileData) AppendRegionData(startLine, startByte, endLine, endByte int, hitCount uint64) {
	region := Region{startLine, startByte, endLine, endByte}
	file.RegionData[region] += hitCount
}

// FileDataSet maintains coverage statistics for multiple files.
type FileDataSet map[string]*FileData

// FileData returns the data for a particular file in the set.
func (fds FileDataSet) FileData(filename string) *FileData {
	if tmp, ok := fds[filename]; ok {
		return tmp
	}

	tmp := NewFileData(filename)
	fds[filename] = tmp
	return tmp
}

// LineCoverage calculates line coverage over all of the files in the set.
func (fds FileDataSet) LineCoverage() Coverage {
	lcov := Coverage{}
	for _, data := range fds {
		lcov = lcov.Add(data.LineCoverage())
	}
	return lcov
}

// RegionCoverage calculates region coverage over all of the files in the set.
func (fds FileDataSet) RegionCoverage() Coverage {
	rcov := Coverage{}
	for _, data := range fds {
		rcov = rcov.Add(data.RegionCoverage())
	}
	return rcov
}

// ConvertRegionToLineData will use hitcounts from region data to infer hit
// counts for line data for all of the files in the set.
//
// This method cannot be called if the line data already exists.
func (fds FileDataSet) ConvertRegionToLineData() {
	for _, data := range fds {
		if len(data.LineData) == 0 {
			data.ConvertRegionToLineData()
		}
	}
}
