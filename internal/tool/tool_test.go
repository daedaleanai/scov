package tool_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"gitlab.com/stone.code/scov/internal/tool"
)

func TempFilename(t *testing.T) (filename string, closer func()) {
	file, err := ioutil.TempFile("", "testing")
	if err != nil {
		t.Fatalf("could not create temporary file: %s", err)
		// unreachable
	}
	name := file.Name()
	file.Close()

	return name, func() {
		err := os.Remove(name)
		if err != nil {
			t.Logf("could not remove temporary file: %s", err)
		}
	}
}

func TestOpen(t *testing.T) {
	tmpfile, closer := TempFilename(t)
	defer closer()

	cases := []struct {
		filename string
		isStdout bool
		keep     bool
	}{
		{"-", true, false},
		{tmpfile, false, false},
		{tmpfile, false, true},
	}

	for _, v := range cases {
		t.Run(v.filename, func(t *testing.T) {
			w, err := tool.Open(v.filename)
			if err != nil {
				t.Fatalf("Could not open file, %s", err)
			}

			if (w.File() == os.Stdout) != v.isStdout {
				t.Errorf("Mismatch on is stdout")
			}

			if v.keep {
				w.Keep(nil)
			}

			err = w.Close()
			if err != nil {
				t.Errorf("Could not close file, %s", err)
			}

			_, err = os.Stat(v.filename)
			if (err == nil) != v.keep {
				t.Errorf("Failed to cleanup file, %s", err)
			}
		})
	}
}

func TestOpenFail(t *testing.T) {
	_, err := tool.Open(".")
	if err == nil {
		t.Errorf("Unexpected success")
	}
}

func TestCombineErrors(t *testing.T) {
	mock1 := errors.New("mock1")
	mock2 := errors.New("mock2")

	cases := []struct{ err1, err2, expected error }{
		{nil, nil, nil},
		{mock1, nil, mock1},
		{nil, mock2, mock2},
		{mock1, mock2, mock1},
	}

	for i, v := range cases {
		out := tool.CombineErrors(v.err1, v.err2)
		if out != v.expected {
			t.Errorf("Case %d: want %v, got %v", i, v.expected, out)
		}
	}
}
