CFLAGS ?= --coverage -O0 -I$(SRCDIR)/methods
LDFLAGS ?= --coverage
GCOVHTML ?= ../gcovhtml

SRCDIR ?= .
VPATH := $(SRCDIR)

example: example.o methods/gauss.o methods/iterate.o

text: example
	./example 0 20
	gcov -i -b example.c methods/*.c
	$(GCOVHTML) -htmldir "" -text coverage.txt *.gcov

html: example
	./example 0 20
	gcov -i -b example.c methods/*.c
	$(GCOVHTML) -title "Example" -srcid=`git describe` *.gcov

clean:
	-$(RM) *.o methods/*.o example
	-$(RM) *.gcno *.gcda methods/*.gcno methods/*.gcda
	-$(RM) *.txt
	-$(RM) *.html *.gcov methods/*.html
