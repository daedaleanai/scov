package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	external = flag.Bool("external", false, "Set whether external files to be included")
	help     = flag.Bool("h", false, "Request help")
	srcdir   = flag.String("srcdir", ".", "Path for the source directory")
	title    = flag.String("title", "GCovHTML", "Title for the HTML pages")
	htmldir  = flag.String("htmldir", ".", "Path for the HTML output")
	text     = flag.String("text", "", "Filename for text report, use - for stdout")
)

func main() {
	// Initialize global maps used to track line and function coverage
	fileData := make(map[string]*FileData)

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	for _, name := range flag.Args() {
		err := loadFile(fileData, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not load data: %s\n", err)
			os.Exit(1)
		}
	}
	fileData = filterFileData(fileData, *external)
	if len(fileData) == 0 {
		fmt.Fprintf(os.Stderr, "error: no file data present\n")
		os.Exit(1)
	}

	if *text != "" {
		err := createTextReport(*text, fileData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not create text report: %s\n", err)
			os.Exit(1)
		}
	} else {
		lcov := Coverage{}
		for _, data := range fileData {
			lcov.Accumulate(data.LineCoverage())
		}
		fmt.Fprintf(os.Stdout, "Coverage: %.1f%%\n", lcov.P())
	}

	if *htmldir != "" {
		err := createHTML(*htmldir, fileData, time.Now().UTC())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not create HTML report: %s\n", err)
			os.Exit(1)
		}
	}
}

func recordType(line string) (string, string) {
	ndx := strings.IndexByte(line, ':')
	if ndx < 0 {
		return line, ""
	}
	return line[:ndx], line[ndx+1:]
}

func loadFile(data map[string]*FileData, name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	currentData := (*FileData)(nil)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t, value := recordType(scanner.Text())
		switch t {
		case "version":
			//fmt.Println("version", value)

		case "file":
			if tmp, ok := data[value]; ok {
				currentData = tmp
			} else {
				currentData = NewFileData(value)
				data[value] = currentData
			}

		case "function":
			funcName, hitCount, err := parseFunctionRecord(value)
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

func parseFunctionRecord(value string) (funcName string, hitCount uint64, err error) {
	values := strings.Split(value, ",")
	if len(values) == 3 {
		// The first field is the line number for the function.
		// We are not using that information.

		hitCount, err = strconv.ParseUint(values[1], 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("can't parse function record: %s", err)
		}
		funcName = values[2]
		return funcName, hitCount, nil
	} else if len(values) == 4 {
		// The first two fields are the line number range for the function.
		// We are not using that information.

		hitCount, err = strconv.ParseUint(values[2], 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("can't parse function record: %s", err)
		}
		funcName = values[3]
		return funcName, hitCount, nil
	}

	return "", 0, fmt.Errorf("can't parse function record")
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

func filterFileData(fileData map[string]*FileData, external bool) map[string]*FileData {
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
