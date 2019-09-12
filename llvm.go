package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
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
	Filename string        `json:"filename"`
	Segments []LLVMSegment `json:"segments"`
}

type LLVMFunction struct {
	Name      string   `json:"name"`
	Count     uint64   `json:"count"`
	Regions   [][]int  `json:"regions"`
	Filenames []string `json:"filenames"`
}

type LLVMSegment struct {
	Line          int
	Column        int
	Count         uint64
	HasCount      bool
	IsRegionEntry bool
}

func (s *LLVMSegment) UnmarshalJSON(data []byte) error {
	// Strip of the outer square brackets for the array.
	if data[0] != '[' || data[len(data)-1] != ']' {
		return errors.New("expected an array")
	}
	data = data[1 : len(data)-1]

	// Get the Line
	ndx := bytes.IndexByte(data, ',')
	if ndx < 1 {
		return errors.New("expected at least 5 elements in the array")
	}
	tmp, err := strconv.ParseInt(string(data[:ndx]), 10, 64)
	if err != nil {
		return err
	}
	s.Line = int(tmp)
	data = data[ndx+1:]

	// Get the Column
	ndx = bytes.IndexByte(data, ',')
	if ndx < 1 {
		return errors.New("expected at least 5 elements in the array")
	}
	tmp, err = strconv.ParseInt(string(data[:ndx]), 10, 64)
	if err != nil {
		return err
	}
	s.Column = int(tmp)
	data = data[ndx+1:]

	// Get the Count
	ndx = bytes.IndexByte(data, ',')
	if ndx < 1 {
		return errors.New("expected at least 5 elements in the array")
	}
	s.Count, err = strconv.ParseUint(string(data[:ndx]), 10, 64)
	if err != nil {
		return err
	}
	data = data[ndx+1:]

	// Get the HasCount
	ndx = bytes.IndexByte(data, ',')
	if ndx < 1 {
		return errors.New("expected at least 5 elements in the array")
	}
	s.HasCount, err = strconv.ParseBool(string(data[:ndx]))
	if err != nil {
		return err
	}
	data = data[ndx+1:]

	// Get the IsRegionEntry
	ndx = bytes.IndexByte(data, ',')
	if ndx >= 1 {
		// Ignore any new elements in the array.
		data = data[:ndx]
	}
	s.IsRegionEntry, err = strconv.ParseBool(string(data))
	if err != nil {
		return err
	}

	return nil
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
			currentData := fds.FileData(w.Filename)

			for i := range w.Segments[:len(w.Segments)-1] {
				startLine := w.Segments[i].Line
				endLine := w.Segments[i+1].Line
				count := w.Segments[i].Count
				hasCount := w.Segments[i].HasCount
				isRegionEntry := w.Segments[i].IsRegionEntry

				if isRegionEntry && hasCount {
					for i := startLine; i <= endLine; i++ {
						applyLCountRecord(currentData, i, uint64(count))
					}
				}
			}
		}
		for _, w := range v.Functions {
			currentData := fds.FileData(w.Filenames[0])

			applyFunctionRecord(currentData, w.Name, w.Regions[0][0], w.Count)
		}
	}

	return nil
}
