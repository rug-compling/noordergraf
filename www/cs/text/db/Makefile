
all:
	@echo
	@echo Gebruik:
	@echo '    make new'
	@echo '    make update'
	@echo

new:
	go run makecsv.go
	go run makedb.go
	chmod 664 $@

update:
	go run makecsv.go
	go run updatedb.go

