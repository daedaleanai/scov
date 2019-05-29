package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLoadLLVMFile(t *testing.T) {
	cases := []struct {
		value string
		ok    bool
	}{
		{"empty", false},
		{"{\"key\":123}", false},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			fds := make(FileDataSet)
			err := loadLLVMFile(fds, strings.NewReader(v.value))
			if ok := err == nil; ok != v.ok {
				if err != nil {
					t.Logf("error: %s", err)
				}
				LogNE(t, "ok", v.ok, ok)
			}
		})
	}
}

func TestLLVMSegmentUnmarshalJSON(t *testing.T) {
	cases := []struct {
		in       string
		expected *LLVMSegment
	}{
		{"[10,1,1,true,true]", &LLVMSegment{10, 1, 1, true, true}},
		{"[10,2,3,false,true]", &LLVMSegment{10, 2, 3, false, true}},
		{"[10,2,3,false,true,{}]", &LLVMSegment{10, 2, 3, false, true}}, // excess elements
		{"[{},2,3,false,true]", nil},
		{"[10,{},3,false,true]", nil},
		{"[10,2,{},false,true]", nil},
		{"[10,2,3,{},true]", nil},
		{"[10,2,3,true,{}]", nil},
		{"[10,2,3,false,]", nil},
		{"[10,2,3,false]", nil},
		{"[10,2,3]", nil},
		{"[10,2]", nil},
		{"[10]", nil},
		{"[]", nil},
		{"10", nil},
	}

	for _, v := range cases {
		t.Run(v.in, func(t *testing.T) {
			out := LLVMSegment{}
			err := json.Unmarshal([]byte(v.in), &out)
			if got := err == nil; got != (v.expected != nil) {
				if err != nil {
					t.Logf("err = %s", err)
				}
				t.Errorf("ok: expected %v, got %v", (v.expected != nil), got)
			}
			if err == nil {
				if out != *v.expected {
					t.Errorf("expected %v, got %v", v.expected, out)
				}
			}
		})
	}
}
