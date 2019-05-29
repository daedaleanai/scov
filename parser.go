package main

import (
	"os"
	"path/filepath"
)

type Parser int

const (
	ParserGCov Parser = iota
	ParserLCov
	ParserGCovJS
	ParserLLVM
	ParserGo
)

func identifyFileType(filename string) (Parser, bool) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".info":
		return ParserLCov, true
	case ".gcov":
		return ParserGCov, true
	case ".gz":
		return ParserGCovJS, true
	case ".json":
		return ParserLLVM, true
	case ".out":
		return ParserGo, true
	default:
		return 0, false
	}
}

func (p Parser) loadFile(data FileDataSet, file *os.File) error {
	switch p {
	case ParserLCov:
		return loadLCovFile(data, file)
	case ParserGCov:
		return loadGCovFile(data, file)
	case ParserGCovJS:
		return loadGCovJSFile(data, file)
	case ParserLLVM:
		return loadLLVMFile(data, file)
	case ParserGo:
		return loadGoFile(data, file)
	}

	panic("Unreachable")
}
