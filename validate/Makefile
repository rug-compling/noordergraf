all: Shapefile Datafile

SRC = ../data/prefix.ttl \
	*.ttl

Shapefile: $(SRC)
	@echo '@prefix sh: <http://www.w3.org/ns/shacl#> .' > $@
	@cat $^ >> $@

Datafile: ../data/prefix.ttl ../data/ns.ttl ../data/*/*.ttl
	@cat $^ > $@
