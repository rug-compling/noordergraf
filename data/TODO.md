
## Datatypes

Datum zonder/met tijd

```
"2014-02-23"^^xsd:date
"2014-02-23T03:18:59"^^xsd:dateTime
```

Tekst met taalcode

```
"Ljouwert"@fy
```

## Plaatsen in plaatsen

Wat te doen met geneste geo-informatie? Voorbeeld (penman):

```
:pob (c2 / village.n.02
     :nam "OOSTWOLD"
     :geo "2749325"
     :mem-of (c3 / municipality.n.01
                 :nam "OLDAMBT"
                 :geo "7535498"))
```

...zo [1] ? (turtle)

```
item:t00000.p1 :pob place:Oostwold2749325 ,
                    place:Oldambt7535498 .
```

... of zo [2] ? (turtle)

```
item:t00000.p1 :pob [
    :primary place:Oostwold2749325  ;
    :secondary place:Oldambt7535498
] .
```

... of zo [3] (turtle) :

```
item:t00000.p1 :pob item:t00000.c2 .

item:t00000.c2 a t:village.n.02 ;
  :place place:Oostwold2749325 ;
  :nam "OOSTWOLD" .

item:t00000.c3 a t:municipality.n.01 ;
  :place place:Oldambt7535498 ;
  :mem item:t00000.c2 ;
  :nam "OLDAMBT" .
```

Belangrijk is om informatie steeds op dezelfde diepte te hebben, dus
optie [2] valt af.
