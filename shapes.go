package shapefile

import (
	"encoding/binary"
	"fmt"
	"io"
)

var L = binary.LittleEndian
var B = binary.BigEndian

func readType(r io.Reader, t_expected ShapeType) error {
	var t ShapeType
	binary.Read(r, binary.LittleEndian, t)
	if t != t_expected {
		return fmt.Errorf("unexpected type: %d expected: %d", t, t_expected)
	}
	return nil
}

type Null struct {
}

func ReadNull(r io.Reader) (n *Null, err error) {
	err = readType(r, NULL_SHAPE)
	return &Null{}, nil
}

type Point struct {
	_ ShapeType // Point is an embedded type
	X float64
	Y float64
}

func ReadPoint(r io.Reader) (p *Point, err error) {
	//err= readType(r, POINT)
	p = &Point{}
	err = binary.Read(r, binary.LittleEndian, p)
	return
}

type Box struct {
	Xmin float64
	Ymin float64
	Xmax float64
	Ymax float64
}
type MultiPoint struct {
	Box       Box
	NumPoints int32
	Points    []*Point
}

func ReadMultiPoint(r io.Reader) (mp *MultiPoint, err error) {
	if err = readType(r, MULTI_POINT); err != nil {
		return
	}
	mp = &MultiPoint{}
	if err = binary.Read(r, binary.LittleEndian, &mp.Box); err != nil {
		return
	}
	if err = binary.Read(r, binary.LittleEndian, &mp.NumPoints); err != nil {
		return
	}

	mp.Points = make([]*Point, mp.NumPoints)
	if err = binary.Read(r, binary.LittleEndian, mp.Points); err != nil {
		return
	}

	return
}

type PolyLine struct {
	Box       Box
	NumParts  int32
	NumPoints int32
	Parts     []int32
	Points    []*Point
}

func ReadPolyLine(r io.Reader) (pl *PolyLine, err error) {
	if err = readType(r, POLY_LINE); err != nil {
		return
	}
	pl = &PolyLine{}
	if err = binary.Read(r, L, &pl.Box); err != nil {
		return
	}
	if err = binary.Read(r, L, &pl.NumParts); err != nil {
		return
	}
	if err = binary.Read(r, L, &pl.NumPoints); err != nil {
		return
	}
	pl.Parts = make([]int32, pl.NumParts)
	if err = binary.Read(r, L, &pl.Parts); err != nil {
		return
	}
	pl.Points = make([]*Point, pl.NumPoints)

	err = binary.Read(r, L, &pl.Parts)

	return

}

type Polygon struct {
	PolyLine
}

func ReadPolygon(r io.Reader) (pg *Polygon, err error) {
	if err = readType(r, POLYGON); err != nil {
		return
	}
	pg = &Polygon{}
	if err = binary.Read(r, L, &pg.Box); err != nil {
		return
	}
	if err = binary.Read(r, L, &pg.NumParts); err != nil {
		return
	}
	if err = binary.Read(r, L, &pg.NumPoints); err != nil {
		return
	}
	pg.Parts = make([]int32, pg.NumParts)
	if err = binary.Read(r, L, &pg.Parts); err != nil {
		return
	}
	pg.Points = make([]*Point, pg.NumPoints)

	err = binary.Read(r, L, &pg.Parts)

	return

}

type PointM struct {
	X float64
	Y float64
	M float64
}

func ReadPointM(r io.Reader) (pm *PointM, err error) {
	if err = readType(r, POINT_M); err != nil {
		return
	}
	pm = &PointM{}
	err = binary.Read(r, L, pm)
	return
}

type MRange struct {
	Mmin float64
	Mmax float64
}
type MultiPointM struct {
	Box       Box
	NumPoints int32
	Points    []*Point
	MRange    MRange
	MArray    []float64
}

func ReadMultiPointM(r io.Reader) (mp *MultiPointM, err error) {
	if err = readType(r, MULTI_POINT_M); err != nil {
		return
	}
	mp = &MultiPointM{}
	if err = binary.Read(r, L, &mp.Box); err != nil {
		return
	}
	if err = binary.Read(r, L, &mp.NumPoints); err != nil {
		return
	}
	mp.Points = make([]*Point, mp.NumPoints)
	if err = binary.Read(r, L, mp.Points); err != nil {
		return
	}
	if err = binary.Read(r, L, &mp.MRange); err != nil {
		return
	}
	mp.MArray = make([]float64, mp.NumPoints)
	err = binary.Read(r, L, mp.MArray)
	return

}

