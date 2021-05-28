package main

import (
	"fmt"
	"io"

	"git.sr.ht/~rj/sgr"
	"git.sr.ht/~rj/sgr/plot"
)

func writeStdoutReport(w io.Writer, report *Report) {
	f := sgr.NewFormatterForWriter(w)

	writeStdoutCoverage(w, f, "Line coverage", report.LCoverage)
	writeStdoutCoverage(w, f, "Func coverage", report.FCoverage)
	writeStdoutCoverage(w, f, "Branch coverage", report.BCoverage)
	writeStdoutCoverage(w, f, "Region coverage", report.RCoverage)
}

func writeStdoutCoverage(w io.Writer, f *sgr.Formatter, name string, cov Coverage) {
	if !cov.Valid() {
		fmt.Fprintf(w, "%15s:  No data\n", name)
		return
	}

	fmt.Fprintf(w, "%15s: [%20s] %5.1f%%  (%d/%d)\n",
		name,
		plot.HorizontalBar(
			f,
			float64(cov.P()*0.01),
			ratingToColor(cov.Rating())),
		cov.P(),
		cov.Hits,
		cov.Total)
}

func ratingToColor(r CoverageRating) sgr.Color {
	switch r {
	case LowCoverage:
		return sgr.Red
	case MediumCoverage:
		return sgr.BrightYellow
	case HighCoverage:
		return sgr.BrightGreen
	default:
		return sgr.Default
	}
}
