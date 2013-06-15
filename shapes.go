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
	binary.Read(r, binary.LittleEndian, &t)
	if t != t_expected {
		return fmt.Errorf("unexpected type: %d expected: %d", t, t_expected)
	}
	return nil
}

func RecordRecordContent(r io.Reader) (content RecordContent, err error) {
	var typ ShapeType
	if err = binary.Read(r, L, &typ); err != nil {
		return
	}
	// this implementation does not enforce the rule that all
	// records in a file must be the same type.
	switch typ {
	case NULL_SHAPE:
		return ReadNull(r)
	case POINT:
		return ReadPoint(r)
	case POLY_LINE:
		return ReadPolyLine(r)
	case POLYGON:
		return ReadPolygon(r)
	case MULTI_POINT:
		return ReadMultiPoint(r)
	case POINT_Z:
		return ReadPointZ(r)
	case POLY_LINE_Z:
		return ReadPolyLineZ(r)
	case POLYGON_Z:
		return ReadPolygonZ(r)
	case MULTI_POINT_Z:
		return ReadMultiPointZ(r)
	case POINT_M:
		return ReadPointM(r)
	case POLY_LINE_M:
		return ReadPolyLineM(r)
	case POLYGON_M:
		return ReadPolygonM(r)
	case MULTI_POINT_M:
		return ReadMultiPointM(r)
	case MULTI_PATCH:
		return ReadMultiPatch(r)
	default:
		err = fmt.Errorf("unknown shape type: %d", typ)
		return
	}
}

type Null struct {
}

func ReadNull(r io.Reader) (n *Null, err error) {
	return &Null{}, nil
}

type Point struct {
	//_ ShapeType // Point is an embedded type
	X float64
	Y float64
}

func (p *Point) String() string {
	return fmt.Sprintf("X: %0.2f Y:%0.2f", p.X, p.Y)
}

