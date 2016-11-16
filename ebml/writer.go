// Copyright 2012 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebml

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	UNKNOWN_SIZE int64 = 0xffffffffffffff
)

type writerListInfo struct {
	id           int
	headerOffset int64
	bodyOffset   int64
}

type Writer struct {
	offset   int64
	writer   io.Writer
	seeker   io.Seeker
	listInfo []writerListInfo
}

func NewWriter(writerSeeker io.WriteSeeker) *Writer {
	return &Writer{offset: 0, writer: io.Writer(writerSeeker), seeker: io.Seeker(writerSeeker), listInfo: []writerListInfo{}}
}

func NewNonSeekableWriter(writer io.Writer) *Writer {
	return &Writer{offset: 0, writer: writer, seeker: nil, listInfo: []writerListInfo{}}
}

func (w *Writer) WriteUnknownSizeHeader(id int) (int, error) {
	return w.writeHeader(id, UNKNOWN_SIZE)
}

func (w *Writer) Write(id int, data interface{}) (int, error) {
	switch v := data.(type) {
	case uint8:
		return w.writeUInt64(id, uint64(v))
	case uint16:
		return w.writeUInt64(id, uint64(v))
	case uint32:
		return w.writeUInt64(id, uint64(v))
	case uint64:
		return w.writeUInt64(id, v)
	case int:
		return w.writeInt64(id, int64(v))
	case int8:
		return w.writeInt64(id, int64(v))
	case int16:
		return w.writeInt64(id, int64(v))
	case int32:
		return w.writeInt64(id, int64(v))
	case int64:
		return w.writeInt64(id, int64(v))
	case float32:
	case float64:
		return w.writeFloat(id, v)
	case string:
		return w.writeBinary(id, []byte(v))
	case []byte:
		return w.writeBinary(id, v)
	}

	panic(fmt.Sprintf("Unexpected type %s", data))
	return 0, errors.New("Unexpected type")
}

func (w *Writer) WriteListStart(id int) {
	headerOffset := w.Offset()
	if _, err := w.WriteUnknownSizeHeader(id); err != nil {
		panic(fmt.Sprintf("Failed to write header. err=%s", err.Error()))
	}
	bodyOffset := w.Offset()

	w.listInfo = append(w.listInfo, writerListInfo{id: id, headerOffset: headerOffset, bodyOffset: bodyOffset})
}

func (w *Writer) WriteListEnd(id int) {
	currentOffset := w.Offset()

	rewroteHeaders := false
	for {
		li := w.listInfo[len(w.listInfo)-1]
		w.listInfo = w.listInfo[:len(w.listInfo)-1]

		if w.seeker != nil {
			if _, err := w.seeker.Seek(li.headerOffset, os.SEEK_SET); err == nil {
				w.offset = li.headerOffset
				rewroteHeaders = true
				if _, err = w.writeHeader8(li.id, currentOffset-li.bodyOffset); err != nil {
					panic(fmt.Sprintf("Header rewrite failed. err=%s", err.Error()))
				}
			}
		}

		if li.id == id {
			break
		}
	}

	if rewroteHeaders {
		_, err := w.seeker.Seek(currentOffset, os.SEEK_SET)
		if err != nil {
			panic(fmt.Sprintf("Seek back to original offset failed. err=%s", err.Error()))
		}
		w.offset = currentOffset
	}
}

func (w *Writer) WriteVoid(size int) (int, error) {
	if size < 2 {
		panic("Can't void a space smaller than 2 bytes.")
	}

	var total = 0
	var err error = nil
	voidSize := size
	if size < 9 {
		voidSize -= 2
		total, err = w.writeHeader(IdVoid, int64(voidSize))
	} else {
		voidSize -= 9
		total, err = w.writeHeader8(IdVoid, int64(voidSize))
	}

	if err != nil || voidSize <= 0 {
		return total, err
	}

	n, err := w.writeToOutput(make([]byte, voidSize))
	return total + n, err
}

func (w *Writer) CanSeek() bool {
	return w.seeker != nil
}

