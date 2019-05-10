package main

import (
	"testing"
)

func TestDARecord(t *testing.T) {
	cases := []struct {
		value    string
		ok       bool
		lineNo   int
		hitCount uint64
	}{
		{"DA:38,3", true, 38, 3},
		{"DA:38", false, 0, 0},
		{"DA:#,3", false, 0, 0},
		{"DA:38,#", false, 0, 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "DA" {
				LogNE(t, "record type", "DA", rt)
			}
			lineNo, hc, err := parseDARecord(value)
			if lineNo != v.lineNo {
				LogNE(t, "line number", v.lineNo, lineNo)
			}
			if hc != v.hitCount {
				LogNE(t, "hit count", v.hitCount, hc)
			}
			if (err == nil) != v.ok {
				LogNE(t, "ok", v.ok, err == nil)
				if err != nil {
					t.Logf("err = %s", err)
				}
			}
		})
	}
}

func TestFNDARecord(t *testing.T) {
	cases := []struct {
		value    string
		ok       bool
		funcName string
		hitCount uint64
	}{
		{"FNDA:3,gauss_get_sum", true, "gauss_get_sum", 3},
		{"FNDA:3", false, "", 0},
		{"FNDA:#,gauss_get_sum", false, "", 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "FNDA" {
				LogNE(t, "record type", "FNDA", rt)
			}
			funcName, hc, err := parseFNDARecord(value)
			if funcName != v.funcName {
				LogNE(t, "function name", v.funcName, funcName)
			}
			if hc != v.hitCount {
				LogNE(t, "hit count", v.hitCount, hc)
			}
			if (err == nil) != v.ok {
				LogNE(t, "ok", v.ok, err == nil)
				if err != nil {
					t.Logf("err = %s", err)
				}
			}
		})
	}
}

func TestBRDARecord(t *testing.T) {
	cases := []struct {
		value  string
		ok     bool
		lineNo int
		status BranchStatus
	}{
		{"BRDA:42,0,0,0", true, 42, BranchNotTaken},
		{"BRDA:42,0,1,3", true, 42, BranchTaken},
		{"BRDA:42,0,1", false, 0, 0},
		{"BRDA:#,0,0,0", false, 0, 0},
		{"BRDA:42,0,0,#", false, 0, 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "BRDA" {
				LogNE(t, "record type", "BRDA", rt)
			}
			lineNo, status, err := parseBRDARecord(value)
			if lineNo != v.lineNo {
				LogNE(t, "line number", v.lineNo, lineNo)
			}
			if status != v.status {
				LogNE(t, "branch status", v.status, status)
			}
			if (err == nil) != v.ok {
				LogNE(t, "ok", v.ok, err == nil)
				if err != nil {
					t.Logf("err = %s", err)
				}
			}
		})
	}
}
