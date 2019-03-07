package main

import (
	"bufio"
	"fmt"
	"io"
)

func writeTextReport(writer io.Writer, data map[string]FileData) error {
	LCov := Coverage{}
	FCov := Coverage{}

	w := bufio.NewWriter(writer)

	for name, data := range data {
		lcov := data.LineCoverage()
		LCov.Update(lcov)
		fcov := data.FuncCoverage()
		FCov.Update(fcov)

		fmt.Fprintf(w, "%5.1f%%\t%5.1f%%\t%s\n", lcov.Percentage(), fcov.Percentage(), name)
	}
	fmt.Fprintf(w, "------\t------\t\n")
	fmt.Fprintf(w, "%5.1f%%\t%5.1f%%\tOverall\n", LCov.Percentage(), FCov.Percentage())
	return w.Flush()
}
