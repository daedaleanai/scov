# GCovHTML
> Generate reports on code coverage using gcov.

GCovHTML collects code coverage using the [`gcov` tool](https://gcc.gnu.org/onlinedocs/gcc/Gcov.html), and then generates reports on the data.  There is a simple text report that calculates the line coverage and function coverage for all of the source files ([example](https://stone.code.gitlab.io/gcovhtml/example/coverage.txt)).  For more detailed information, there is an HTML report ([example](https://stone.code.gitlab.io/gcovhtml/example/)).  The HTML report includes line coverage, function coverage, and possibly branch coverage.  Annotated source files are also created ([example](https://stone.code.gitlab.io/gcovhtml/example/example.c.html)).

## Getting started

To generate a report on code coverage, you should already have an instrumented binary.  Details of this are linked to your compiler, but for recent versions of `gcc`, adding the command-line flag `--coverage` to both compiling and linking should work.  Running the binary will then generate data files with the collected profiling information.  Afterwards, you will need to run `gcov` to convert the data to the intermediate text format.  Finally, you can run `gcovhtml` to create the HTML reports.

```shell
gcc --coverage -g -O0 -o ./example [source files]
./example
gcov -i [source files]
gcovhtml -title "My Report" -htmldir ./html *.gcov
```

This will create a folder, and insert the HTML files into that folder.  Open `index.html` to get an overview of the code coverage, and follow the links for the annotated source files.

If you add the option `-b` when running `gcov`, the reports will include information on the branch coverage.

Note that `gcov` version 7 or higher is required.  Earlier versions do not support the `-i` command-line flag.

### Options

**-exclude [regexp]**  	Exclude source files that match the regular expression.

**-external**   Set whether external files to be included.

**-h**	Request help.

**-htmldir [folder]**  	Path for the HTML output (default ".").

**-srcdir [folder]**  	Path for the source directory (default ".").

**-srcid [string]**    	String to identify revision of the source.  This could be either `git describe` or `hg id`.

**-text [filename]**   	Filename for text report, use - to direct the report to stdout.

**-title string**    	Title for the HTML pages (default "GCovHTML").

**-v**  Request version information.

## Installation

### From Source

To build a copy of `gcovhtml`, you will need a copy of [Go](https://golang.org/).  For instructions on how to install Go, please refer to the language's website.  No dependencies beyond the standard library are required, although `gcovhtml` does require version 1.10 or higher to run the  automated testing.

Once Go is installed, you can clone the repository and build the application.

```shell
go get -u gitlab.com/stone.code/gcovhtml
go install gitlab.com/stone.code/gcovhtml
```

### Binaries

Binaries can be downloaded for any [release](https://gitlab.com/stone.code/gcovhtml/releases).  Builds are included for multiple platforms.  If your platform is missing, please [open an issue](https://gitlab.com/stone.code/gcovhtml/issues).

Additionally, binaries are built with every commit, and they can be downloaded from the [pipelines pages](https://gitlab.com/stone.code/gcovhtml/pipelines).

## Contributing

Development of this project is ongoing.  If you find a bug or have any suggestions, please [open an issue](https://gitlab.com/stone.code/gcovhtml/issues).

If you'd like to contribute, please fork the repository and make changes.  Pull requests are welcome.

## Related projects

- [gcov](https://gcc.gnu.org/onlinedocs/gcc/Gcov.html):  Use the `gcov` tool in conjunction with GCC to test code coverage in your programs.
- [llvm-cov](http://llvm.org/docs/CommandGuide/llvm-cov.html): The `llvm-cov` tool shows code coverage information for programs that are instrumented to emit profile data. 
- [lcov](http://ltp.sourceforge.net/coverage/lcov.php):  LCOV is a graphical front-end for GCC's coverage testing tool `gcov`.
- [Gcovr](https://pypi.org/project/gcovr/):  Gcovr provides a utility for managing the use of the GNU `gcov` utility and generating summarized code coverage results.

## Licensing

This project is licensed under the [3-Clause BSD License](https://opensource.org/licenses/BSD-3-Clause).  See the LICENSE in the repository.
