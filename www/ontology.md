
## Plot

Example: [tomb:T00000](tomb/T00000)

- **Plot** --- **Cross** | **Tomb**
    - dcModified → `xsd:dateTime`
    - commemorator → **[Person](#Person)**
    - comment → `xsd:string` | `rdf:langString`
    - creator
        - **Creator**
            - location → **[Location](#Location)**
            - name → **[Name](#Name)**
    - date → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
    - figure
        - **Figure** --- **Bust** | **Portrait** | **Statue**
            - theme → **[Person](#Person)**
    - photo
        - **Photo**
            - dcCreator → `rdfs:Resource` | `xsd:string`
            - dcDate → `xsd:dateTime`
            - dcLicense → `rdfs:Resource` | `xsd:string`
            - file → img:...
            - geo
                - **<a name="GeoPoint">GeoPoint</a>**
                    - geoLat → `xsd:decimal`
                    - geoLong → `xsd:decimal`
            - nd → `ll:`
    - quote
        - **Quote**
            - comment → `xsd:string` | `rdf:langString`
            - quoteSubject → **[Person](#Person)**
            - reference
                - **Reference**
                    - --- **BibleReference**
                        - bibleSource → https:\/\/www\.bible\.com\/bible\/...
                        - chapter → `xsd:integer`
                        - comment → `xsd:string` | `rdf:langString`
                        - part → **[BiblePart](#BiblePart)**
                        - verse → `xsd:string` | `rdf:langString`
                    - --- **BookReference**
                        - sameAs → wikidata:...
                        - author
                            - **Author**
                                - sameAs → wikidata:...
                                - authorName → `xsd:string` | `rdf:langString`
                        - bookTitle → `xsd:string` | `rdf:langString`
                        - comment → `xsd:string` | `rdf:langString`
                        - date → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                    - --- **FolkReference**
                        - seeAlso → `rdfs:Resource`
                        - sameAs → wikidata:...
                    - --- **HymnReference**
                        - comment → `xsd:string` | `rdf:langString`
                        - hymnSource → `rdfs:Resource`
                        - hymnTitle → `xsd:string` | `rdf:langString`
                    - --- **MusicReference**
                        - album
                            - **Album**
                                - albumTitle → `xsd:string` | `rdf:langString`
                                - comment → `xsd:string` | `rdf:langString`
                                - date → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                                - discogs → discogs:...
                        - albumTrack → `xsd:string`
                        - artist
                            - **Artist**
                                - sameAs → wikidata:...
                                - artistName → `xsd:string` | `rdf:langString`
                        - comment → `xsd:string` | `rdf:langString`
                        - trackTitle → `xsd:string` | `rdf:langString`
            - text → `xsd:string` | `rdf:langString`
    - site → **[Site](#Site)**
    - symbol
        - **Symbol**
            - comment → `xsd:string` | `rdf:langString`
            - symbolType  → **[SymbolType](#SymbolType)**
    - text → `xsd:string` | `rdf:langString`
    - subject
        - **<a name="Person">Person</a>** --- **Female** | **Male** | ...[meer](ns#Person)
            - sameAs → wikidata:...
            - age
                - **Age**
                    - duration → `xsd:duration`
                    - inferred → `xsd:boolean`
                    - text → `xsd:string` | `rdf:langString`
            - dateOfBirth
                - **<a name="Date">Date</a>**
                    - date → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                    - inferred → `xsd:boolean`
            - dateOfDeath → **[Date](#Date)**
            - location
                - **<a name="Location">Location</a>** --- **Building** | **City** | **Colony** | ...[meer](ns#Location)
                    - comment → `xsd:string` | `rdf:langString`
                    - coordinates
                        - **Coordinates**
                            - geo → **[GeoPoint](#GeoPoint)**
                            - nd → `ll:`
                    - partOf → **[Location](#Location)**
                    - place → **[Place](#Place)**
                    - text → `xsd:string` | `rdf:langString`
            - name
                - **<a name="Name">Name</a>**
                    - fullName → `xsd:string`
                    - fullNameSE → `xsd:string`
                    - prefix → `xsd:string`
                    - suffix → `xsd:string`
                    - text → `xsd:string` | `rdf:langString`
            - occupation
                - **<a name="Occupation">Occupation</a>**
                    - after → **[Occupation](#Occupation)** | **[Relation](#Relation)**
                    - before → **[Occupation](#Occupation)** | **[Relation](#Relation)**
                    - begin → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                    - comment → `xsd:string` | `rdf:langString`
                    - end → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                    - hisco → hisco:...
                    - location → **[Location](#Location)**
                    - organisation
                        - **Organisation** --- **Battalion** | **Community** | **Company** | ...[meer](ns#Organisation)
                            - location → **[Location](#Location)**
                            - organisationName → `xsd:string` | `rdf:langString`
                    - period
                        - **<a name="Period">Period</a>**
                            - comment → `xsd:string` | `rdf:langString`
                            - duration → `xsd:duration`
                            - text → `xsd:string` | `rdf:langString`
                    - text → `xsd:string` | `rdf:langString`
            - placeOfBirth → **[Location](#Location)**
            - placeOfDeath → **[Location](#Location)**
            - property
                - **Property**
                    - --- **Decoration**
                        - begin → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                        - decorationTitle → `xsd:string` | `rdf:langString`
                        - end → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear` WEG ??
                        - text → `xsd:string` | `rdf:langString`
            - quantity
                - **Quantity**
                    - number → `xsd:integer`
            - relation
                - **<a name="Relation">Relation</a>** --- **Aunt** | **Brother** | **BrotherInLaw** | ...[meer](ns#Relation)
                    - after → **[Relation](#Relation)** | **[Occupation](#Occupation)**
                    - before → **[Relation](#Relation)** | **[Occupation](#Occupation)**
                    - begin → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                    - comment → `xsd:string` | `rdf:langString`
                    - end → `xsd:date` | `xsd:gYearMonth` | `xsd:gYear`
                    - inferred → `xsd:boolean`
                    - period → **[Period](#Period)**
                    - target → **[Person](#Person)**
                    - text → `xsd:string` | `rdf:langString`
            - text → `xsd:string` | `rdf:langString`

## BiblePart

Example: [bible:Genesis](bible/Genesis)

- **<a name="BiblePart">BiblePart</a>**
    - biblePartName → `rdf:langString`
    - sameAs → wikidata:...
    - dcModified → `xsd:dateTime`

## Place

Example: [place:Aalden2760142](place/Aalden2760142)

- **<a name="Place">Place</a>**
    - dcModified → `xsd:dateTime`
    - sameAs → wikidata:... | https:\/\/sws\.geonames\.org\/...
    - geo → **[GeoPoint](#GeoPoint)**
    - nd → `ll:`
    - placeName → `xsd:string`

## Site

Example: [site:NLdr7811heEmmen](site/NLdr7811heEmmen)

- **<a name="Site">Site</a>**
    - dcModified → `xsd:dateTime`
    - seeAlso → `rdfs:Resource`
    - sameAs → wikidata:... | https:\/\/sws\.geonames\.org\/...
    - geo → **[GeoPoint](#GeoPoint)**
    - nd → `ll:`
    - siteAddress
        - **PostalAddress**
            - addressCountry → `xsd:string`
            - addressLocality → `xsd:string`
            - addressRegion → `xsd:string`
            - postalCode → `xsd:string`
            - streetAddress → `xsd:string`
    - siteName → `xsd:string`

## SymbolType

Example: [symbol:Anchor](symbol/Anchor)

- **<a name="SymbolType">SymbolType</a>**
    - dcModified → `xsd:dateTime`
    - seeAlso → https:\/\/www\.dodenakkers\.nl\/naslag\/symboliek\/...
    - sameAs → wikidata:...
    - pictogram
        - **Pictogram**
            - dcCreator → `xsd:string` | `rdfs:Resource`
            - dcLicense → `rdfs:Resource` | `xsd:string`
            - dcSource → `rdfs:Resource`
            - dcTitle → `xsd:string` | `rdf:langString`
            - file → picto:...
            - pictogramParts WEG ??
                - **PictogramPart**
                    - dc:creator → `xsd:string` | `rdfs:Resource`
                    - dc:source → `rdfs:Resource`
                    - dc:title → `xsd:string` | `rdf:langString`


<pre>




















































</pre>
