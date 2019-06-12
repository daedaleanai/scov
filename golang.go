package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func loadGoFile(fds FileDataSet, file *os.File) error {
	mode := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := scanner.Text()
		if mode == "" {
			if !strings.HasPrefix(record, "mode:") {
				return errors.New("format error: missing mode record")
			}
			mode = strings.TrimSpace(record[5:])
		} else {
			filename, start, end, _, hitCount, err := parseGoRecord(record)
			if err != nil {
				return err
			}

			currentData := fds.FileData(filename)
			for i := start.Line; i <= end.Line; i++ {
				applyLCountRecord(currentData, i, hitCount)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type Position struct {
	Line   int
	Column int
}

func (p Position) IsZero() bool {
	return p.Line == 0 && p.Column == 0
}

func parseGoRecord(record string) (string, Position, Position, int, uint64, error) {
	ndx := strings.IndexByte(record, ':')
	if ndx < 1 {
		return "", Position{}, Position{}, 0, 0, errors.New("could not find separator ':' in record")
	}
	filename := record[:ndx]
	record = record[ndx+1:]

	ndx = strings.IndexByte(record, ',')
	if ndx < 1 {
		return "", Position{}, Position{}, 0, 0, errors.New("could not find separator ',' in record")
	}
	start, err := parseGoPosition(record[:ndx])
	if err != nil {
		return "", Position{}, Position{}, 0, 0, err
	}
	record = record[ndx+1:]

	ndx = strings.IndexByte(record, ' ')
	if ndx < 1 {
		return "", Position{}, Position{}, 0, 0, errors.New("could not find separator ' ' in record")
	}
	end, err := parseGoPosition(record[:ndx])
	if err != nil {
		return "", Position{}, Position{}, 0, 0, err
	}
	record = record[ndx+1:]

	ndx = strings.IndexByte(record, ' ')
	if ndx < 1 {
		return "", Position{}, Position{}, 0, 0, errors.New("could not find separator ' ' in record")
	}
	nos, err := strconv.ParseInt(record[:ndx], 10, 64)
	if err != nil {
		return "", Position{}, Position{}, 0, 0, err
	}
	record = record[ndx+1:]

	count, err := strconv.ParseUint(record, 10, 64)
	return filename, start, end, int(nos), count, err
}

func parseGoPosition(field string) (Position, error) {
	ndx := strings.IndexByte(field, '.')
	if ndx < 1 {
		return Position{}, errors.New("could not parse position: missing field separator '.'")
	}

	line, err := strconv.ParseInt(field[:ndx], 10, 64)
	if err != nil {
		return Position{}, errors.New("could not parse position: " + err.Error())
	}
	col, err := strconv.ParseInt(field[ndx+1:], 10, 64)
	if err != nil {
		return Position{}, errors.New("could not parse position: " + err.Error())
	}

	return Position{int(line), int(col)}, nil
}
