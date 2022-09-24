#!/usr/bin/env python3

# Dit script maakt de indexen voor 'free text search' aan in AllegroGraph
#
# Dit script wordt aangeroepen vanuit /net/noordergraf/data/input

import os

AGRAPH_USER     = os.environ.get('AGRAPH_USER')
AGRAPH_PASSWORD = os.environ.get('AGRAPH_PASSWORD')

from franz.openrdf.connect import ag_connect

conn = ag_connect('noordergraf', create=False, user=AGRAPH_USER, password=AGRAPH_PASSWORD)

conn.createFreeTextIndex("all", predicates=None, indexLiterals=True, indexResources=False, indexFields=["object"], minimumWordSize=2, stopWords=[], wordFilters=["drop-accents"], innerChars=["alphanumeric"], borderChars=None, tokenizer="default")

conn.createFreeTextIndex("fullname", predicates=["<https://noordergraf.rug.nl/ns#fullNameSE>"], indexLiterals=True, indexResources=False, indexFields=["object"], minimumWordSize=1, stopWords=[], wordFilters=[], innerChars=["alphanumeric"], borderChars=None, tokenizer="default")
