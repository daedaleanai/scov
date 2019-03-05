package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	lineCountData map[string]map[int]uint64
	funcCountData map[string]map[string]uint64
	external      = flag.Bool("external", false, "Set whether external files to be included")
	help          = flag.Bool("h", false, "Request help")
	srcdir        = flag.String("srcdir", ".", "Path for the source directory")
	outdir        = flag.String("outdir", ".", "Path for the output")
	title         = flag.String("title", "LCovHTML", "Title for the HTML pages")
)

func main() {
	// Initialize global maps used to track line and function coverage
	lineCountData = make(map[string]map[int]uint64)
	funcCountData = make(map[string]map[string]uint64)

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *external {
		fmt.Println("including external files")
	}

	for _, name := range flag.Args() {
		loadFile(name)
	}

	err := buildText(os.Stdout)
	if err != nil {
		panic(err)
	}
	err = buildHtml(*outdir)
	if err != nil {
		panic(err)
	}

}

func recordType(line string) (string, string) {
	ndx := strings.IndexByte(line, ':')
	if ndx < 0 {
		return line, ""
	}
	return line[:ndx], line[ndx+1:]
}

func loadFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	currentFile := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t, value := recordType(scanner.Text())
		switch t {
		case "version":
			//fmt.Println("version", value)

		case "file":
			currentFile = value

		case "function":
			funcName, hitCount, err := parseFunctionRecord(value)
			if err != nil {
				return err
			}
			applyFunctionRecord(currentFile, funcName, hitCount)

		case "lcount":
			lineNo, hitCount, err := parseLCountRecord(value)
			if err != nil {
				return err
			}
			applyLCountRecord(currentFile, lineNo, hitCount)

		default:
			panic("unknown record")
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

func applyFunctionRecord(file string, funcName string, hitCount uint64) {
	if m, ok := funcCountData[file]; ok {
		m[funcName] += hitCount
	} else {
		m := make(map[string]uint64)
		m[funcName] += hitCount
		funcCountData[file] = m
	}
}

func parseLCountRecord(value string) (lineNo int, hitCount uint64, err error) {
	values := strings.Split(value, ",")
	if len(values) != 2 {
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

func applyLCountRecord(file string, lineNo int, hitCount uint64) {
	if m, ok := lineCountData[file]; ok {
		m[lineNo] += hitCount
	} else {
		m := make(map[int]uint64)
		m[lineNo] += hitCount
		lineCountData[file] = m
	}
}

func lineCoverageForFile(data map[int]uint64) (int, int) {
	a, b := 0, 0

	for _, v := range data {
		if v != 0 {
			a++
		}
		b++
	}
	return a, b
}

func funcCoverageForFile(data map[string]uint64) (int, int) {
	a, b := 0, 0

	for _, v := range data {
		if v != 0 {
			a++
		}
		b++
	}
	return a, b
}
