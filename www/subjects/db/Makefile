
include ../bin/Makefile.cfg

data.sqlite: makedb.go questions.csv
	rm -f data.sqlite
	$(GO) fmt .
	$(GO) run -v $<
