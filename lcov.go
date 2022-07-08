package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func loadLCovFile(fds FileDataSet, file *os.File) error {
	currentData := (*FileData)(nil)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t, value := recordType(scanner.Text())
		switch t {
		case "TN": // Title
			// ignore

		case "SF": // Source file
			currentData = fds.FileData(value)

		case "FN": // Function
			funcName, funcStart, err := parseFNRecord(value)
			if err != nil {
				return err
			}
			currentData.AppendFunctionData(funcName, funcStart, 0)

		case "FNDA": // Function data
			funcName, hitCount, err := parseFNDARecord(value)
			if err != nil {
				return err
			}
			currentData.AppendFunctionData(funcName, 0, hitCount)

		case "DA": // Line data
			lineNo, hitCount, err := parseDARecord(value)
			if err != nil {
				return err
			}
			currentData.AppendLineCountData(lineNo, hitCount)

		case "BRDA": // Branch data
			lineNo, branchStatus, err := parseBRDARecord(value)
			if err != nil {
				return err
			}
			currentData.AppendBranchData(lineNo, branchStatus)

		default:
			// Unknown records are ignored.  If future versions of the file
			// format introduce new records, we don't want to have an error.
			//
			// We are also skipping known records:
			//    LF  lines found
			//    LH  lines hit
			//    FNF functions found
			//    FNH functions hit
			//    BRF branches found
			//    BRH branches hit
			//
			// The above records provides summaries of the counts in the data
			// records, but we will calculate that ourselves.
		}
	}

	return scanner.Err()
}

func parseDARecord(value string) (lineNo int, hitCount uint64, err error) {
	buffer := [4]string{}
	values := splitOnComma(buffer[:], value)

	if l := len(values); l != 2 {
		return 0, 0, fmt.Errorf("can't parse DA record")
	}

	lineNoTmp, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse DA record: %s", err)
	}
	lineNo = int(lineNoTmp)
	hitCount, err = strconv.ParseUint(values[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse DA record: %s", err)
	}

	return lineNo, hitCount, nil
}

func parseFNRecord(value string) (funcName string, funcStart int, err error) {
	buffer := [4]string{}
	values := splitOnComma(buffer[:], value)

	if len(values) != 2 {
		return "", 0, fmt.Errorf("can't parse function record")
	}

	line, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("can't parse function record: %s", err)
	}
	funcName = values[1]
	return funcName, int(line), nil
}

func parseFNDARecord(value string) (funcName string, hitCount uint64, err error) {
	buffer := [4]string{}
	values := splitOnComma(buffer[:], value)

	if len(values) != 2 {
		return "", 0, fmt.Errorf("can't parse function data record")
	}

	hitCount, err = strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("can't parse function data record: %s", err)
	}
	funcName = values[1]
	return funcName, hitCount, nil
}

func parseBRDARecord(value string) (lineNo int, status BranchStatus, err error) {
	buffer := [4]string{}
	values := splitOnComma(buffer[:], value)

	if len(values) != 4 {
		return 0, 0, fmt.Errorf("can't parse branch record")
	}

	lineNoTmp, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse branch record: %s", err)
	}
	lineNo = int(lineNoTmp)

	var hitCount int64 = 0

	if values[3] != "-" {
		hitCount, err = strconv.ParseInt(values[3], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("can't parse branch record: %s", err)
		}
	}

	if hitCount > 0 {
		return lineNo, BranchTaken, nil
	}
	return lineNo, BranchNotTaken, nil
}
