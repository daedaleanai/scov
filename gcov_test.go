package main

import (
	"testing"
)

func TestParseFunctionRecord(t *testing.T) {
	cases := []struct {
		value    string
		ok       bool
		name     string
		line     int
		hitCount uint64
	}{
		{"function:222,6,main", true, "main", 222, 6},
		{"function:222,223,6,main", true, "main", 222, 6},
		{"function:222,0,main", true, "main", 222, 0},
		{"function:222,223,0,main", true, "main", 222, 0},
		{"function:#,6,main", false, "", 0, 0},
		{"function:222,#,main", false, "", 0, 0},
		{"function:#,223,6,main", false, "", 0, 0},
		{"function:222,223,#,main", false, "", 0, 0},
		{"function:222,main", false, "", 0, 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "function" {
				LogNE(t, "record type", "function", rt)
			}
			name, line, hc, err := parseFunctionRecord(value)
			if name != v.name {
				LogNE(t, "function name", v.name, name)
			}
			if line != v.line {
				LogNE(t, "function line", v.name, name)
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

func TestParseLCountRecord(t *testing.T) {
	cases := []struct {
		value    string
		ok       bool
		lineNo   int
		hitCount uint64
	}{
		{"lcount:32,18", true, 32, 18},
		{"lcount:32,0", true, 32, 0},
		{"lcount:32,18,0", true, 32, 18},
		{"lcount:32,0,1", true, 32, 0},
		{"lcount:32,#", false, 0, 0},
		{"lcount:#,18", false, 0, 0},
		{"lcount:32", false, 0, 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "lcount" {
				LogNE(t, "record type", "lcount", rt)
			}
			lineNo, hc, err := parseLCountRecord(value)
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

func TestParseBranchRecord(t *testing.T) {
	cases := []struct {
		value  string
		ok     bool
		lineNo int
		status BranchStatus
	}{
		{"branch:176,taken", true, 176, BranchTaken},
		{"branch:176,nottaken", true, 176, BranchNotTaken},
		{"branch:176,notexec", true, 176, BranchNotExec},
		{"branch:176,tkn", false, 0, 0},
		{"branch:176", false, 0, 0},
		{"branch:#,taken", false, 0, 0},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			rt, value := recordType(v.value)
			if rt != "branch" {
				LogNE(t, "record type", "lcount", rt)
			}
			lineNo, status, err := parseBranchRecord(value)
			if lineNo != v.lineNo {
				LogNE(t, "function name", v.lineNo, lineNo)
			}
			if status != v.status {
				LogNE(t, "hit count", v.status, status)
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
