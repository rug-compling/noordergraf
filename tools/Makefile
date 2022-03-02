
% : %.go
	go build $<

all: chkll chktags jpg2ttl mksymidx name2soundex soundex-test xml2ttl

name2soundex: name2soundex.go ../go/nlsoundex/nlsoundex.go

soundex-test: soundex-test.go ../go/nlsoundex/nlsoundex.go
