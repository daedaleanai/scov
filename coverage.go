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

// Accumulate revises the counts to accumulate data over several sub-scopes.
func (c *Coverage) Accumulate(delta Coverage) {
	c.Hits += delta.Hits
	c.Total += delta.Total
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

// BranchStatus indicates whether a branch was taken, not taken, or if the
// conditional was never executed.
type BranchStatus uint8

const (
	BranchTaken BranchStatus = iota
	BranchNotTaken
	BranchNotExec
)

type FileData struct {
	Filename   string
	LineData   map[int]uint64
	FuncData   map[string]uint64
	BranchData map[int][]BranchStatus
}

func NewFileData(filename string) *FileData {
	return &FileData{
		Filename:   filename,
		LineData:   make(map[int]uint64),
		FuncData:   make(map[string]uint64),
		BranchData: make(map[int][]BranchStatus),
	}
}

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

func (file *FileData) FuncCoverage() Coverage {
	a, b := 0, 0

	for _, v := range file.FuncData {
		if v != 0 {
			a++
		}
		b++
	}
	return Coverage{a, b}
}

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

type FileDataSet map[string]*FileData

func (fds FileDataSet) FileData(filename string) *FileData {
	if tmp, ok := fds[filename]; ok {
		return tmp
	}

	tmp := NewFileData(filename)
	fds[filename] = tmp
	return tmp
}

func (fds FileDataSet) LineCoverage() Coverage {
	lcov := Coverage{}
	for _, data := range fds {
		lcov.Accumulate(data.LineCoverage())
	}
	return lcov
}
