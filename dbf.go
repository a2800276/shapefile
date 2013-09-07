package shapefile

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// DBF is documented here: http://www.clicketyclick.dk/databases/xbase/format/dbf.html

type DBFFile struct {
	DBFFileHeader    *DBFFileHeader
	FieldDescriptors []FieldDescriptor
	Entries          [][]interface{}
}

func NewDBFFile(r io.Reader) (dbf *DBFFile, err error) {
	dbf = &DBFFile{}
	if dbf.DBFFileHeader, err = NewDBFFileHeader(r); err != nil {
		return
	}
	len_fd := dbf.DBFFileHeader.LenHeader - 32 // the fixed portion of the header is 32 bytes
	num_fd := (int)(len_fd / 32)               // each field descriptor are 32 bytes each, see below.

	var fd FieldDescriptor
	for i := 0; i != num_fd; i++ {
		if err = binary.Read(r, L, &fd); err != nil {
			return
		}
		dbf.FieldDescriptors = append(dbf.FieldDescriptors, fd)
	}
	bullshitByte := make([]byte, 1)
	var n int
	if n, err = r.Read(bullshitByte); err != nil || n != 1 {
		if err == nil {
			err = fmt.Errorf("couldn't read bullshit byte!")
		}
		return
	}
	err = dbf.readEntries(r)
	return
}

func (dbf *DBFFile) readEntries(r io.Reader) (err error) {
	countRead := (uint32)(0)

	rawEntry := make([]byte, dbf.DBFFileHeader.LenRecord)
	var n int
	for {
		if n, err = r.Read(rawEntry); (err != nil) || n != (int)(dbf.DBFFileHeader.LenRecord) {
			if err == nil {
				err = fmt.Errorf("expected %d bytes, read: %d", dbf.DBFFileHeader.LenRecord, n)
			}
			return
		}
		if 0x2a == rawEntry[0] { // record deleted
			continue
		}

		entry := make([]interface{}, len(dbf.FieldDescriptors))
		var offset = 1

		for i, desc := range dbf.FieldDescriptors {
			rawField := rawEntry[offset : offset+(int)(desc.FieldLength)]
			offset += (int)(desc.FieldLength)

			switch desc.FieldType {
			case Character:
				entry[i] = (string)(rawField)
			case Number:
				numberStr := strings.TrimLeft((string)(rawField), " ")
				if entry[i], err = strconv.ParseInt(numberStr, 10, 64); err != nil {
					return
				}
			case Float:
				numberStr := strings.TrimLeft((string)(rawField), " ")
				if entry[i], err = strconv.ParseFloat(numberStr, 64); err != nil {
					return
				}

			default:
				err = fmt.Errorf("unsupported type: %c", desc.FieldType)
			}
		}
		dbf.Entries = append(dbf.Entries, entry)

		countRead++
		if countRead == dbf.DBFFileHeader.NumRecords {
			break
		}
	} // for
	return
}

// http://www.clicketyclick.dk/databases/xbase/format/dbf.html#DBF_STRUCT

type DBFFileHeader struct {
	Version        byte
	LastUpdate     [3]uint8 // YY MM DD (YY = years since 1900)
	NumRecords     uint32   // LittleEndian
	LenHeader      uint16
	LenRecord      uint16
	_              [2]byte // reserved
	IncompleteTx   byte
	EncFlag        byte
	FreeRecThread  uint32 // ...
	_              [8]byte
	MDXFlag        byte
	LanguageDriver byte
	_              [2]byte
}

func (hdr *DBFFileHeader) String() string {
	str := fmt.Sprintf("Version     : %d\n", hdr.Version)
	str += fmt.Sprintf("Last Update : %d %d %d\n", hdr.LastUpdate[0], hdr.LastUpdate[1], hdr.LastUpdate[2])
	str += fmt.Sprintf("Num Records : %d\n", hdr.NumRecords)
	str += fmt.Sprintf("Len Header  : %d\n", hdr.LenHeader)
	str += fmt.Sprintf("Len Record  : %d\n", hdr.LenRecord)
	return str
}

func NewDBFFileHeader(r io.Reader) (hdr *DBFFileHeader, err error) {
	hdr = &DBFFileHeader{}
	err = binary.Read(r, L, hdr)
	return
}

type FieldType byte

const (
	Character FieldType = 'C'
	Number              = 'N'
	Logical             = 'L'
	Date                = 'D'
	Memo                = 'M'
	Float               = 'F'
	// VarChar = ???
	Binary        = 'B'
	General       = 'G'
	Picture       = 'P'
	Currency      = 'Y'
	DateTime      = 'T'
	Integer       = 'I'
	VariField     = 'V'
	VarCharVar    = 'X'
	Timestamp     = '@'
	Double        = 'O' // 8 bytes
	Autoincrement = '+'
)

type FieldDescriptor struct {
	FieldName_     [11]byte
	FieldType      FieldType
	FieldDataAddr  uint32
	FieldLength    uint8
	DecimalCount   uint8
	_              [2]byte
	WorkAreaID     byte
	_              [2]byte
	FlagSetField   byte
	_              [7]byte
	IndexFieldFlag byte
}

func (f *FieldDescriptor) String() string {
	str := fmt.Sprintf("Name : %s\n", f.FieldName())
	str += fmt.Sprintf("Type : %c\n", f.FieldType)
	str += fmt.Sprintf("Len  : %d\n", f.FieldLength)
	str += fmt.Sprintf("Count: %d\n", f.DecimalCount)
	return str
}
func (f *FieldDescriptor) FieldName() string {
	for i, b := range f.FieldName_ {
		if b == '\000' {
			return (string)(f.FieldName_[0:i])
		}
	}
	return ""
}
