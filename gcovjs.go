package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
)

type GCovData struct {
	DataFile   string     `json:"data_file"`
	GCCVersion string     `json:"gcc_version"`
	Files      []GCovFile `json:"files"`
}

type GCovFile struct {
	File      string         `json:"file"`
	Functions []GCovFunction `json:"functions"`
	Lines     []GCovLine     `json:"lines"`
}

type GCovFunction struct {
	Name           string `json:"name"`
	StartLine      int    `json:"start_line"`
	ExecutionCount uint64 `json:"execution_count"`
}

type GCovLine struct {
	LineNumber int          `json:"line_number"`
	Count      uint64       `json:"count"`
	Branches   []GCovBranch `json:"branches"`
}

type GCovBranch struct {
	Fallthrough bool   `json:"fallthrough"`
	Count       uint64 `json:"count"`
	Throw       bool   `json:"throw"`
}

func loadGCovJSFile(fds FileDataSet, file io.Reader) error {
	jsonData := GCovData{}

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	err = json.NewDecoder(gz).Decode(&jsonData)
	if err != nil {
		return err
	}

	for _, v := range jsonData.Files {
		filename := v.File
		currentData := fds.FileData(filename)

		for _, u := range v.Functions {
			applyFunctionRecord(currentData, u.Name, u.StartLine, u.ExecutionCount)
		}

		for _, u := range v.Lines {
			applyLCountRecord(currentData, u.LineNumber, u.Count)

			for _, b := range u.Branches {
				status := BranchNotTaken
				if b.Count > 0 {
					status = BranchTaken
				}
				applyBranchRecord(currentData, u.LineNumber, status)
			}
		}
	}

	return nil
}
