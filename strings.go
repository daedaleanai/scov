package main

import (
	"strings"
)

// recordType parses the record type, and the fields, from a line for both
// gcov and lcov records.
func recordType(line string) (string, string) {
	ndx := strings.IndexByte(line, ':')
	if ndx < 0 {
		return line, ""
	}
	return line[:ndx], line[ndx+1:]
}

// splitOnComma is an optimized version of strings.Split for the special case
// required when parsing gcov and lcov files.  We don't need a multi-byte
// separator, so we can optimize on IndexByte for searching.  Additionally,
// we can reuse a buffer to avoid unnecessary allocations.
func splitOnComma(buffer []string, line string) []string {
	// Start with an empty buffer.  We just want the backing.
	buffer = buffer[:0]

	for {
		m := strings.IndexByte(line, ',')
		if m < 0 {
			break
		}
		buffer = append(buffer, line[:m])
		line = line[m+1:]
	}
	buffer = append(buffer, line)
	return buffer
}
