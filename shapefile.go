
package shapefile

import (
	"io"
)

type Shapefile struct {
	header *MainFileHeader
	records []*Record
}

type Record struct {
	header *MainFileRecordHeader
	content RecordContent
}

type RecordContent interface{}

func NewShapefile(rdr io.Reader) (s *Shapefile, err error) {
	s = &Shapefile{}

	var h *MainFileHeader
	if h, err = NewMainFileHeaderFromReader(rdr); err != nil {
		return
	}

	s.header = h
	i := s.header.FileLength - 50 // length of header = 100 bytes = 50 words
	var rh *MainFileRecordHeader
	var rec *Record
	for {
		if i <=0 {
			break
		}
		if rh, err = NewMainFileRecordHeaderFromReader(rdr); err != nil {
			return
		}
		i = i-rh.ContentLength-4
		rec = &Record{}
		rec.header = rh
		if rec.content, err = RecordRecordContent(rdr); err != nil {
			return
		}
		s.records = append(s.records, rec)
	}
	return
}
