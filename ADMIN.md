# Handleiding noordergraf.rug.nl

Dit bestand geeft een overzicht van de onderdelen van de server voor
de website Noordergraf.

Voor een lijst van directory's en bestanden, zie `MANIFEST.txt`

## Configuratie van Apache2


De site *noordergraf.rug.nl* wordt gedefinieerd in het bestand
`/etc/apache2/sites-available/noordergraf.conf`. Dit bestand is
gesymlinkt naar `/etc/apache2/sites-enabled/`. In die laatste directory
staan geen andere configuraties gelinkt. `noordergraf.conf` definieert
de servername, de serverroot `/net/noordergraf/www/`, gebruik van https,
en nog een paar dingen.

In `/etc/apache2/envvars` staat `APACHE_RUN_GROUP=noordergraf`, wat
ervoor zorgt de Apache dezelfde lees- en schrijfrechten heeft als de
gebruikers in groep *noordergraf*.

In `/etc/apache2/mods-available/dir.conf` wordt als enige waarde voor
`DirectoryIndex` het bestand `index` gedefinieerd. Omdat er Multiviews
gebruikt worden (zie beneden) wordt door Apache automatisch het
index-bestand met de juiste extensie gekozen.

In `/etc/apache2/mods-available/ssl.conf` is `SSLProtocol TLSv1.2`
gedefinieerd. Geen idee waarom.

In `/etc/apache2/conf-available/noordergraf.conf` gesymlinkt naar
`/etc/apache2/conf-enabled/`, zijn een aantal opties voor de website
gezet. In het bijzonder de volgende dingen:

 - Voor `ns.ttl` en directory's met Turtle-bestanden wordt de handler
   `/net/noordergraf/handlers/data/data` gedefinieerd. Requests voor
   deze bestanden worden dus niet direct door Apache2 zelf
   afgehandeld, maar eerst verwerkt door de handler. Komen er
   directory's bij met Turtle-bestanden, dan moet `noordergraf.conf` (in
   `conf-available/`) aangepast worden, en mogelijk ook de handler zelf.
 - Voor bestanden met de extensie `.md` wordt de handler `/bin/markdown`
   in de DocumentRoot gedefinieert. Deze handler zet
   markdown-bestanden om in HTML. \
   TODO: Dit is makkelijk, maar niet het meest efficiënt. Iedere keer
   dat een markdown-bestand wordt opgevraagd wordt het opnieuw omgezet
   naar HTML, terwijl het markdown-bestand zelden verandert. Dus
   misschien moet dit veranderd worden. \
   Het is wel handig als we de wiki van github lokaal willen hosten,
   omdat die helemaal uit markdown bestaat.


## Configuratie van AllegroGraph

