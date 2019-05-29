package main

import (
	"encoding/json"
	"errors"
	"io"
)

type LLVMData struct {
	Version string      `json:"version"`
	Type    string      `json:"type"`
	Data    []LLVMData2 `json:"data"`
}

type LLVMData2 struct {
	Files     []LLVMFile     `json:"files"`
	Functions []LLVMFunction `json:"functions"`
}

type LLVMFile struct {
	Filename string  `json:"filename"`
	Segments [][]int `json:"segments"`
}

type LLVMFunction struct {
	Name      string   `json:"name"`
	Count     uint64   `json:"count"`
	Regions   [][]int  `json:"regions"`
	Filenames []string `json:"filenames"`
}

func loadLLVMFile(fds FileDataSet, file io.Reader) error {
	data := LLVMData{}

	err := json.NewDecoder(file).Decode(&data)
	if err != nil {
		return err
	}

	if data.Type != "llvm.coverage.json.export" {
		return errors.New("incorrect type for JSON data from LLVM: " + data.Type)
	}

	for _, v := range data.Data {
		for _, w := range v.Files {
			filename := normalizeSourceFilename(w.Filename)
			currentData := fds.FileData(filename)

			for i := range w.Segments[:len(w.Segments)-1] {
				startLine := w.Segments[i][0]
				endLine := w.Segments[i+1][0]
				count := w.Segments[i][2]
				hasCount := (w.Segments[i][3] != 0)
				isRegionEntry := (w.Segments[i][4] != 0)

				if isRegionEntry && hasCount {
					for i := startLine; i <= endLine; i++ {
						applyLCountRecord(currentData, i, uint64(count))
					}
				}
			}
		}
		for _, w := range v.Functions {
			filename := normalizeSourceFilename(w.Filenames[0])
			currentData := fds.FileData(filename)

			applyFunctionRecord(currentData, w.Name, w.Regions[0][0], w.Count)
		}
	}

	return nil
}
