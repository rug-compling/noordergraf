Invoer data in AllegroGraph waarbij alle bestaande data wordt gewist:

```
cat prefix.ttl */*.ttl | agtool load --overwrite --bulk --duplicates delete --fti all --automate-nd-datatype-mappings --input turtle --optimize --base-uri https://noordergraf.rug.nl/ noordergraf -
```

of:

```
(
    cat prefix.ttl
    echo '<https://noordergraf.rug.nl/tomb/> {'
    cat tomb/*.ttl
    echo '}'
    echo '<https://noordergraf.rug.nl/site/> {'
    cat site/*.ttl
    echo '}'
    echo '<https://noordergraf.rug.nl/place/> {'
    cat place/*.ttl
    echo '}'
) | agtool load --overwrite --bulk --duplicates delete --fti all --automate-nd-datatype-mappings --input trig --optimize --base-uri https://noordergraf.rug.nl/ noordergraf -
```
