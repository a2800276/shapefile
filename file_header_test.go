package shapefile

import (
	"os"
	"testing"
)

const testfile = "test/Geometrie_Wahlkreise_18DBT.shp"

func TestMainFileHeaderRead(t *testing.T) {
	file, _ := os.Open(testfile)
	hdr, err := NewMainFileHeaderFromReader(file)
	println(hdr.String())
	println(err)
}
