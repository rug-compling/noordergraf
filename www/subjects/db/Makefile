
all:
	@echo
	@echo Gebruik: make new
	@echo
	@echo WAARSCHUWING: dit zal alle al verzamelde data wissen
	@echo

new: data.sqlite

data.sqlite: makedb.go questions.csv
	rm -f $@
	go run -v $<
	chmod 664 $@

questions.csv: makecsv.go
	go run $<
