---
title: Noordergraf -- {{Overzicht|Overview}}
class: pad
---

# {{Overzicht|Overview}}

{{Zie ook|See also}}: [details](ns)

## Class

{{Alles (elke instance van een class) kan een 'comment' hebben:|}}
{{|Everything (each class instance) can have a 'comment':}}

- **Class** --- **Item** | **Place** | **Site** | ...[etc](ns#Class)
    - comment → xsd:string | rdf:langString

{{Bij alles wat volgt is 'comment' steeds weggelaten.|}}
{{|In everything that follows, 'comment' is omitted.}}

## Item

{{Voorbeeld|Example}}: [item:T00000](item/T00000)

- **Item** --- **Cross** | **Tomb**
    - commemorator → **[Person](#Person)**
    - creator
        - **Creator**
            - location → **[Location](#Location)**
            - name → **[Name](#Name)**
    - date → xsd:date | xsd:gYearMonth | xsd:gYear
    - dcModified → xsd:dateTime
    - figure
        - **Figure** --- **Bust** | **Portrait** | **Statue**
            - theme → **[Person](#Person)**
    - geo
        - **<a name="GeoPoint">GeoPoint</a>**
            - geoLat → xsd:decimal
            - geoLong → xsd:decimal
    - lot → xsd:string
    - nd → ll:
    - photo
        - **Photo**
            - dcCreator → rdfs:Resource | xsd:string
            - dcDate → xsd:dateTime
            - dcLicense → rdfs:Resource | xsd:string
            - file → `img:...`
    - quote
        - **Quote**
            - quoteSubject → **[Person](#Person)**
            - text → xsd:string | rdf:langString
            - reference
                - **Reference**
                    - --- **BibleReference**
                        - bibleSource → `https://www.bible.com/bible/...`
                        - book → **[BibleBook](#BibleBook)**
                        - chapter → xsd:integer
                        - verse → xsd:string | rdf:langString
                    - --- **BookReference**
                        - author
                            - **Author**
                                - authorName → xsd:string | rdf:langString
                                - sameAs → `wikidata:...`
                        - bookTitle → xsd:string | rdf:langString
                        - date → xsd:date | xsd:gYearMonth | xsd:gYear
                        - sameAs → `wikidata:...`
                    - --- **FolkReference**
                        - seeAlso → rdfs:Resource
                        - sameAs → `wikidata:...`
                    - --- **HymnReference**
                        - hymnSource → rdfs:Resource
                        - hymnTitle → xsd:string | rdf:langString
                    - --- **MusicReference**
                        - album
                            - **Album**
                                - albumTitle → xsd:string | rdf:langString
                                - date → xsd:date | xsd:gYearMonth | xsd:gYear
                                - discogs → `discogs:...`
                        - albumTrack → xsd:string
                        - artist
                            - **Artist**
                                - artistName → xsd:string | rdf:langString
                                - sameAs → `wikidata:...`
                        - trackTitle → xsd:string | rdf:langString
    - site → **[Site](#Site)**
    - symbol
        - **Symbol**
            - symbolType  → **[SymbolType](#SymbolType)**
    - text → xsd:string | rdf:langString
    - todo → **[ToDo](#ToDo)**
    - subject
        - **<a name="Person">Person</a>** --- **Female** | **Male** | ...[meer](ns#Person)
            - age
                - **Age**
                    - duration → xsd:duration
                    - inferred → xsd:boolean
                    - text → xsd:string | rdf:langString
            - dateOfBirth
                - **<a name="Date">Date</a>**
                    - date → xsd:date | xsd:gYearMonth | xsd:gYear
                    - inferred → xsd:boolean
            - dateOfDeath → **[Date](#Date)**
            - location
                - **<a name="Location">Location</a>** --- **Building** | **City** | **Colony** | ...[meer](ns#Location)
                    - coordinates
                        - **Coordinates**
                            - geo → **[GeoPoint](#GeoPoint)**
                            - nd → ll:
                    - partOf → **[Location](#Location)**
                    - place → **[Place](#Place)**
                    - text → xsd:string | rdf:langString
            - name
                - **<a name="Name">Name</a>**
                    - fullName → xsd:string
                    - fullNameSE → xsd:string
                    - prefix → xsd:string
                    - suffix → xsd:string
                    - text → xsd:string | rdf:langString
            - occupation
                - **<a name="Occupation">Occupation</a>**
                    - after → **[Occupation](#Occupation)** | **[Decoration](#Decoration)** | **[Relation](#Relation)**
                    - before → **[Occupation](#Occupation)** | **[Decoration](#Decoration)** | **[Relation](#Relation)**
                    - begin → xsd:date | xsd:gYearMonth | xsd:gYear
                    - end → xsd:date | xsd:gYearMonth | xsd:gYear
                    - hisco → `hisco:...`
                    - location → **[Location](#Location)**
                    - organisation
                        - **Organisation** --- **Battalion** | **Community** | **Company** | ...[meer](ns#Organisation)
                            - location → **[Location](#Location)**
                            - organisationName → xsd:string | rdf:langString
                    - period
                        - **<a name="Period">Period</a>**
                            - duration → xsd:duration
                            - text → xsd:string | rdf:langString
                    - text → xsd:string | rdf:langString
            - placeOfBirth → **[Location](#Location)**
            - placeOfDeath → **[Location](#Location)**
            - property
                - **Property**
                    - --- **Decoration**
                        - begin → xsd:date | xsd:gYearMonth | xsd:gYear
                        - end → xsd:date | xsd:gYearMonth | xsd:gYear WEG ??
                        - decorationTitle → xsd:string | rdf:langString
                        - text → xsd:string | rdf:langString
            - quantity
                - **Quantity**
                    - number → xsd:integer
            - relation
                - **<a name="Relation">Relation</a>** --- **Aunt** | **Brother** | **BrotherInLaw** | ...[meer](ns#Relation)
                    - after → **[Relation](#Relation)** | **[Decoration](#Decoration)** | **[Occupation](#Occupation)**
                    - before → **[Relation](#Relation)** | **[Decoration](#Decoration)** | **[Occupation](#Occupation)**
                    - begin → xsd:date | xsd:gYearMonth | xsd:gYear
                    - end → xsd:date | xsd:gYearMonth | xsd:gYear
                    - inferred → xsd:boolean
                    - period → **[Period](#Period)**
                    - target → **[Person](#Person)**
                    - text → xsd:string | rdf:langString
            - sameAs → `wikidata:...`
            - text → xsd:string | rdf:langString

## BibleBook

{{Voorbeeld|Example}}: [bible:Matthew](bible/Matthew)

- **<a name="BibleBook">BibleBook</a>**
    - bibleBookName → rdf:langString
    - dcModified → xsd:dateTime
    - order → xsd:integer
    - sameAs → `wikidata:...`

## Place

{{Voorbeeld|Example}}: [place:Aalden2760142](place/Aalden2760142)

- **<a name="Place">Place</a>**
    - dcModified → xsd:dateTime
    - geo → **[GeoPoint](#GeoPoint)**
    - nd → ll:
    - placeName → xsd:string
    - sameAs → `wikidata:...` | `https://sws.geonames.org/...`

## Site

{{Voorbeeld|Example}}: [site:NLdr7811heEmmen](site/NLdr7811heEmmen)

- **<a name="Site">Site</a>**
    - dcModified → xsd:dateTime
    - geo → **[GeoPoint](#GeoPoint)**
    - nd → ll:
    - sameAs → `wikidata:...` | `https://sws.geonames.org/...`
    - seeAlso → rdfs:Resource
    - siteAddress
        - **PostalAddress**
            - streetAddress → xsd:string
            - addressLocality → xsd:string
            - addressRegion → xsd:string
            - addressCountry → xsd:string
            - postalCode → xsd:string
    - siteName → xsd:string

## SymbolType

{{Voorbeeld|Example}}: [symbol:Anchor](symbol/Anchor)

- **<a name="SymbolType">SymbolType</a>**
    - dcModified → xsd:dateTime
    - pictogram
        - **Pictogram**
            - dcCreator → xsd:string | rdfs:Resource
            - dcLicense → rdfs:Resource | xsd:string
            - dcSource → rdfs:Resource
            - dcTitle → xsd:string | rdf:langString
            - file → `picto:...`
    - sameAs → `wikidata:...`
    - seeAlso → `dodenakkers:...`
    - symbolDescription → rdf:langString

## ToDo

{{Voorbeeld|Example}}: [todo:text](todo/Text)

- **<a name="ToDo">ToDo</a>**
    - task → rdf:langString
