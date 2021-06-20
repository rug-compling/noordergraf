
## Datatypes

Datum zonder/met tijd

    "2014-02-23"^^xsd:date
    "2014-02-23T03:18:59"^^xsd:dateTime

Coördinaten

    "+53.2192+6.5667"^^<[spherical type]>

Tekst met taalcode

    "Ljouwert"@fy

## Plaatsen in plaatsen

Wat te doen met geneste geo-informatie? Voorbeeld (penman):

    :pob (c2 / village.n.02
         :nam "OOSTWOLD"
         :geo "2749325"
         :mem-of (c3 / municipality.n.01
                     :nam "OLDAMBT"
                     :geo "7535498"))

...zo? (turtle)

    item:t00000.p1 :pob place:Oostwold2749325 ,
                        place:Oldambt7535498 .

... of zo? (turtle)

    item:t00000.p1 :pob [
        :primary place:Oostwold2749325  ;
        :secondary place:Oldambt7535498
    ] .