type PolyLineM struct {
	Box       Box
	NumParts  int32
	NumPoints int32
	Parts     []int32
	Points    []*Point
	MRange    MRange
	MArray    []float64
}

func ReadPolyLineM(r io.Reader) (pl *PolyLineM, err error) {
	if err = readType(r, POLY_LINE_M); err != nil {
		return
	}
	pl = &PolyLineM{}
	if err = binary.Read(r, L, pl.Box); err != nil {
		return
	}
	if err = binary.Read(r, L, &pl.NumParts); err != nil {
		return
	}
	if err = binary.Read(r, L, &pl.NumPoints); err != nil {
		return
	}
	pl.Parts = make([]int32, pl.NumParts)
	if err = binary.Read(r, L, pl.Parts); err != nil {
		return
	}
	pl.Points = make([]*Point, pl.NumPoints)
		if err = binary.Read(r, L, pl.Points); err != nil {
		return
	}

	if err = binary.Read(r, L, pl.MRange); err != nil {
		return
	}
	pl.MArray = make([]float64, pl.NumPoints)

	err = binary.Read(r, L, pl.MArray)
	return
}

type PolygonM struct {
	PolyLineM
}

func ReadPolygonM(r io.Reader) (pg *PolygonM, err error) {
	if err = readType(r, POLYGON_M); err != nil {
		return
	}
	pg = &PolygonM{}
	if err = binary.Read(r, L, pg.Box); err != nil {
		return
	}
	if err = binary.Read(r, L, &pg.NumParts); err != nil {
		return
	}
	if err = binary.Read(r, L, &pg.NumPoints); err != nil {
		return
	}
	pg.Parts = make([]int32, pg.NumParts)
	if err = binary.Read(r, L, pg.Parts); err != nil {
		return
	}
	pg.Points = make([]*Point, pg.NumPoints)
	if err = binary.Read(r, L, pg.Points); err != nil {
		return
	}
	if err = binary.Read(r, L, pg.MRange); err != nil {
		return
	}
	pg.MArray = make([]float64, pg.NumPoints)

	err = binary.Read(r, L, pg.MArray)
	return
}


type PointZ struct {
	X float64
	Y float64
	Z float64
	M float64
}

func ReadPointZ(r io.Reader)(p *PointZ, err error) {
	if err = readType(r, POINT_Z); err != nil {
		return
	}
	p = &PointZ{}
	err = binary.Read(r, L, &p)
	return

}

type ZRange struct {
	Zmin float64
	Zmax float64
}
type MultiPointZ struct {
	Box       Box
	NumPoints int32
	Points    []*Point
	ZRange    ZRange
	ZArray    []float64
	MRange    MRange
	MArray    []float64
}
type PolyLineZ struct {
	Box       Box
	NumParts  int32
	NumPoints int32
	Parts     []int32
	Points    []*Point
	ZRange    ZRange
	ZArray    []float64
	MRange    MRange
	MArray    []float64
}
type PolygonZ struct {
	PolyLineZ
}

type PartType int32

const (
	TRIANGLE_STRIP PartType = iota
	TRIANGLE_FAN
	OUTER_RING
	INNER_RING
	FIRST_RING
	RING
)

func (p PartType) String() string {
	switch p {
	case TRIANGLE_STRIP:
		return "TRIANGLE_STRIP"
	case TRIANGLE_FAN:
		return "TRIANGLE_FAN"
	case OUTER_RING:
		return "OUTER_RING"
	case INNER_RING:
		return "INNER_RING"
	case FIRST_RING:
		return "FIRST_RING"
	case RING:
		return "RING"
	default:
		return "UNKNOWN"
	}
}

type MultiPatch struct {
	Box       Box
	NumParts  int32
	NumPoints int32
	Parts     []int32
	PartTypes []PartType
	Points    []*Point
	ZRange    ZRange
	MRange    MRange
	MArray    []float64
}
