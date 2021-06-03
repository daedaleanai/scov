package main

import (
	"bufio"
	"fmt"
	"io"

	"git.sr.ht/~rj/sgr"
)

func createTextReport(filename string, report *Report) error {
	w, err := Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = writeTextReport(w.File(), report)
	w.Keep(err)
	return err
}

func writeTextReport(writer io.Writer, report *Report) error {
	w := bufio.NewWriter(writer)
	f := sgr.NewFormatterForWriter(writer)

	// Head
	_, _ = fmt.Fprintf(w, "%v\n%v\n",
		f.Bold(" Lines\t Funcs\tBranch\tRegion"),
		f.Dim("------\t------\t------\t------"))

	// Body
	for _, i := range report.Files {
		fmt.Fprintf(w, "%6.1f%%\t%5.1f%%\t%5.1f%%\t%5.1f%%\t%s\n",
			i.LCoverage,
			i.FCoverage,
			i.BCoverage,
			i.RCoverage,
			i.Name)
	}

	// Foot
	fmt.Fprintf(w, "%v\n%v\n",
		f.Dim("------\t------\t------\t------"),
		f.Boldf("%5.1f%%\t%5.1f%%\t%5.1f%%\t%5.1f%%\tOverall",
			report.LCoverage,
			report.FCoverage,
			report.BCoverage,
			report.RCoverage))
	return w.Flush()
}
