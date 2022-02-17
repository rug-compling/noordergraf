all: Shapefile Datafile

Shapefile: ../data/prefix.ttl *.ttl
	echo '@prefix sh: <http://www.w3.org/ns/shacl#> .' > $@
	cat $^ >> $@

Datafile: ../data/*.ttl ../data/*/*.ttl
	/net/noordergraf/data/concat -n > $@
