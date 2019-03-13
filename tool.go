package main

import (
	"os"
)

// Writer wraps an output, and handles various details with tool output, such
// as possibily directing output to standard out, and deleting an incomplete
// output file if there was an error.
type Writer struct {
	file  *os.File
	flags uint8
}

const (
	flag_unowned = 1
	flag_keep    = 2
)

// Open create a new Writer, directing output to standard out if the filename
// is "-".
func Open(filename string) (Writer, error) {
	if filename == "-" {
		return Writer{
			file:  os.Stdout,
			flags: flag_unowned,
		}, nil
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return Writer{}, err
	}

	return Writer{
		file: file,
	}, nil
}

// Close the writer, but only if the file is owned by the writer.  Close will
// possibly remove the output file if there was an error (see Keep).
func (w *Writer) Close() error {
	if (w.flags & flag_unowned) != 0 {
		return nil
	}

	filename := w.file.Name()
	err := w.file.Close()

	if (w.flags & flag_keep) == 0 {
		err2 := os.Remove(filename)
		if err2 != nil && err != nil {
			err = err2
		}
	}

	return err
}

// File returns the underlying file.
func (w *Writer) File() *os.File {
	return w.file
}

// Keep flags the output file as complete if err is nil.  If not called, or if
// err is not nil, the file will be remove when closed.
func (w *Writer) Keep(err error) error {
	if err == nil {
		w.flags |= flag_keep
	}
	return err
}
