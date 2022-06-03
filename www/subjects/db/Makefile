
data.sqlite: makedb.go questions.csv
	rm -f $@
	go run -v $<
	chmod 664 $@

questions.csv: makecsv.go
	go run $<
