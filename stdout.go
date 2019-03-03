package main

import (
	"fmt"
	"os"
	"strings"
)

func buildText(out *os.File) error {
	LCov := Coverage{}
	FCov := Coverage{}

	for name, data := range lineCountData {
		if *external || !strings.HasPrefix(name, "/") {
			a, b := lineCoverageForFile(data)
			lcov := Coverage{a, b}
			LCov.Update(a, b)
			a, b = funcCoverageForFile(funcCountData[name])
			fcov := Coverage{a, b}
			FCov.Update(a, b)

			fmt.Fprintf(out, "%.1f\t%.1f\t%s\n", lcov.Percentage(), fcov.Percentage(), name)
		}
	}
	fmt.Fprintf(out, "-----\t-----\t\n")
	fmt.Fprintf(out, "%.1f\t%.1f\tOverall\n", LCov.Percentage(), FCov.Percentage())
	return nil
}
