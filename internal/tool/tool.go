package tool

import (
	"os"
)

// File wraps an output, and handles various details with tool output, such
// as possibly directing output to standard out, and deleting an incomplete
// output file if there was an error.
type File struct {
	file  *os.File
	flags uint8
}

const (
	flagUnowned = 1
	flagKeep    = 2
)

// Open create a new Writer, directing output to standard out if the filename
// is "-".
func Open(filename string) (File, error) {
	if filename == "-" {
		return File{
			file:  os.Stdout,
			flags: flagUnowned,
		}, nil
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return File{}, err
	}

	return File{
		file: file,
	}, nil
}

// CombineErrors reports the first non-nil error in the list.
func CombineErrors(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	return err1
}

// Closes the underlying file, but only if the file is owned by the File.  Close
// will possibly remove the output file if there was an error (see Keep).
func (w *File) Close() error {
	if (w.flags & flagUnowned) != 0 {
		return nil
	}

	filename := w.file.Name()
	err := w.file.Close()

	if (w.flags & flagKeep) == 0 {
		err2 := os.Remove(filename)
		return CombineErrors(err, err2)
	}

	return err
}

// File returns the underlying file.
func (w *File) File() *os.File {
	return w.file
}

// Keep flags the output file as complete if err is nil.  If not called, or if
// err is not nil, the file will be remove when closed.
func (w *File) Keep(err error) {
	if err == nil {
		w.flags |= flagKeep
	}
}