De sources voor [AllegroGraph](https://allegrograph.com/) zijn geïnstalleerd in
`/opt/agraph-7.3.0/`. AllegroGraph draait vanuit
`/opt/agraph/`. Installatie van AllegroGraph gaat met `install-agraph
/opt/agraph` vanuit `/opt/agraph-7.3.0/`. Dit moet als gebruiker
`agraph`. Het `agraph`-account is aangemaakt als een systeem-account
(laag UID), maar wel met een home-directory en een login-shell,
anders werkt `agtool` niet, om een duistere reden.

Bij de installatie wordt het configuratiebestand
`/opt/agraph/lib/agraph.cfg` aangemaakt, dat later met de hand is
aangepast. Let daar op als je een nieuwere versie van AllegroGraph
installeert.

Namespaces zijn hier gedefinieerd:

 - `/opt/agraph/data/settings/default-namespaces`

Dit zijn de standaard namespaces, een aantal specifiek voor de
werking van AllegroGraph, en namespaces gedefinieerd voor Noordergraf.
Door deze definities kun je query's uitvoeren zonder zelf de
namespaces te definiëren in je query.

De namespaces worden bijgewerkt vanuit het script
`/net/noordergraf/data/input` (dat de data in AllegroGraph opnieuw
invoert), door aanroep van het script
`/net/noordergraf/data/_bin/prefix2allegro.py`. Dit laatste script voegt de
standaard en Allegro-specifieke namespaces samen met de namespaces voor
Noordergraf die gedefinieerd zijn in `/net/noordergraf/data/prefix.ttl`.

De poort 10036 is van buiten af bereikbaar, en wordt gebruikt voor
zowel de web-interface als de API. Dit gaat over https.

Op de machine noordergraf zelf kun je ook poort 10035 gebruiken. Dit
gebruikt plain http.

AllegroGraph wordt na het booten automatisch gestart mbv
`/etc/systemd/system/agraph.service`, dat een symlink is naar
`/opt/etc/systemd/agraph.service`.

`/opt/etc/agraph/secret.sh` bevat de login-gegevens van het
admin-account in AllegroGraph.

`/opt/etc/ssl/agraph.pem` bevat de certificaten die nodig zijn voor
het gebruik van https door AllegroGraph. Dit moet bijgwerkt worden als
het certificaat voor de website is bijgewerkt. Zie:
`/etc/ssl/private/key-noordergraf.rug.nl.README`


## Hulpmiddelen voor de invoer van data

In `/net/noordergraf/tools/` staan een aantal programma's die van pas
komen bij het bewerken van data.

TODO: documenteren in `/net/noordergraf/tools/` .

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

## Meer hulpmiddelen

Het programma [penman](https://penman.readthedocs.io/) staat in `/usr/local/bin/`.

De
[AllegroGraph Python client](https://franz.com/agraph/support/documentation/current/python/)
voor Python 3 is geïnstalleerd.


## Corrigeren van bulk data

In `/net/noordergraf/www/cs/` staan een paar subdirectory's voor de
handmatige bewerking van enkele aspecten van alle data. Dit zijn
voorbeelden. Dit is gemodelleerd op het voorbeeld in
https://www.let.rug.nl/kleiweg/crowdsourcing/ waar ook gedocumenteerd
is hoe het systeem in elkaar zit.

Dit is alleen voor het verzamelen van bewerkingen door gebruikers. De
toepassing van de verzamelde gegevens is nog niet geïmplementeerd, en
zal per situatie verschillen.


## Controle van data

Wanneer je nieuwe data hebt ingevoerd, of oude data bewerkt, is het
goed te controlleren of de data voldoet aan het formaat dat we hebben
vastgelegd voor Noordergraf. Dat kan door in de directory
`/net/noordergraf/validate/` het script `Validate.sh` te draaien.

Dit maakt gebruik van de
[TopBraid SHACL API](https://github.com/TopQuadrant/shacl) die
geïnstalleerd is onder `/opt/shacl`.

TODO: Het script controleert alle data, en is dus niet heel snel. Er
moet een script komen dat je gebruikt om één of enkele
Turtle-bestanden te controlleren.

Wanneer de data wordt ingevoerd in AllegroGraph (zie beneden) worden
nog andere controles uitgevoerd.


## Formatteren van data

TODO: Automatisch formatteren van Turtle zou kunnen door vertaling
naar `rdfxml-abbrev` met `ttl2xml` en terug naar Turtle met `xml2ttl`.


## Invoer van data in AllegroGraph

Zie voor procedure: `/net/noordergraf/data/README.md`

In `/net/noordergraf/data/_bin/` staan een aantal hulpprogramma's die
gebruikt worden bij het invoeren van data in AllegroGraph door
`/net/noordergraf/data/input`. Sommige van deze tools voeren controles
uit, andere doen voorbewerkingen zoals het maken van lijsten met
gegevens die gebruikt worden door de handler van de webserver. Bovenin
de broncode van elke tool staat waar het voor dient.

TODO:

Een van de twee scripts, `concat`, die gebruikt worden voor het invoeren van data,
zal niet functioneren als er veel data bij komt, omdat er dan te lange
command lines in het script worden gebruikt. Dat script zal dus
vervangen moeten worden door iets anders.

TODO:

Nieuwe data wordt ingevoerd door eerst alles in AllegroGraph te
wissen, en daarna alle data opnieuw in te voeren. Dit gaat nog
redelijk snel, maar is niet erg efficiënt, zeker niet als er veel meer
data is dan nu.

Er zou een script moeten komen waarmee je één Turtle-bestand kunt
bijwerken of toevoegen. Dat kan omdat de inhoud van Turtle-bestanden
niet als triples maar als quads worden ingevoerd. De vierde
waarde is de naam van het bestand, zonder extensie. Je kunt hierop
filteren, en dus alleen de data uit één bestand wissen uit
AllegroGraph voordat je dit ene bestand (opnieuw) invoert.


## Backup van data in git

Het grootste deel van alles onder `/net/noordergraf/` wordt geback-upt in
git (subdirectory `/net/noordergraf/.git/`). Belangrijkste uitzondering
zijn de foto's van graven. Zie `.gitignore` in diverse directory's.

Na een grote bewerking van data doe je een *git add* / *git commit*
om de data zoals ze op dat moment zijn vast te leggen in git. Zo hou
je een overzicht van wat er is veranderd, en kunnen fouten later
hersteld worden.

De data in git wordt weer geback-upt naar github door *git push*.

Voor het browsen van veranderingen in de data kun je het best op
github kijken, al zijn er ook grafische tools om zoiets lokaal te
doen.


## Wiki op github

Er is een lokale kloon van de
[wiki op github](https://github.com/rug-compling/noordergraf/wiki) in
`/net/noordergraf/wiki/`.
Deze kun je bijwerken door in die directory *git pull* te doen.

We zouden deze bestanden eventueel beschikbaar kunnen maken op de
website van noordergraf, bijvoorbeeld door het geheel te verplaatsen
naar `/net/noordergraf/www/wiki/`.


## Website

De statische deel van de website staat in `/net/noordergraf/www/`. De
startpagina is `index.md` voor webbrowsers, of een van de andere
index-bestanden, afhankelijk van waar een applicatie om vraagt.

TODO: overige pagina's beschrijven

TODO: tweetaligheid

TODO: de handler voor de verwerking van turle-bestanden: uitvoer
afhankelijk van gevraagde content-type (html, rdf (default), etc.),
zie ook Multiviews in `/net/noordergraf/www/.htaccess` .


## Onze versie van soundex

In `/net/noordergraf/go/nlsoundex/`


TODO: beschrijving van algoritme

TODO: hoe kunnen mensen dit gebruiken in hun query's?


## Bestandsrechten

TODO: `/net/noordergraf/setall.sh`
