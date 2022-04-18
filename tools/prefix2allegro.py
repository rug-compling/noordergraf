#!/usr/bin/env python3

# Dit programma set /net/noordergraf/data/prefix.ttl om in het formaat
# gebruikt door de server van AllegroGraph
#
# Dit programma wordt aangeroepen vanuit /net/noordergraf/data/input

import sys

data = '''
@prefix dc:      <http://purl.org/dc/elements/1.1/> .
@prefix dcterms: <http://purl.org/dc/terms/> .
@prefix err:     <http://www.w3.org/2005/xqt-errors#> .
@prefix fn:      <http://www.w3.org/2005/xpath-functions#> .
@prefix foaf:    <http://xmlns.com/foaf/0.1/> .
@prefix fti:     <http://franz.com/ns/allegrograph/2.2/textindex/> .
@prefix keyword: <http://franz.com/ns/keyword#> .
@prefix nd:      <http://franz.com/ns/allegrograph/5.0/geo/nd#> .
@prefix ndfn:    <http://franz.com/ns/allegrograph/5.0/geo/nd/fn#> .
@prefix owl:     <http://www.w3.org/2002/07/owl#> .
@prefix rdf:     <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs:    <http://www.w3.org/2000/01/rdf-schema#> .
@prefix skos:    <http://www.w3.org/2004/02/skos/core#> .
@prefix xs:      <http://www.w3.org/2001/XMLSchema#> .
@prefix xsd:     <http://www.w3.org/2001/XMLSchema#> .)
'''.split('\n')

with open('/net/noordergraf/data/prefix.ttl', 'rt') as fp:
    for line in fp:
        data.append(line)

lines = []
seen = {}
for line in data:
    a = line.split()
    if len(a) != 4:
        continue
    if a[0] != '@prefix':
        continue
    if a[1] == 'll:':
        continue
    if a[1] in seen:
        continue
    seen[a[1]] = True
    key = a[2].split('/')[2] + a[1]
    lines.append('{}\t("{}" "{}")\n'.format(key, a[1][:-1], a[2][1:-1]))

lines.sort()

sys.stdout.write('(\n')
for line in lines:
    sys.stdout.write(line.split('\t')[1])
sys.stdout.write(')\n')
