package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func recordType(line string) (string, string) {
	ndx := strings.IndexByte(line, ':')
	if ndx < 0 {
		return line, ""
	}
	return line[:ndx], line[ndx+1:]
}

func loadGCovFile(fds FileDataSet, file *os.File) error {
	currentData := (*FileData)(nil)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t, value := recordType(scanner.Text())
		switch t {
		case "version":
			//fmt.Println("version", value)

		case "file":
			currentData = fds.FileData(value)

		case "function":
			funcName, _, hitCount, err := parseFunctionRecord(value)
			if err != nil {
				return err
			}
			applyFunctionRecord(currentData, funcName, hitCount)

		case "lcount":
			lineNo, hitCount, err := parseLCountRecord(value)
			if err != nil {
				return err
			}
			applyLCountRecord(currentData, lineNo, hitCount)

		case "branch":
			lineNo, branchStatus, err := parseBranchRecord(value)
			if err != nil {
				return err
			}
			applyBranchRecord(currentData, lineNo, branchStatus)

		default:
			return fmt.Errorf("unrecognized record: %s", t)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func parseFunctionRecord(value string) (funcName string, funcStart int, hitCount uint64, err error) {
	values := strings.Split(value, ",")
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

func applyFunctionRecord(data *FileData, funcName string, hitCount uint64) {
	data.FuncData[funcName] += hitCount
}

func parseLCountRecord(value string) (lineNo int, hitCount uint64, err error) {
	values := strings.Split(value, ",")
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

func applyLCountRecord(data *FileData, lineNo int, hitCount uint64) {
	data.LineData[lineNo] += hitCount
}

func parseBranchRecord(value string) (lineNo int, status BranchStatus, err error) {
	values := strings.Split(value, ",")
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

func applyBranchRecord(data *FileData, lineNo int, status BranchStatus) {
	tmp := data.BranchData[lineNo]
	tmp = append(tmp, status)
	data.BranchData[lineNo] = tmp
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
