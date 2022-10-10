## Data invoer

Dit kan alleen op machine `noordergraf` in deze directory


Invoer alle data in AllegroGraph (oude data wordt gewist):

    ./concat | ./input


Bij melding over fout in, bijvoorbeeld, regelnummer 12345:

    ./concat | ./show 12345

TODO: Een script waarmee je een enkel turtle-bestand in de database
kunt bijwerken.

TODO: Het scrip `concat` zal niet functioneren als er veel data bij
komt, omdat er dan te lange command lines in het script worden
gebruikt. Dat script zal dus vervangen moeten worden door iets anders.


## Directory's

Wanneer hier directory's met Turtle-bestanden bijkomen of weggaan, dan
ook de volgende bestanden aanpassen:
 
 - `/etc/apache2/conf-available/noordergraf.conf`
 - `/net/noordergraf/handlers/data/data.go`
 - `/net/noordergraf/www/index.nt`
 - `/net/noordergraf/www/index.penman`
 - `/net/noordergraf/www/index.rdf`
 - `/net/noordergraf/www/index.ttl`
