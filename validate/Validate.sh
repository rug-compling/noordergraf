#!/bin/bash

cd /net/noordergraf/validate

make

shaclvalidate.sh -shapesfile Shapefile -datafile Datafile | grep -v '^@prefix'
