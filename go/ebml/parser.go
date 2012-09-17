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
	//"log"
)

type ParserClient interface {
	OnHeader(offset int64, hdr []byte, id int, size int64) bool
	OnBody(offset int64, body []byte) bool
	OnElementEnd(offset int64, id int) bool
}

type listInfo struct {
	id          int
	size        int64
	bytesParsed int64
}

type Parser struct {
	buf              *bytes.Buffer
	offset           int64
	bytesLeft        int64
	currentId        int
	lists            []*listInfo
	client           ParserClient
	listMap          map[int]bool
	unknownSizeIdMap map[int]map[int]bool
	parserError      bool
}

func (li *listInfo) AddBytes(byteCount int64) bool {
	li.bytesParsed += byteCount
	return li.bytesParsed == li.size
}

func NewParser(listElements []int, unknownSizeInfo map[int][]int, client ParserClient) *Parser {
	listMap := map[int]bool{}
	for i := 0; i < len(listElements); i++ {
		listMap[listElements[i]] = true
	}
	unknownSizeIdMap := map[int]map[int]bool{}
	for id, parentAndPeerIds := range unknownSizeInfo {
		idMap := map[int]bool{}

		for i := range parentAndPeerIds {
			parentOrPeerId := parentAndPeerIds[i]
			idMap[parentOrPeerId] = true

			if parentOrPeerId != id {
				for j := range unknownSizeInfo[parentOrPeerId] {
					ancestorId := unknownSizeInfo[parentOrPeerId][j]
					idMap[ancestorId] = true
				}
			}
		}
		unknownSizeIdMap[id] = idMap
	}

	return &Parser{buf: bytes.NewBuffer([]byte{}), offset: 0, bytesLeft: 0, client: client, listMap: listMap, unknownSizeIdMap: unknownSizeIdMap, parserError: false}
}

func (b *Parser) Append(buf []byte) bool {
	if b.parserError {
		return false
	}

	b.buf.Write(buf)

	for b.buf.Len() > 0 {
		if b.bytesLeft == 0 {
			totalParsed, id, size := b.readHeader(b.buf.Bytes())
			if totalParsed <= 0 {
				break
			}

			// Check to see if this ID indicates the end of
			// a list with an unknown size.
			if !b.checkForAncestorId(id) {
				b.parserError = true
				return false
			}

			//log.Printf("%d id %s size %d depth %d\n",
			//	b.offset, idToName[id], size,
			//	len(b.lists))

			if b.isList(id) {
				if _, ok := b.unknownSizeIdMap[id]; (size == -1) && !ok {
					b.parserError = true
					return false
				}

				// Consume the header.
				if !b.consumeHeader(totalParsed, id, size) {
					b.parserError = true
					return false
				}

				b.lists = append(b.lists, &listInfo{id: id, size: size, bytesParsed: 0})
				continue
			}

			if size == -1 {
				b.parserError = true
				return false
			}

			// Consume the header.
			if !b.consumeHeader(totalParsed, id, size) {
				b.parserError = true
				return false
			}
			b.bytesLeft = size
		}

		bytesToConsume := b.buf.Len()
		if b.bytesLeft <= int64(bytesToConsume) {
			bytesToConsume = int(b.bytesLeft)
		}

		// Consume element body.
		b.bytesLeft -= int64(bytesToConsume)
		if !b.consumeBody(bytesToConsume) {
			b.parserError = true
			return false
		}
	}
	return true
}

func (b *Parser) EndOfData() {
	for len(b.lists) > 0 {
		li := b.lists[len(b.lists)-1]
		if li.size != -1 {
			break
		}

		li.size = li.bytesParsed
		if !b.consumeBytes(0) {
			return
		}
	}
}

func (b *Parser) checkForAncestorId(id int) bool {
	for len(b.lists) > 0 {
		li := b.lists[len(b.lists)-1]
		if li.size != -1 {
			break
		}

		ancestors := b.unknownSizeIdMap[li.id]
		if _, isAncestor := ancestors[id]; !isAncestor {
			break
		}

		li.size = li.bytesParsed
		if !b.consumeBytes(0) {
			return false
		}
	}
	return true
}

func (b *Parser) consumeHeader(headerSize int, id int, size int64) bool {
	b.currentId = id
	if !b.client.OnHeader(b.offset, b.buf.Next(headerSize), id, size) {
		return false
	}

	return b.consumeBytes(headerSize)
}

func (b *Parser) consumeBody(byteCount int) bool {
	if byteCount > 0 && !b.client.OnBody(b.offset, b.buf.Next(byteCount)) {
		return false
	}

	if b.bytesLeft == 0 {
		if !b.client.OnElementEnd(b.offset, b.currentId) {
			return false
		}

		if len(b.lists) > 0 {
			b.currentId = b.lists[len(b.lists)-1].id
		}
	}
	return b.consumeBytes(byteCount)
}

func (b *Parser) consumeBytes(byteCount int) bool {
	if byteCount > 0 {
		b.offset += int64(byteCount)
	}

	listByteCount := int64(byteCount)
	for len(b.lists) > 0 {
		li := b.lists[len(b.lists)-1]

		if !li.AddBytes(listByteCount) {
			break
		}

		if !b.client.OnElementEnd(b.offset, li.id) {
			return false
		}

		//log.Printf("list end %s\n", idToName[li.id]);
		listByteCount = li.size
		b.lists = b.lists[:len(b.lists)-1]
	}

	return true
}

func (b *Parser) readNumber(buf []byte, isSize bool) (int, int64) {
	if len(buf) < 1 {
		return 0, 0
	}

	maxSize := 4
	if isSize {
		maxSize = 8
	}

	totalBytes := 0
	mask := uint8(0x80)
	for i := 0; i < maxSize; i++ {

		if (buf[0] & mask) != 0 {
			totalBytes = 1 + i
			mask = ^mask
			break
		}
		mask >>= 1
		mask |= 0x80
	}

	if totalBytes == 0 {
		return -1, 0
	}

	if len(buf) < totalBytes {
		return 0, 0
	}

	allOnes := (buf[0] == 0xff)
	value := int64(buf[0])
	if isSize {
		value &= int64(mask)
		allOnes = (buf[0] & mask) == mask
	}

	if totalBytes == 1 {
		return totalBytes, value
	}

	for i := 1; i < totalBytes; i++ {
		value <<= 8
		value |= int64(buf[i])
		if buf[i] != 0xff {
			allOnes = false
		}
	}

	if allOnes {
		value = -1
	}

	return totalBytes, value
}

func (b *Parser) isList(id int) bool {
	return b.listMap[id]
}

func (b *Parser) readHeader(buf []byte) (int, int, int64) {
	idBytes, id := b.readNumber(buf, false)
	if idBytes <= 0 {
		return idBytes, 0, 0
	}

	sizeBuf := buf[idBytes:]
	sizeBytes, size := b.readNumber(sizeBuf, true)
	if sizeBytes <= 0 {
		return sizeBytes, 0, 0
	}

	return idBytes + sizeBytes, int(id), size
}