func ReadPoint(r io.Reader) (p *Point, err error) {
	err = readType(r, POINT)
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

func (b *Box) String() string {
	return fmt.Sprintf("Xmin: %0.2f Xmax: %0.2f Ymin: %0.2f Ymax: %0.2f", b.Xmin, b.Xmax, b.Ymin, b.Ymax)

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
	Points    []Point
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
	pl.Points = make([]Point, pl.NumPoints)

	err = binary.Read(r, L, &pl.Parts)

	return

}

func (p *PolyLine) String() string {
	str := fmt.Sprintf("Box: \n%s", p.Box.String())
	str += fmt.Sprintf("NumParts: %d\n", p.NumParts)
	str += fmt.Sprintf("NumPoints: %d\n", p.NumPoints)
	for i, p := range p.Parts {
		str += fmt.Sprintf("part %d : %d\n", i, p)
	}
	for i, p := range p.Points {
		str += fmt.Sprintf("point %d : %s", i, p.String())
	}
	return str
}

type Polygon struct {
	PolyLine
}

func (p *Polygon) String() string {
	str := fmt.Sprintf("Box: %s\n", p.Box.String())
	str += fmt.Sprintf("NumParts: %d\n", p.NumParts)
	str += fmt.Sprintf("NumPoints: %d\n", p.NumPoints)
	for i, p := range p.Parts {
		str += fmt.Sprintf("part %d : %d\n", i, p)
	}
	for i, p := range p.Points {
		str += fmt.Sprintf("point %d : %s\n", i, p.String())
	}
	return str

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
	if err = binary.Read(r, L, pg.Parts); err != nil {
		return
	}
	pg.Points = make([]Point, pg.NumPoints)

	err = binary.Read(r, L, pg.Points)

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
	MRange    MRange    // optional
	MArray    []float64 // optional
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
	MRange    MRange    // optional
	MArray    []float64 // optional
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

func ReadPointZ(r io.Reader) (p *PointZ, err error) {
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
	Points    []Point
	ZRange    ZRange
	ZArray    []float64
	MRange    MRange    // optional
	MArray    []float64 // optional
}

func ReadMultiPointZ(r io.Reader) (mp *MultiPointZ, err error) {
	mp = &MultiPointZ{}
	if err = binary.Read(r, L, mp.Box); err != nil {
		return
	}
	if err = binary.Read(r, L, &mp.NumPoints); err != nil {
		return
	}
	mp.Points = make([]Point, mp.NumPoints)
	if err = binary.Read(r, L, &mp.Points); err != nil {
		return
	}
	if err = binary.Read(r, L, mp.ZRange); err != nil {
		return
	}
	mp.ZArray = make([]float64, mp.NumPoints)
	if err = binary.Read(r, L, mp.ZArray); err != nil {
		return
	}
	if err = binary.Read(r, L, mp.MRange); err != nil {
		return
	}
	mp.MArray = make([]float64, mp.NumPoints)

	err = binary.Read(r, L, mp.MArray)
	return
}

type PolyLineZ struct {
	Box       Box
	NumParts  int32
	NumPoints int32
	Parts     []int32
	Points    []Point
	ZRange    ZRange
	ZArray    []float64 //optional
	MRange    MRange    // optional
	MArray    []float64 //optional
}

func ReadPolyLineZ(r io.Reader) (pl *PolyLineZ, err error) {
	pl = &PolyLineZ{}

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
	pl.Points = make([]Point, pl.NumPoints)

	if err = binary.Read(r, L, pl.ZRange); err != nil {
		return
	}

	pl.ZArray = make([]float64, pl.NumPoints)
	if err = binary.Read(r, L, pl.ZArray); err != nil {
		return
	}
	if err = binary.Read(r, L, pl.MRange); err != nil {
		return
	}
	pl.MArray = make([]float64, pl.NumPoints)
	if err = binary.Read(r, L, pl.MArray); err != nil {
		return
	}
	return

}

type PolygonZ struct {
	PolyLineZ
}

func ReadPolygonZ(r io.Reader) (pg *PolygonZ, err error) {
	pg = &PolygonZ{}

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
	pg.Points = make([]Point, pg.NumPoints)

	if err = binary.Read(r, L, pg.ZRange); err != nil {
		return
	}

	pg.ZArray = make([]float64, pg.NumPoints)
	if err = binary.Read(r, L, pg.ZArray); err != nil {
		return
	}
	if err = binary.Read(r, L, pg.MRange); err != nil {
		return
	}
	pg.MArray = make([]float64, pg.NumPoints)
	if err = binary.Read(r, L, pg.MArray); err != nil {
		return
	}
	return

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
	Points    []Point
	ZRange    ZRange    // optional
	MRange    MRange    // optional
	MArray    []float64 // optional
}

func ReadMultiPatch(r io.Reader) (mp *MultiPatch, err error) {
	mp = &MultiPatch{}
	if err = binary.Read(r, L, mp.Box); err != nil {
		return
	}
	if err = binary.Read(r, L, &mp.NumParts); err != nil {
		return
	}
	if err = binary.Read(r, L, &mp.NumPoints); err != nil {
		return
	}
	mp.Parts = make([]int32, mp.NumParts)
	if err = binary.Read(r, L, &mp.Parts); err != nil {
		return
	}
	mp.PartTypes = make([]PartType, mp.NumParts)
	if err = binary.Read(r, L, &mp.PartTypes); err != nil {
		return
	}
	mp.Points = make([]Point, mp.NumPoints)
	if err = binary.Read(r, L, mp.Points); err != nil {
		return
	}
	if err = binary.Read(r, L, mp.ZRange); err != nil {
		return
	}
	if err = binary.Read(r, L, mp.MRange); err != nil {
		return
	}
	mp.MArray = make([]float64, mp.NumPoints)
	err = binary.Read(r, L, mp.NumPoints)
	return

}
