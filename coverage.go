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

type FileData struct {
	Filename string
	LineData map[int]uint64
	FuncData map[string]uint64
}

func NewFileData(filename string) FileData {
	return FileData{
		Filename: filename,
		LineData: make(map[int]uint64),
		FuncData: make(map[string]uint64),
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
