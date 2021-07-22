Invoer data in AllegroGraph waarbij alle bestaande data wordt gewist:

```
cat prefix.ttl ns.ttl */*.ttl | agtool load --overwrite --bulk --duplicates delete --fti all --automate-nd-datatype-mappings --input turtle --optimize --base-uri https://noordergraf.rug.nl/ noordergraf -
```
