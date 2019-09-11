CFLAGS ?= -O0 -I$(SRCDIR)/methods
LDFLAGS ?=
LLVM_PROFDATA ?= llvm-profdata-8 
LLVM_COV ?= llvm-cov-8 
SCOV ?= ../scov
ifeq ($(findstring clang,$(CC)),clang)
CFLAGS +=  -fprofile-instr-generate -fcoverage-mapping
LDFLAGS += -fprofile-instr-generate
else
CFLAGS += --coverage
LDFLAGS += --coverage
endif

SRCDIR ?= .
VPATH := $(SRCDIR)

example: example.o methods/gauss.o methods/iterate.o

text: example
	./example 0 20
ifeq ($(findstring clang,$(CC)),clang)
	$(LLVM_PROFDATA) merge -o default.prof default.profraw
	$(LLVM_COV) export -instr-profile default.prof ./example > default.json
else
	gcov -i -b example.c methods/*.c
endif
	$(SCOV) -htmldir "" -text coverage.txt .

html: example
	./example 0 20
ifeq ($(findstring clang,$(CC)),clang)
	$(LLVM_PROFDATA) merge -o default.prof default.profraw
	$(LLVM_COV) export -instr-profile default.prof ./example > default.json
else
	gcov -i -b example.c methods/*.c
endif
	$(SCOV) -title "Example" -htmljs=true -srcid=`git describe` -testid="./example 0 20" .

markdown: example
	./example 0 20
ifeq ($(findstring clang,$(CC)),clang)
	$(LLVM_PROFDATA) merge -o default.prof default.profraw
	$(LLVM_COV) export -instr-profile default.prof ./example > default.json
else
	gcov -i -b example.c methods/*.c
endif
	$(SCOV) -title "Example" -htmldir "" -markdown coverage.md -srcid=`git describe` -testid="./example 0 20" .

lcov: example
	./example 0 20
	-$(RM) *.gcov
	lcov --capture --directory . --output-file coverage.info
	genhtml coverage.info --output-directory ./lcov-html

clean:
	-$(RM) *.o methods/*.o example
	-$(RM) *.gcno *.gcda methods/*.gcno methods/*.gcda
	-$(RM) *.gcov *.info *.prof *.profraw *.json.gz
	-$(RM) *.txt *.md
	-$(RM) -rf *.html *.js methods/*.html lcov-html
