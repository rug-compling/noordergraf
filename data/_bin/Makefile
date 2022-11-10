
% : %.go
	go build $<

all: chkdirs chkll chklocation chktags mkbibidx mkprops mksymidx name2soundex soundex-test

name2soundex: name2soundex.go ../../go/nlsoundex/nlsoundex.go

soundex-test: soundex-test.go ../../go/nlsoundex/nlsoundex.go
