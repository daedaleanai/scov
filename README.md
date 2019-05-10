# SCov
> Generate reports on code coverage.

SCov collects code coverage data generated by instrumented binaries, using either [`gcov`](https://gcc.gnu.org/onlinedocs/gcc/Gcov.html) or [`llvm-cov`](http://llvm.org/docs/CommandGuide/llvm-cov.html), and then generates reports on the data.  There is a simple text report that calculates the line coverage and function coverage for all of the source files ([example](https://stone.code.gitlab.io/scov/example/coverage.txt)).  For more detailed information, there is an HTML report ([example](https://stone.code.gitlab.io/scov/example/)).  The HTML report includes line coverage, function coverage, and possibly branch coverage.  Annotated source files are also created ([example](https://stone.code.gitlab.io/scov/example/example.c.html)).

SCov is also quite a bit faster than `lcov` for generating reports.  Timing was measured for two sample projects, which had sizes of 0.3~kloc and 10~kloc.  It is only a few points, but the measurements show that SCov should be more than 30x faster.  For smaller code bases, `lcov` also has a significant start-up cost.

## Getting started

To generate a report on code coverage, you require an instrumented binary.  Details of this are linked to your compiler and linker, but examples for `gcc` and `clang` are included below.  Afterwards, running the instrumented binary will generate data files with the collected profiling data.  Typically, this involves running test code, but regular use works as well.  Finally, the profiling data can be exported, and processed by `scov` to create the reports.

### Using gcc

For recent versions of `gcc`, an instrumented binary can be built by adding the command-line flag `--coverage` to both compiling and linking.  After running the instrumented binary, run `gcov` to export the data to the intermediate text format.  Finally, you can run `scov` to create the reports.

```shell
gcc --coverage -g -O0 -o ./example [source files]
./example
gcov -i [source files]
scov -title "My Report" -htmldir ./html *.gcov
```

This will create a folder, and insert the HTML files into that folder.  Open `index.html` to get an overview of the code coverage, and follow the links for the annotated source files.

If you add the option `-b` when running `gcov`, the reports will include information on the branch coverage.

Note that `gcov` version 7 or higher is required.  Earlier versions do not support the `-i` command-line flag.

Starting with `gcov` version 9, the format of the output files has changed.  You will need to replace `*.gcov` with `*.json.gz`.

### Using clang

For recent versions of `clang`, adding the command-line flags `fprofile-instr-generate` and `-fcoverage-mapping` when compiling, and `-fprofile-instr-generate` when linking, will build an instrumented binary. After rnning the instrumented binary, you will need to process and extract the data using LLVM's tools.  Finally, you can run `scov` to create the reports.

```shell
clang -fprofile-instr-generate -fcoverage-mapping -g -O0 -o ./example [source files]
./example
llvm-profdata merge -o default.prof default.profraw
llvm-cov export -format lcov -instr-profile default.prof ./example > default.info
scov -title "My Report" -htmldir ./html default.info
```

This will create a folder, and insert the HTML files into that folder.  Open `index.html` to get an overview of the code coverage, and follow the links for the annotated source files.

Note that `llvm` version 8 or higher is required.  Earlier versions do not support the `-format` command-line flag when exporting data.

## Options

**-exclude [regexp]**  	Exclude source files that match the regular expression.

**-external**   Set whether external files to be included.

**-h**	Request help.

**-htmldir [folder]**  	Path for the HTML output (default ".").

**-htmljs**    	Use javascript to enhance reports.

**-srcdir [folder]**  	Path for the source directory (default ".").

**-srcid [string]**    	String to identify revision of the source.  This could be either `git describe` or `hg id`.

**-text [filename]**   	Filename for text report, use - to direct the report to stdout.

**-title string**    	Title for the HTML pages (default "SCov").

**-v**  Request version information.

## Installation

### From Source

To build a copy of `scov`, you will need a copy of [Go](https://golang.org/).  For instructions on how to install Go, please refer to the language's website.  No dependencies beyond the standard library are required, although `scov` does require version 1.10 or higher to run the  automated testing.

Once Go is installed, you can clone the repository and build the application.

```shell
go get -u gitlab.com/stone.code/scov
go install gitlab.com/stone.code/scov
```

### Binaries

Binaries can be downloaded for any [release](https://gitlab.com/stone.code/scov/releases).  Builds are included for multiple platforms.  If your platform is missing, please [open an issue](https://gitlab.com/stone.code/scov/issues).

Additionally, binaries are built with every commit, and they can be downloaded from the [pipelines pages](https://gitlab.com/stone.code/scov/pipelines).

## Contributing

Development of this project is ongoing.  If you find a bug or have any suggestions, please [open an issue](https://gitlab.com/stone.code/scov/issues).

If you'd like to contribute, please fork the repository and make changes.  Pull requests are welcome.

## Related projects

- [gcov](https://gcc.gnu.org/onlinedocs/gcc/Gcov.html):  Use the `gcov` tool in conjunction with GCC to test code coverage in your programs.
- [llvm-cov](http://llvm.org/docs/CommandGuide/llvm-cov.html): The `llvm-cov` tool shows code coverage information for programs that are instrumented to emit profile data. 
- [lcov](http://ltp.sourceforge.net/coverage/lcov.php):  LCOV is a graphical front-end for GCC's coverage testing tool `gcov`.
- [Kcov](https://simonkagstrom.github.io/kcov/):  Kcov is a code coverage tester for compiled programs, Python scripts and shell scripts.
- [Gcovr](https://pypi.org/project/gcovr/):  Gcovr provides a utility for managing the use of the GNU `gcov` utility and generating summarized code coverage results.
- [OpenCppCoverage](https://github.com/OpenCppCoverage/OpenCppCoverage):  OpenCppCoverage is an open source code coverage tool for C++ under Windows.

## Licensing

This project is licensed under the [3-Clause BSD License](https://opensource.org/licenses/BSD-3-Clause).  See the LICENSE in the repository.
