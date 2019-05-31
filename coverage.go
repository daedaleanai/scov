package main

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

// Q returns the percentag eof lines or functions that were not executed.
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

// FileData maintains coverage statistics for a single file.
type FileData struct {
	Filename   string
	LineData   map[int]uint64
	FuncData   map[string]FuncData
	BranchData map[int][]BranchStatus
}

// NewFileData initializes a new FileData.
func NewFileData(filename string) *FileData {
	return &FileData{
		Filename:   filename,
		LineData:   make(map[int]uint64),
		FuncData:   make(map[string]FuncData),
		BranchData: make(map[int][]BranchStatus),
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
