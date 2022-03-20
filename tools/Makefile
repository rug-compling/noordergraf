
% : %.go
	go build $<

all: chkll chkplace chktags jpg2ttl mkprops mksymidx name2soundex soundex-test xml2ttl

name2soundex: name2soundex.go ../go/nlsoundex/nlsoundex.go

soundex-test: soundex-test.go ../go/nlsoundex/nlsoundex.go
