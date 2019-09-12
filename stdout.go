package main

import (
	"bufio"
	"fmt"
	"io"
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

	// Head
	w.WriteString(" Lines\t Funcs\n")
	w.WriteString("------\t------\n")

	// Body
	for _, i := range report.Files {
		fmt.Fprintf(w, "%5.1f%%\t%5.1f%%\t%s\n", i.LCoverage.P(), i.FCoverage.P(), i.Name)
	}

	// Foot
	w.WriteString("------\t------\n")
	fmt.Fprintf(w, "%5.1f%%\t%5.1f%%\tOverall\n", report.LCoverage.P(), report.FCoverage.P())
	return w.Flush()
}
