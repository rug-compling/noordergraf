#!/bin/bash

cd /net/noordergraf/validate

make Datafile

echo '@prefix sh: <http://www.w3.org/ns/shacl#> .' > Partshapefile
cat ../data/prefix.ttl "$@" >> Partshapefile

shaclvalidate.sh -shapesfile Partshapefile -datafile Datafile | grep -v '^@prefix'
