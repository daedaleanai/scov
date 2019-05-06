CFLAGS ?= -O0 -I$(SRCDIR)/methods
LDFLAGS ?=
LLVM_PROFDATA ?= llvm-profdata-8 
LLVM_COV ?= llvm-cov-8 
GCOVHTML ?= ../gcovhtml
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
	$(LLVM_COV) export -format lcov -instr-profile default.prof ./example > default.info
	$(GCOVHTML) -htmldir "" -text coverage.txt *.info
else
	gcov -i -b example.c methods/*.c
	$(GCOVHTML) -htmldir "" -text coverage.txt *.gcov
endif

html: example
	./example 0 20
ifeq ($(findstring clang,$(CC)),clang)
	$(LLVM_PROFDATA) merge -o default.prof default.profraw
	$(LLVM_COV) export -format lcov -instr-profile default.prof ./example > default.info
	$(GCOVHTML) -title "Example" -htmljs=true -srcid=`git describe` *.info
else
	gcov -i -b example.c methods/*.c
	$(GCOVHTML) -title "Example" -htmljs=true -srcid=`git describe` *.gcov
endif

clean:
	-$(RM) *.o methods/*.o example
	-$(RM) *.gcno *.gcda methods/*.gcno methods/*.gcda
	-$(RM) *.info *.prof *.profraw *.json.gz
	-$(RM) *.txt
	-$(RM) *.html *.gcov methods/*.html
