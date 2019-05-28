package main

import (
	"bytes"
	"compress/gzip"
	"strings"
	"testing"
)

func gzipString(t *testing.T, s string) string {
	buf := bytes.NewBuffer(nil)
	data := gzip.NewWriter(buf)
	_, err := data.Write([]byte(s))
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = data.Close()
	if err != nil {
		t.Fatalf("%s", err)
	}
	return buf.String()
}

func TestLoadGCovJSFile(t *testing.T) {
	cases := []struct {
		value string
		ok    bool
	}{
		{"empty", false},
		{gzipString(t, ""), false},
		{"{\"key\":123}", false},
		{gzipString(t, "{\"key\":123}"), true},
	}

	for _, v := range cases {
		t.Run(v.value, func(t *testing.T) {
			fds := make(FileDataSet)
			err := loadGCovJSFile(fds, strings.NewReader(v.value))
			if ok := err == nil; ok != v.ok {
				if err != nil {
					t.Logf("error: %s", err)
				}
				LogNE(t, "ok", v.ok, ok)
			}
		})
	}
}
