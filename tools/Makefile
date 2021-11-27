
% : %.go
	go build $<

all: jpg2ttl name2soundex soundex-test xml2ttl

name2soundex: name2soundex.go ../go/nlsoundex/nlsoundex.go

soundex-test: soundex-test.go ../go/nlsoundex/nlsoundex.go
