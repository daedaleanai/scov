package main

import (
	"fmt"
	"io"
	"strings"
)

func buildText(out io.Writer, data map[string]FileData) error {
	LCov := Coverage{}
	FCov := Coverage{}

	for name, data := range data {
		if *external || !strings.HasPrefix(name, "/") {
			lcov := data.LineCoverage()
			LCov.Update(lcov)
			fcov := data.FuncCoverage()
			FCov.Update(fcov)

			_, err := fmt.Fprintf(out, "%5.1f%%\t%5.1f%%\t%s\n", lcov.Percentage(), fcov.Percentage(), name)
			if err != nil {
				return err
			}
		}
	}
	fmt.Fprintf(out, "------\t------\t\n")
	fmt.Fprintf(out, "%5.1f%%\t%5.1f%%\tOverall\n", LCov.Percentage(), FCov.Percentage())
	return nil
}
