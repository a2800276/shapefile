package shapefile

import (
	"io"
)

type Shapefile struct {
	Header  *MainFileHeader
	Records []*Record
}

type Record struct {
	Header  *MainFileRecordHeader
	Content RecordContent
}

type RecordContent interface{}

func NewShapefile(rdr io.Reader) (s *Shapefile, err error) {
	s = &Shapefile{}

	var h *MainFileHeader
	if h, err = NewMainFileHeaderFromReader(rdr); err != nil {
		return
	}

	s.Header = h
	i := s.Header.FileLength - 50 // length of header = 100 bytes = 50 words
	var rh *MainFileRecordHeader
	var rec *Record
	for {
		if i <= 0 {
			break
		}
		if rh, err = NewMainFileRecordHeaderFromReader(rdr); err != nil {
			return
		}
		i = i - rh.ContentLength - 4
		rec = &Record{}
		rec.Header = rh
		if rec.Content, err = RecordRecordContent(rdr); err != nil {
			return
		}
		s.Records = append(s.Records, rec)
	}
	return
}
