package shapefile

import (
	"os"
	"testing"
)

const dbf_test_fn = "test/Geometrie_Wahlkreise_18DBT.dbf"
const dbf_test2_fn = "test/ne_50m_admin_0_countries.dbf"

func TestDBFSimple2(t *testing.T) {
	// dbf parser barfed on something in this file,
	// start debug ...
	file, err := os.Open(dbf_test2_fn)
	if err != nil {
		t.Errorf("Failed opening file: %s", err.Error())
	}
	_, err = NewDBFFile(file)

	if err != nil {
		t.Error(err)
	}
}

func TestDBFHeadSimple(t *testing.T) {
	file, _ := os.Open(dbf_test_fn)
	defer file.Close()
	hdr, err := NewDBFFileHeader(file)

	if err != nil {
		t.Fail()
	}
	if hdr.Version != 3 {
		t.Fail()
	}
	if hdr.NumRecords != 299 {
		t.Fail()
	}
}

func TestDBFSimple(t *testing.T) {
	file, _ := os.Open(dbf_test_fn)
	defer file.Close()
	f, err := NewDBFFile(file)

	if err != nil {
		t.Errorf(err.Error())
	}
	if 4 != len(f.FieldDescriptors) {
		t.Errorf("fielddesc len != 4")
	}
	if f.FieldDescriptors[0].FieldName() != "WKR_NR" {
		t.Errorf("fieldname 0 not WKR_NR")
	}
	if f.FieldDescriptors[3].FieldName() != "LAND_NAME" {
		t.Errorf("fieldname 3 not LAND_NAME")
	}

	if 299 != len(f.Entries) {
		t.Errorf("incorrect number of entries: %d", len(f.Entries))
	}

	//	for _, fd := range f.FieldDescriptors {
	//		println(fd.String())
	//	}

	//	for _, entry := range f.Entries {
	//		for _, e := range entry {
	//			switch e.(type) {
	//			case string:
	//				println(e.(string))
	//			case int64:
	//				println(e.(int64))
	//			default:
	//				println("?")
	//			}
	//		}
	//	}

}
