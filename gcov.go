package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func loadGCovFile(fds FileDataSet, file *os.File) error {
	currentData := (*FileData)(nil)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t, value := recordType(scanner.Text())
		switch t {
		case "version":
			// Don't need to track the version, as the record variations can be
			// recognized when the fields are split.  May want to print the
			// version for debugging, but otherwise we can ignore this record.
			//
			// fmt.Println("version", value)

		case "file":
			currentData = fds.FileData(value)

		case "function":
			funcName, funcStart, hitCount, err := parseFunctionRecord(value)
			if err != nil {
				return err
			}
			currentData.AppendFunctionData(funcName, funcStart, hitCount)

		case "lcount":
			lineNo, hitCount, err := parseLCountRecord(value)
			if err != nil {
				return err
			}
			currentData.AppendLineCountData(lineNo, hitCount)

		case "branch":
			lineNo, branchStatus, err := parseBranchRecord(value)
			if err != nil {
				return err
			}
			currentData.AppendBranchData(lineNo, branchStatus)

		default:
			// Unknown records are ignored.  If future versions of the file
			// format introduce new records, we don't want to have an error.
		}
	}

	return scanner.Err()
}

func parseFunctionRecord(value string) (funcName string, funcStart int, hitCount uint64, err error) {
	buffer := [4]string{}
	values := splitOnComma(buffer[:], value)

	if len(values) == 3 {
		tmp, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return "", 0, 0, fmt.Errorf("can't parse function record: %s", err)
		}
		funcStart = int(tmp)

		hitCount, err = strconv.ParseUint(values[1], 10, 64)
		if err != nil {
			return "", 0, 0, fmt.Errorf("can't parse function record: %s", err)
		}
		funcName = values[2]
		return funcName, funcStart, hitCount, nil
	} else if len(values) == 4 {
		// The first two fields are the line number range for the function.
		// We are using the start.
		tmp, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return "", 0, 0, fmt.Errorf("can't parse function record: %s", err)
		}
		funcStart = int(tmp)

		hitCount, err = strconv.ParseUint(values[2], 10, 64)
		if err != nil {
			return "", 0, 0, fmt.Errorf("can't parse function record: %s", err)
		}
		funcName = values[3]
		return funcName, funcStart, hitCount, nil
	}

	return "", 0, 0, fmt.Errorf("can't parse function record")
}

func parseLCountRecord(value string) (lineNo int, hitCount uint64, err error) {
	buffer := [4]string{}
	values := splitOnComma(buffer[:], value)

	if l := len(values); l != 2 && l != 3 {
		return 0, 0, fmt.Errorf("can't parse lcount record")
	}

	lineNoTmp, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse lcount record: %s", err)
	}
	lineNo = int(lineNoTmp)
	hitCount, err = strconv.ParseUint(values[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse lcount record: %s", err)
	}

	return lineNo, hitCount, nil
}

func parseBranchRecord(value string) (lineNo int, status BranchStatus, err error) {
	buffer := [4]string{}
	values := splitOnComma(buffer[:], value)

	if len(values) != 2 {
		return 0, 0, fmt.Errorf("can't parse branch record")
	}

	lineNoTmp, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse branch record: %s", err)
	}
	lineNo = int(lineNoTmp)

	switch values[1] {
	case "taken":
		return lineNo, BranchTaken, nil
	case "nottaken":
		return lineNo, BranchNotTaken, nil
	case "notexec":
		return lineNo, BranchNotExec, nil
	}

	return 0, 0, fmt.Errorf("can't parse branch record: unrecognized branch status")
}

func filterExternalFileData(fileData map[string]*FileData, external bool) map[string]*FileData {
	if external {
		return fileData
	}

	for key := range fileData {
		if filepath.IsAbs(key) {
			delete(fileData, key)
		}
	}
	return fileData
}
