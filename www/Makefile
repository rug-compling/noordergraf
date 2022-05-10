
%.en : %.IN
	perl -p -e 's/\{\{.*?\|(.*?)\}\}/\1/g' $< > $@

%.nl : %.IN
	perl -p -e 's/\{\{(.*?)\|.*?\}\}/\1/g' $< > $@


all: \
	mappob.html.en \
	mappob.html.nl \
	mappod.html.en \
	mappod.html.nl \
	mapsites.html.en \
	mapsites.html.nl
