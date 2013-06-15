package shapefile

import (
	"encoding/binary"
	"os"
	"testing"
)

const testfile = "test/Geometrie_Wahlkreise_18DBT.shp"
const testfileTrunc = "test/truncated.shp"
const testfileInv = "test/invalidhdr.shp"

func TestMainFileHeaderRead(t *testing.T) {
	file, _ := os.Open(testfile)
	defer file.Close()
	hdr, _ := NewMainFileHeaderFromReader(file)
	expected := `FileLength 152120
Version 1000
ShapeType POLYGON
Xmin 280399.43
Ymin 5235857.07
Xmax 921298.68
Ymax 6101309.73
Zmin 0.00
Zmax 0.00
Mmin 0.00
Mmax 0.00
`
	if expected != hdr.String() {
		t.Fail()
	}
}

func TestMainFileHeaderReadNotEnough(t *testing.T) {
	file, _ := os.Open(testfileTrunc)
	defer file.Close()
	_, err := NewMainFileHeaderFromReader(file)
	if "can't read UNUSED" != err.Error() {
		t.Fail()
	}

}
func TestMainFileHeaderReadInvalid(t *testing.T) {
	file, _ := os.Open(testfileInv)
	defer file.Close()
	_, err := NewMainFileHeaderFromReader(file)
	if "invalid fileCode: 654966784" != err.Error() {
		t.Fail()
	}
}

type ttt struct {
	X int32
	Y int32
}
type tttt struct {
	X int32
	B [40]byte
	A ttt
}

func TestBinaryTesting(t *testing.T) {
	file, _ := os.Open(testfile)
	defer file.Close()
	test := tttt{}
	binary.Read(file, binary.LittleEndian, &test)
	println(test.X)
	println(test.A.X)
	println(test.A.Y)
}

func TestTest(t *testing.T) {
	file, _ := os.Open(testfile)
	defer file.Close()
	Parse(file)
}
