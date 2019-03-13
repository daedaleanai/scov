package main

import (
	"bufio"
	"fmt"
	"io"
)

func createTextReport(filename string, data map[string]FileData) error {
	w, err := Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = writeTextReport(w.File(), data)
	w.Keep(err)
	return err
}

func writeTextReport(writer io.Writer, data map[string]FileData) error {
	LCov := Coverage{}
	FCov := Coverage{}

	w := bufio.NewWriter(writer)

	for name, data := range data {
		lcov := data.LineCoverage()
		LCov.Accumulate(lcov)
		fcov := data.FuncCoverage()
		FCov.Accumulate(fcov)

		fmt.Fprintf(w, "%5.1f%%\t%5.1f%%\t%s\n", lcov.P(), fcov.P(), name)
	}
	fmt.Fprintf(w, "------\t------\t\n")
	fmt.Fprintf(w, "%5.1f%%\t%5.1f%%\tOverall\n", LCov.P(), FCov.P())
	return w.Flush()
}
