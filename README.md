# Shapefile - elementary Shapefile processing.

!!! This library is far from complete - Work in Progress !!!

Currently, the code can read shapefile: `.shp` files, though only tested
for a single set of files containing Polygons.

It can read `.dbf` files, though only a very limited subset ('C' and 'N'
datadiles)

Not supported are the `.shx` or any of the additional meta data files
not specified in the [ESRI
Whitepaper](http://www.esri.com/library/whitepapers/pdfs/shapefile.pdf)
Because I've not been able to find a proper format spec.

Basically, it serves a very small purpose that I needed and it's
published here for a.) backup b.) to save someone work who wants to get
started working on a proper library.

I may continue work on this at some point in the future. 

## TODO

- interface and doc
- write support
- random access to records via `.shx`
- figure out ancilliary file formats (.prj, .sbn, .shp.xml, ...)
- find more complete / diverse sample data for testing
- export / convert to other formats (geojson?)



