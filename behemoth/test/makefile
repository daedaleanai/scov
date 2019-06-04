CFLAGS ?= -O0 -I$(SRCDIR)/methods
LDFLAGS ?=
LLVM_PROFDATA ?= llvm-profdata-8 
LLVM_COV ?= llvm-cov-8 
LOC ?= 100000
FILES ?= 10
SCOV ?= ../../scov
ifeq ($(findstring clang,$(CC)),clang)
CFLAGS +=  -fprofile-instr-generate -fcoverage-mapping
LDFLAGS += -fprofile-instr-generate
else
CFLAGS += --coverage
LDFLAGS += --coverage
endif
SRC := rand$(LOC).c $(addsuffix .c,$(addprefix file,$(shell seq 1 ${FILES})))
OBJ := $(SRC:.c=.o)

SRCDIR ?= .
VPATH := $(SRCDIR)

rand$(LOC): $(OBJ)
	$(CC) $(LDFLAGS) -o $@ $^

$(SRC):
	../behemoth -loc $(LOC) -files $(FILES) > $@

html: rand$(LOC) $(SCOV)
	./rand$(LOC)
ifeq ($(findstring clang,$(CC)),clang)
	$(LLVM_PROFDATA) merge -o default.prof default.profraw
	$(LLVM_COV) export -instr-profile default.prof ./example > default.json
else
	time gcov -i -b *.c
endif
	time $(SCOV) -title "Example" -htmljs=true -srcid=`git describe` -testid="./rand" .

lcov: rand$(LOC)
	./rand$(LOC)
	-$(RM) *.gcov
	time lcov --capture --directory . --output-file rand$(LOC).info
	time genhtml rand$(LOC).info --output-directory .

clean:
	-$(RM) rand* file* *.gcno *.gcda
	-$(RM) -rf *.html *.js *.png *.css ./test/
