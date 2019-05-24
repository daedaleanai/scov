package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
			value = normalizeSourceFilename(value)
			currentData = fds.FileData(value)

		case "FNF": // Functions founds
			// ignore

		case "FNH": // Functions hit
			// ignore

		case "FN": // Function
			funcName, funcStart, err := parseFNRecord(value)
			if err != nil {
				return err
			}
			applyFunctionRecord(currentData, funcName, funcStart, 0)

		case "FNDA": // Function data
			funcName, hitCount, err := parseFNDARecord(value)
			if err != nil {
				return err
			}
			applyFunctionRecord(currentData, funcName, 0, hitCount)

		case "LF": // Lines founds
			// ignore

		case "LH": // Lines hit
			// ignore

		case "DA": // Line data
			lineNo, hitCount, err := parseDARecord(value)
			if err != nil {
				return err
			}
			applyLCountRecord(currentData, lineNo, hitCount)

		case "BRF": // Branches found
			// ignore

		case "BRH": // Branches hit
			// ignore

		case "BRDA": // Branch data
			lineNo, branchStatus, err := parseBRDARecord(value)
			if err != nil {
				return err
			}
			applyBranchRecord(currentData, lineNo, branchStatus)

		default:
			// Unknown records are ignored.  If future versions of the file
			// format introduce new records, we don't want to have an error.		}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func normalizeSourceFilename(filename string) string {
	base, err := filepath.Abs(*srcdir)
	if err != nil {
		fmt.Println("Error:  ", err.Error())
		panic(err)
	}

	if strings.HasPrefix(filename, base) {
		filename, _ = filepath.Rel(base, filename)
	}
	return filename
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

	hitCount, err := strconv.ParseInt(values[3], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse branch record: %s", err)
	}

	if hitCount > 0 {
		return lineNo, BranchTaken, nil
	}
	return lineNo, BranchNotTaken, nil
}
