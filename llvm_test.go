package main

import (
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
