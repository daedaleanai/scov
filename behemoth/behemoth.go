package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	loc       = flag.Int("loc", 1, "Target for the number of lines of code")
	seed      = flag.Int("seed", 0, "Seed for random number generator")
	files     = flag.Uint("files", 0, "Number of sub-files to create")
	funcCount = 0
)

func main() {
	flag.Parse()

	if *seed == 0 {
		rand.Seed(time.Now().UnixNano())
	}

	writeHeaders(os.Stdout)
	writeDummyFunc(os.Stdout)
	writeMain(os.Stdout)
	for i := uint(0); i < *files; i++ {
		writeSource(i + 1)
	}
}

func writeHeaders(out io.Writer) {
	fmt.Fprintf(out, "#include <stdlib.h>\n\n")
}

func writeDummyFunc(out io.Writer) {
	fmt.Fprintf(out, "void dummy( int p ) {\n")
	fmt.Fprintf(out, "\treturn;\n")
	fmt.Fprintf(out, "}\n\n")
}

func writeMain(out io.Writer) {
	if *files > 0 {
		for i := uint(0); i < *files; i++ {
			fmt.Fprintf(out, "extern int file%d( int p );\n", i+1)
		}
		_, _ = out.Write([]byte{'\n'})
	}

	fmt.Fprintf(out, "int main( int argc, char const* argv[] ) {\n")
	fmt.Fprintf(out, "\tint p = 0;\n")
	fmt.Fprintf(out, "\tif ( argc>1 ) {\n")
	fmt.Fprintf(out, "\t\tp = atoi( argv[1] );\n")
	fmt.Fprintf(out, "\t}\n\n")

	if *files > 0 {
		for i := uint(0); i < *files; i++ {
			fmt.Fprintf(out, "\tfile%d( p );\n", i+1)
		}
	} else {
		stmt := makeStatements(out, *loc)
		for _, v := range stmt {
			fmt.Fprintf(out, "\t%s\n", v)
		}
	}
	fmt.Fprintf(out, "\treturn 0;\n}\n\n")
}

func writeSource(ndx uint) {
	out, err := os.Create(fmt.Sprintf("file%d.c", ndx))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	fmt.Fprintf(out, "extern void dummy( int p );\n\n")

	stmt := makeStatements(out, *loc/int(*files))
	fmt.Fprintf(out, "int file%d( int p ) {\n", ndx)
	for _, v := range stmt {
		fmt.Fprintf(out, "\t%s\n", v)
	}
	fmt.Fprintf(out, "\treturn 0;\n}\n\n")
}

func makeStatements(out io.Writer, count int) []string {
	stmt := []string{}

	if count <= 0 {
		return nil
	}

	if count <= 10 {
		for i := 0; i < count; i++ {
			p := rand.Float32()
			if p < 0.8 {
				// dummy statement
				stmt = append(stmt, "dummy( p );")
			} else {
				// if statement
				buf := bytes.Buffer{}
				fmt.Fprintf(&buf, "if ( p != %d ) {\n", rand.Int()%10)
				buf.WriteString("\t\tdummy( p );\n")
				buf.WriteString("}")
				stmt = append(stmt, buf.String())
				i++
			}
		}
	} else {
		for i := 0; i < count; i++ {
			p := rand.Float32()
			if p < 0.7 {
				// dummy statement
				stmt = append(stmt, "dummy( p );")
			} else if p < 0.9 {
				// function call
				name := encodeName()
				stmt = append(stmt, name+"( p );")

				writeFunc(out, name, count/10)
				i += count/10 - 1
			} else {
				tmp := makeStatements(out, count/20)
				i += count/20 - 1
				buf := bytes.Buffer{}
				fmt.Fprintf(&buf, "if ( p != %d ) {\n", rand.Int()%10)
				for _, v := range tmp {
					fmt.Fprintf(&buf, "\t%s\n", v)
				}
				buf.WriteString("\t}")
				stmt = append(stmt, buf.String())
			}
		}
	}
	return stmt
}

func isReserved(name string) bool {
	return name == "p" ||
		name == "asm" || name == "do" || name == "if"
}

func encodeName() string {
	v := funcCount
	funcCount++

	buf := [20]byte{}
	s := buf[:0]
	for v >= 26 {
		s = append(s, 'a'+byte(v%26))
		v /= 26
	}
	s = append(s, 'a'+byte(v%26))

	name := string(s)
	return "bm_" + name
}

func writeFunc(out io.Writer, name string, count int) {
	stmt := makeStatements(out, count)

	fmt.Fprintf(out, "static int %s( int p ) {\n", name)
	for _, v := range stmt {
		fmt.Fprintf(out, "\t%s\n", v)
	}
	fmt.Fprintf(out, "}\n\n")
}
