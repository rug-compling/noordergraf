## Hulpmiddelen voor de invoer van data

Hier staan een aantal programma's die van pas
komen bij het bewerken van data.

 - `jpg2ttl` zet (exif-data in) een jpeg-bestand om naar een
   Turtle-fragment, dat je kunt opnemen in een Turtle-bestand voor een
   grafsteen.
 - `penman2turtle` zet een Penman-bestand om in een Turtle-bestand.
 - `ttl2json`, `ttl2jsont`, `ttl2rdf`, `ttl2tri`, `ttl2xml` zijn scripts die
   bestanden in Turtle-formaat omzetten in diverse andere formaten.
 - `xml2ttl` zet bestanden in het formaat `rdfxml-abbrev` om naar Turtle.

Het is lastig om automatisch bulkbewerkingen uit te voeren op
Turtle-bestanden. De aanpak is dan:

 1. Omzetten van Turtle naar `rdfxml-abbrev` met `ttl2xml`.
 2. Bewerking van XML-bestanden met behulp van XPath.
 3. Terug omzetten naar Turtle met `xml2ttl`.

