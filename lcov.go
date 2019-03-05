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
	err = buildHtml(".")
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
			processFunction(currentFile, value)

		case "lcount":
			processLCount(currentFile, value)

		default:
			panic("unknown record")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processFunction(file string, value string) {
	values := strings.Split(value, ",")
	lineCount, _ := strconv.ParseUint(values[2], 10, 64)
	funcName := values[3]

	if m, ok := funcCountData[file]; ok {
		m[funcName] += lineCount
	} else {
		m := make(map[string]uint64)
		m[funcName] += lineCount
		funcCountData[file] = m
	}
}

func processLCount(file string, value string) {
	values := strings.Split(value, ",")
	lineNo, _ := strconv.ParseInt(values[0], 10, 64)
	lineCount, _ := strconv.ParseUint(values[1], 10, 64)

	if m, ok := lineCountData[file]; ok {
		m[int(lineNo)] += lineCount
	} else {
		m := make(map[int]uint64)
		m[int(lineNo)] += lineCount
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