func (w *Writer) Offset() int64 {
	if w.seeker != nil {
		offset, err := w.seeker.Seek(0, os.SEEK_CUR)
		if err == nil && offset != w.offset {
			panic(fmt.Sprintf("Offset mismatch %d %d\n", offset, w.offset))
		}
	}
	return w.offset
}

func (w *Writer) SetOffset(offset int64) bool {
	if w.seeker == nil {
		return false
	}

	if _, err := w.seeker.Seek(offset, os.SEEK_SET); err != nil {
		return false
	}
	w.offset = offset
	return true
}

func (w *Writer) WriteToOutput(p []byte) (int, error) {
	return w.writeToOutput(p)
}

func (w *Writer) writeToOutput(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err == nil {
		w.offset += int64(n)
	}
	return n, err
}

func (w *Writer) writeBinary(id int, body []byte) (int, error) {
	header_bytes, err := w.writeHeader(id, int64(len(body)))
	if err != nil {
		return header_bytes, err
	}

	body_bytes, err := w.writeToOutput(body)
	return header_bytes + body_bytes, err
}

func (w *Writer) writeInt64(id int, value int64) (int, error) {
	buf := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
	count := 0

	if value == 0x80000000 {
		count = 7
	} else {
		var tmp uint64
		var threshold uint64
		if value > 0 {
			tmp = uint64(value)
			threshold = 0x80
		} else {
			tmp = uint64(-value)
			threshold = 0x100
		}

		for ; tmp >= threshold && count < 7; count++ {
			threshold <<= 8
		}
	}

	for i := count; i >= 0; i-- {
		buf[i] = byte(value & 0xff)
		value >>= 8
	}
	return w.writeBinary(id, buf[:count+1])
}

func (w *Writer) writeUInt64(id int, value uint64) (int, error) {
	buf := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
	count := 0
	var threshold uint64 = 0x100

	for ; value >= threshold && count < 7; count++ {
		threshold <<= 8
	}

	for i := count; i >= 0; i-- {
		buf[i] = byte(value & 0xff)
		value >>= 8
	}
	return w.writeBinary(id, buf[:count+1])
}

func (w *Writer) writeFloat(id int, value interface{}) (int, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		return 0, err
	}
	return w.writeBinary(id, buf.Bytes())
}

func (w *Writer) writeHeader(id int, size int64) (int, error) {
	id_bytes, err := w.writeId(id)
	if err != nil {
		return id_bytes, err
	}

	size_bytes, err := w.writeSize(size, false)
	return id_bytes + size_bytes, err
}

func (w *Writer) writeHeader8(id int, size int64) (int, error) {
	id_bytes, err := w.writeId(id)
	if err != nil {
		return id_bytes, err
	}

	size_bytes, err := w.writeSize(size, true)
	return id_bytes + size_bytes, err
}

func (w *Writer) writeId(id int) (int, error) {
	buf := [4]byte{0, 0, 0, 0}
	count := 0
	mask := 0xff

	for ; id > mask && count < 3; count++ {
		mask = (mask << 7) | 0x7f
	}

	for i := count; i >= 0; i-- {
		buf[i] = byte(id & 0xff)
		id >>= 8
	}

	return w.writeToOutput(buf[:count+1])
}

func (w *Writer) writeSize(size int64, use8bytes bool) (int, error) {
	buf := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
	count := 0
	var mask int64 = 0x7f
	var sizeFlag int64 = 0x80

	if use8bytes {
		count = 7
		mask = 0xffffffffffffff
		sizeFlag = 0x01
	} else {
		for ; size > (mask-1) && count < 7; count++ {
			mask = (mask << 7) | 0x7f
			sizeFlag >>= 1
		}
	}
	for i := count; i > 0; i-- {
		buf[i] = byte(size & 0xff)
		size >>= 8
	}
	buf[0] = byte(sizeFlag | (size & (sizeFlag - 1)))
	return w.writeToOutput(buf[:count+1])
}
