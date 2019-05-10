package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	help     = flag.Bool("h", false, "Request help")
	version  = flag.Bool("v", false, "Request version information")
	external = flag.Bool("external", false, "Set whether external files to be included")
	exclude  = flag.String("exclude", "", "Exclude source files that match the regular expression")
	srcdir   = flag.String("srcdir", ".", "Path for the source directory")
	srcid    = flag.String("srcid", "", "String to identify revision of source")
	title    = flag.String("title", "GCovHTML", "Title for the HTML pages")
	htmldir  = flag.String("htmldir", ".", "Path for the HTML output")
	htmljs   = flag.Bool("htmljs", false, "Use javascript to enhance reports")
	text     = flag.String("text", "", "Filename for text report, use - to direct the report to stdout")
)

var (
	versionInformation = "(development)"
)

func main() {
	// Initialize global maps used to track line and function coverage
	fileData := make(map[string]*FileData)

	flag.Parse()
	if ok := handleRequestFlags(os.Stdout, *help, *version); ok {
		os.Exit(0)
	}

	for _, name := range flag.Args() {
		err := loadFile(fileData, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not load data: %s\n", err)
			os.Exit(1)
		}
	}
	fileData = filterExternalFileData(fileData, *external)
	fileData = filterExcludedFileData(os.Stderr, fileData, *exclude)
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

func handleRequestFlags(out io.Writer, help, version bool) bool {
	if help {
		flag.CommandLine.SetOutput(out)
		flag.Usage()
		return true
	}
	if version {
		fmt.Fprintf(out, "gcovhtml %s\n", versionInformation)
		return true
	}

	return false
}

func loadFile(data FileDataSet, filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// If this is a directory, we will need to process the directory entries.
	stat, err := file.Stat()
	if err != nil {
		return err
	} else if stat.IsDir() {
		return loadFilesFromDir(data, file)
	}

	parser, ok := identifyFileType(filename)
	if !ok {
		return fmt.Errorf("unrecognized file extension: %s", filepath.Ext(filename))
	}

	return parser.loadFile(data, file)
}

func loadFilesFromDir(data FileDataSet, file *os.File) error {
	names, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, v := range names {
		if _, ok := identifyFileType(v); ok {
			err := loadFile(data, filepath.Join(file.Name(), v))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func filterExcludedFileData(out io.Writer, fileData FileDataSet, filter string) FileDataSet {
	if filter == "" {
		return fileData
	}

	re, err := regexp.Compile(filter)
	if err != nil {
		fmt.Fprintf(out, "warning: did not apply filter to exclude files: %s\n", err)
		return fileData
	}

	for key := range fileData {
		if re.FindString(key) != "" {
			delete(fileData, key)
		}
	}
	return fileData
}
