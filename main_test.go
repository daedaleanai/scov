package main

import (
	"bytes"
	"fmt"
	"testing"
)

func LogNE(t *testing.T, field string, expected, got interface{}) {
	t.Errorf("Values for %s not equal, expected %v, got %v", field, expected, got)
}

func TestHandleRequestFlags(t *testing.T) {
	cases := []struct {
		help    bool
		version bool
		ok      bool
	}{
		{false, false, false},
		{true, false, true},
		{false, true, true},
	}

	for i, v := range cases {
		buffer := bytes.NewBuffer(nil)
		ok := handleRequestFlags(buffer, v.help, v.version)

		if ok != (buffer.Len() > 0) {
			t.Errorf("Case %d: mismatch between output and return status", i)
		}
		if ok != v.ok {
			t.Errorf("Case %d: expected %v, got %v", i, v.ok, ok)
		}
	}
}

func TestFilterExternalFileData(t *testing.T) {
	const file1 = "binc.cpp"
	const file2 = "/usr/include/a.h"
	const file3 = "/usr/include/b.h"

	cases := []struct {
		in       []string
		external bool
		expected int
	}{
		{[]string{file1, file2, file3}, false, 1},
		{[]string{file1, file2, file3}, true, 3},
		{[]string{file1}, false, 1},
		{[]string{file1}, true, 1},
	}

	for i, v := range cases {
		name := fmt.Sprintf("Case %d", i)
		t.Run(name, func(t *testing.T) {
			data := make(map[string]*FileData)
			for _, v := range v.in {
				data[v] = NewFileData(v)
			}

			data = filterExternalFileData(data, v.external)
			if out := len(data); out != v.expected {
				LogNE(t, "file count", v.expected, out)
			}
		})
	}
}

func TestFilterExcludedFileData(t *testing.T) {
	const file1 = "binc.cpp"
	const file2 = "/usr/include/a.h"
	const file3 = "/usr/include/b.h"

	cases := []struct {
		in       []string
		filter   string
		ok       bool
		expected int
	}{
		{[]string{file1, file2, file3}, "", true, 3},
		{[]string{file1, file2, file3}, ".+", true, 0},
		{[]string{file1, file2, file3}, "^/", true, 1},
		{[]string{file1, file2, file3}, "\\.h$", true, 1},
		{[]string{file1, file2, file3}, "\\.cpp", true, 2},
		{[]string{file1, file2, file3}, "[]", false, 3},
	}

	for i, v := range cases {
		name := fmt.Sprintf("Case %d", i)
		t.Run(name, func(t *testing.T) {
			data := make(map[string]*FileData)
			for _, v := range v.in {
				data[v] = NewFileData(v)
			}

			out := bytes.NewBuffer(nil)
			data = filterExcludedFileData(out, data, v.filter)
			if ok := out.Len() == 0; ok != v.ok {
				LogNE(t, "ok", v.ok, ok)
			}
			if out := len(data); out != v.expected {
				LogNE(t, "file count", v.expected, out)
			}
		})
	}
}
