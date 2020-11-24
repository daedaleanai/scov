package main

import (
	"sort"
	"time"
)

// A Report contains the configuration for a report, along with the calculated results.
type Report struct {
	// Metadata about the run
	Title      string
	TestID     string
	SrcID      string
	SrcDir     string
	ProjectURL string

	// Configuration
	AllowHTMLScripting bool

	LCoverage Coverage
	FCoverage Coverage
	BCoverage Coverage
	RCoverage Coverage
	Files     []FileStatistics
	Funcs     []FuncStatistics
	Date      time.Time
}

// NewReport initializes a new report.
func NewReport(title string) *Report {
	return &Report{
		Title: title,
		Date:  time.Now().UTC(),
	}
}

// NewTestReport initializes a new report for use in testing.
func NewTestReport() *Report {
	return &Report{
		Title: "SCov",
		Date:  time.Date(2006, 01, 02, 15, 4, 5, 6, time.UTC),
	}
}

// FileStatistics is used to capture coverage statistics for a file.
type FileStatistics struct {
	Name      string
	LCoverage Coverage
	FCoverage Coverage
	BCoverage Coverage
	RCoverage Coverage
}

// FuncStatistics is used to capture data for a function.
type FuncStatistics struct {
	Name      string
	Filename  string
	StartLine int
	HitCount  uint64
}

// CollectStatistics iterates over the file data, and assembles coverage
// statistics for the set.  It also assembles coverage statistics for each
// source file, and each function.
func (r *Report) CollectStatistics(data map[string]*FileData) {
	// Preallocate space for our statistics
	files := make([]FileStatistics, 0, len(data))
	funcs := make([]FuncStatistics, 0, len(data))

	LCov := Coverage{}
	FCov := Coverage{}
	BCov := Coverage{}
	RCov := Coverage{}
	for filename, data := range data {
		stats := FileStatistics{Name: filename}

		stats.LCoverage = data.LineCoverage()
		LCov = LCov.Add(stats.LCoverage)
		stats.FCoverage = data.FuncCoverage()
		FCov = FCov.Add(stats.FCoverage)
		stats.BCoverage = data.BranchCoverage()
		BCov = BCov.Add(stats.BCoverage)
		stats.RCoverage = data.RegionCoverage()
		RCov = RCov.Add(stats.RCoverage)

		files = append(files, stats)

		for name, data := range data.FuncData {
			funcs = append(funcs, FuncStatistics{
				Name:      name,
				Filename:  filename,
				StartLine: data.StartLine,
				HitCount:  data.HitCount,
			})
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].Name < funcs[j].Name
	})

	r.LCoverage = LCov
	r.FCoverage = FCov
	r.BCoverage = BCov
	r.RCoverage = RCov
	r.Files = files
	r.Funcs = funcs
}

// UnixDate returns the date of the report formatted to the format time.UnixDate.
func (r *Report) UnixDate() string {
	return r.Date.Format(time.UnixDate)
}
