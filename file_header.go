package shapefile

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MainFileHeader struct {
	//	FileCode   int32
	//	Unused     [20]byte
	FileLength int32 // length in 16 bit words
	Version    int32
	ShapeType  ShapeType
	Xmin       float64
	Ymin       float64
	Xmax       float64
	Ymax       float64
	Zmin       float64
	Zmax       float64
	Mmin       float64
	Mmax       float64
}

func (h *MainFileHeader) String() string {
	str := fmt.Sprintf("FileLength %d\n", h.FileLength)
	str += fmt.Sprintf("Version %d\n", h.Version)
	str += fmt.Sprintf("ShapeType %s\n", h.ShapeType.String())
	str += fmt.Sprintf("Xmin %0.2f\n", h.Xmin)
	str += fmt.Sprintf("Ymin %0.2f\n", h.Ymin)
	str += fmt.Sprintf("Xmax %0.2f\n", h.Xmax)
	str += fmt.Sprintf("Ymax %0.2f\n", h.Ymax)
	str += fmt.Sprintf("Zmin %0.2f\n", h.Zmin)
	str += fmt.Sprintf("Zmax %0.2f\n", h.Zmax)
	str += fmt.Sprintf("Mmin %0.2f\n", h.Mmin)
	str += fmt.Sprintf("Mmax %0.2f\n", h.Mmax)
	return str
}

type ShapeType int32

const (
	NULL_SHAPE    ShapeType = 0
	POINT         ShapeType = 1
	POLY_LINE     ShapeType = 3
	POLYGON       ShapeType = 5
	MULTI_POINT   ShapeType = 8
	POINT_Z       ShapeType = 11
	POLY_LINE_Z   ShapeType = 13
	POLYGON_Z     ShapeType = 15
	MULTI_POINT_Z ShapeType = 18
	POINT_M       ShapeType = 21
	POLY_LINE_M   ShapeType = 23
	POLYGON_M     ShapeType = 25
	MULTI_POINT_M ShapeType = 28
	MULTI_PATCH   ShapeType = 31
)

func (s ShapeType) String() string {
	switch s {
	case NULL_SHAPE:
		return "NULL_SHAPE"
	case POINT:
		return "POINT"
	case POLY_LINE:
		return "POLY_LINE"
	case POLYGON:
		return "POLYGON"
	case MULTI_POINT:
		return "MULTI_POINT"
	case POINT_Z:
		return "POINT_Z"
	case POLY_LINE_Z:
		return "POLY_LINE_Z"
	case POLYGON_Z:
		return "POLYGON_Z"
	case MULTI_POINT_Z:
		return "MULTI_POINT_Z"
	case POINT_M:
		return "POINT_M"
	case POLY_LINE_M:
		return "POLY_LINE_M"
	case POLYGON_M:
		return "POLYGON_M"
	case MULTI_POINT_M:
		return "MULTI_POINT_M"
	case MULTI_PATCH:
		return "MULTI_PATCH"
	default:
		return "UNKNOWN"
	}
}

type MainFileRecordHeader struct {
	RecordNumber  int32
	ContentLength int32
}

func (h *MainFileRecordHeader) String() string {
	str := fmt.Sprintf("RecordNumber %d\n", h.RecordNumber)
	str += fmt.Sprintf("ContentLength %d\n", h.ContentLength)
	return str
}

func Parse (r io.Reader) {
	hdr, _ := NewMainFileHeaderFromReader(r)
	println(hdr.String())
	h2, _ := NewMainFileRecordHeaderFromReader(r)
	println(h2.String())
	pg, _ := ReadPolygon(r)
	println(pg.String())
	h2, _ = NewMainFileRecordHeaderFromReader(r)
	println(h2.String())
	pg, _ = ReadPolygon(r)
	println(pg.String())
	h2, _ = NewMainFileRecordHeaderFromReader(r)
	println(h2.String())
	pg, _ = ReadPolygon(r)
	println(pg.String())
}

func NewMainFileHeaderFromReader(r io.Reader) (hdr *MainFileHeader, err error) {

	var fileCode int32
	if err = binary.Read(r, binary.BigEndian, &fileCode); err != nil {
		return
	}
	if fileCode != 9994 {
		return nil, fmt.Errorf("invalid fileCode: %d", fileCode)
	}
	unused := make([]byte, 20)
	var n int
	if n, err = r.Read(unused); err != nil {
		return
	} else if n != 20 {
		return nil, fmt.Errorf("can't read UNUSED")
	}

	hdr = &MainFileHeader{}
	if err = binary.Read(r, binary.BigEndian, &hdr.FileLength); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Version); err != nil {
		return
	}
	if hdr.Version != 1000 {
		return nil, fmt.Errorf("Version must be 1000, is: %d", hdr.Version)
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.ShapeType); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Xmin); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Ymin); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Xmax); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Ymax); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Zmin); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Zmax); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Mmin); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &hdr.Mmax); err != nil {
		return
	}

	return
}

func NewMainFileRecordHeaderFromReader(r io.Reader) (hdr *MainFileRecordHeader, err error) {
	hdr = &MainFileRecordHeader{}
	if err = binary.Read(r, binary.BigEndian, &hdr.RecordNumber); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &hdr.ContentLength); err != nil {
		return
	}
	return
}
