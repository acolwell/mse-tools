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

package isobmff

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ParserClient interface {
	OnHeader(offset int64, hdr []byte, id string, size int64) bool
	OnBody(offset int64, body []byte) bool
	OnElementEnd(offset int64, id string) bool
	OnEndOfData(offset int64)
}

type Parser struct {
	buf         *bytes.Buffer
	offset      int64
	bytesLeft   int64
	currentId   string
	client      ParserClient
	parserError bool
}

func (p *Parser) Append(buf []byte) bool {
	if p.parserError {
		return false
	}

	p.buf.Write(buf)

	for p.buf.Len() > 0 {
		if p.bytesLeft == 0 {
			totalParsed, id, size := p.readHeader(p.buf.Bytes())
			if totalParsed <= 0 {
				break
			}

			if !p.consumeHeader(totalParsed, id, size) {
				p.parserError = true
				return false
			}
			p.bytesLeft = size - int64(totalParsed)
		}

		bytesToConsume := int64(p.buf.Len())
		if p.bytesLeft <= bytesToConsume {
			bytesToConsume = p.bytesLeft
		}

		p.bytesLeft -= bytesToConsume
		if !p.consumeBody(int(bytesToConsume)) {
			p.parserError = true
			return false
		}
	}
	return true
}

func (p *Parser) EndOfData() {
	p.client.OnEndOfData(p.offset)
}

func (p *Parser) readHeader(buf []byte) (int, string, int64) {
	if len(buf) < 8 {
		return 0, "", 0
	}

	size := int64(binary.BigEndian.Uint32(buf[0:4]))
	id := bytes.NewBuffer(buf[4:8]).String()

	if size < 8 {
		fmt.Printf("Unsupported box size %d.\n", size)
		return -1, "", 0
	}

	if id == "uuid" {
		fmt.Printf("uuid boxes not supported.\n")
		return -1, "", 0
	}

	return 8, id, size
}

func (p *Parser) consumeHeader(headerSize int, id string, size int64) bool {
	p.currentId = id
	if !p.client.OnHeader(p.offset, p.buf.Next(headerSize), id, size) {
		return false
	}

	return p.consumeBytes(int64(headerSize))
}

func (p *Parser) consumeBody(byteCount int) bool {
	if byteCount > 0 && !p.client.OnBody(p.offset, p.buf.Next(byteCount)) {
		return false
	}

	if !p.consumeBytes(int64(byteCount)) {
		return false
	}

	if p.bytesLeft == 0 {
		if !p.client.OnElementEnd(p.offset, p.currentId) {
			return false
		}
	}
	return true
}

func (p *Parser) consumeBytes(byteCount int64) bool {
	if byteCount > 0 {
		p.offset += byteCount
	}

	return true
}

func NewParser(client ParserClient) *Parser {
	return &Parser{buf: bytes.NewBuffer([]byte{}), offset: 0, bytesLeft: 0, client: client, parserError: false}
}
