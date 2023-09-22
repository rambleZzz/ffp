package ipscan

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net"
	"os"
)

const (
	INDEX_LEN       = 7
	REDIRECT_MODE_1 = 0x01
	REDIRECT_MODE_2 = 0x02
)

type QQwry struct {
	db *os.File
}

type QQwryResult struct {
	IP      string
	Country string
	City    string
}

func NewQQwry(filepath string) (qqwry *QQwry, err error) {
	db, err := os.OpenFile(filepath, os.O_RDONLY, 0400)
	if err != nil {
		return
	}
	qqwry = &QQwry{db}
	return
}

func GbkToUtf8(s []byte) string {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

func (q *QQwry) ReadMode(offset uint32) byte {
	q.db.Seek(int64(offset), 0)
	mode := make([]byte, 1)
	q.db.Read(mode)
	return mode[0]
}

func (q *QQwry) ReadArea(offset uint32) []byte {
	mode := q.ReadMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		areaOffset := q.ReadUInt24()
		if areaOffset == 0 {
			return []byte("")
		} else {
			return q.ReadString(areaOffset)
		}
	}
	return q.ReadString(offset)
}

func (q *QQwry) ReadString(offset uint32) []byte {
	q.db.Seek(int64(offset), 0)
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		q.db.Read(buf)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

func (q *QQwry) SearchIndex(ip uint32) uint32 {
	header := make([]byte, 8)
	q.db.Seek(0, 0)
	q.db.Read(header)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	for {
		mid := q.GetMiddleOffset(start, end)
		q.db.Seek(int64(mid), 0)
		buf := make([]byte, INDEX_LEN)
		q.db.Read(buf)
		_ip := binary.LittleEndian.Uint32(buf[:4])

		if end-start == INDEX_LEN {
			offset := BytesToUInt32(buf[4:])
			q.db.Read(buf)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			} else {
				return 0
			}
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip {
			start = mid
		} else if _ip == ip {
			return BytesToUInt32(buf[4:])
		}
	}
}

func (q *QQwry) ReadUInt24() uint32 {
	buf := make([]byte, 3)
	q.db.Read(buf)
	return BytesToUInt32(buf)
}

func (q *QQwry) GetMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / INDEX_LEN) >> 1
	return start + records*INDEX_LEN
}

func BytesToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

func (q *QQwry) Find(ip string) (result QQwryResult, err error) {
	ipv4 := net.ParseIP(ip).To4()
	ipv4long := binary.BigEndian.Uint32(ipv4)
	offset := q.SearchIndex(ipv4long)
	if offset <= 0 {
		return
	}
	var country, area []byte
	mode := q.ReadMode(offset + 4)
	if mode == REDIRECT_MODE_1 {
		countryOffset := q.ReadUInt24()
		mode = q.ReadMode(countryOffset)
		if mode == REDIRECT_MODE_2 {
			c := q.ReadUInt24()
			country = q.ReadString(c)
			countryOffset += 4
		} else {
			country = q.ReadString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = q.ReadArea(countryOffset)
	} else if mode == REDIRECT_MODE_2 {
		countryOffset := q.ReadUInt24()
		country = q.ReadString(countryOffset)
		area = q.ReadArea(offset + 8)
	} else {
		country = q.ReadString(offset + 4)
		area = q.ReadArea(offset + uint32(5+len(country)))
	}
	result = QQwryResult{
		IP:      ip,
		Country: GbkToUtf8(country),
		City:    GbkToUtf8(area),
	}
	return result, nil
}

func (q *QQwry) Close() error {
	return q.db.Close()
}
