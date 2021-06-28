Invoer data in AllegroGraph waarbij alle bestaande data wordt gewist:

```
#agtool load --overwrite --optimize --bulk --graph :source --duplicates delete --fti all --automate-nd-datatype-mappings noordergraf */*.ttl
agtool load --overwrite --bulk --duplicates delete --fti all --automate-nd-datatype-mappings noordergraf */*.ttl
```

----

Voorbeeld van gedeeltelijke update

Eerst data uit oude bestand verwijderen:

```
echo 'CLEAR GRAPH <file://sites/sites.ttl>' | agtool query noordergraf -
```

Daarna bijgewerkt bestand opnieuw invoeren:

```
agtool load --optimize --bulk --graph :source --duplicates delete --fti all --automate-nd-datatype-mappings noordergraf sites/sites.ttl
```

----
