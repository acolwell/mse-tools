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
)

type ElementParserClient interface {
	OnListStart(offset int64, id int) bool
	OnListEnd(offset int64, id int) bool
	OnBinary(id int, value []byte) bool
	OnInt(id int, value int64) bool
	OnUint(id int, value uint64) bool
	OnFloat(id int, value float64) bool
	OnString(id int, value string) bool
}

type ElementParser struct {
	id      int
	buf     *bytes.Buffer
	client  ElementParserClient
	typeMap map[int]int
}

func (p *ElementParser) OnHeader(offset int64, hdr []byte, id int, size int64) bool {
	p.id = id
	p.buf.Truncate(0)

	if elementType, present := p.typeMap[p.id]; present && elementType == TypeList {
		return p.client.OnListStart(offset, id)
	}
	return true
}

func (p *ElementParser) OnBody(offset int64, body []byte) bool {
	bytesWritten, err := p.buf.Write(body)
	return (bytesWritten == len(body)) && (err == nil)
}

func (p *ElementParser) OnElementEnd(offset int64, id int) bool {
	if elementType, present := p.typeMap[id]; present {
		switch elementType {
		case TypeList:
			return p.client.OnListEnd(offset, id)
		case TypeBinary:
			return p.ParseBinary(p.id, p.buf.Bytes())
			break
		case TypeUint:
			return p.ParseUint(p.id, p.buf.Bytes())
			break
		case TypeInt:
			return p.ParseInt(p.id, p.buf.Bytes())
			break
		case TypeFloat:
			return p.ParseFloat(p.id, p.buf.Bytes())
			break
		case TypeString:
			return p.ParseString(p.id, p.buf.Bytes())
			break
		case TypeUTF8:
			return p.ParseUTF8(p.id, p.buf.Bytes())
			break
		}
	}
	return p.ParseBinary(p.id, p.buf.Bytes())
}

func (p *ElementParser) ParseBinary(id int, body []byte) bool {
	return p.client.OnBinary(id, body)
}

func (p *ElementParser) ParseUint(id int, body []byte) bool {
	if len(body) == 0 || len(body) > 8 {
		return false
	}
	var value uint64 = 0
	for i := 0; i < len(body); i += 1 {
		value = (value << 8) | uint64(body[i])
	}
	return p.client.OnUint(id, value)
}

func (p *ElementParser) ParseInt(id int, body []byte) bool {
	if len(body) == 0 || len(body) > 8 {
		return false
	}

	var value int64 = 0
	if body[0]&0x80 != 0 {
		value = -1
	}

	for i := 0; i < len(body); i += 1 {
		value = (value << 8) | int64(body[i])
	}
	return p.client.OnInt(id, value)
}

func (p *ElementParser) ParseFloat(id int, body []byte) bool {
	var buf = bytes.NewBuffer(body)
	if len(body) == 4 {
		var value float32
		err := binary.Read(buf, binary.BigEndian, &value)
		if err != nil {
			return false
		}
		return p.client.OnFloat(id, float64(value))
	} else if len(body) == 8 {
		var value float64
		err := binary.Read(buf, binary.BigEndian, &value)
		if err != nil {
			return false
		}
		return p.client.OnFloat(id, value)
	}
	return false
}

func (p *ElementParser) ParseString(id int, body []byte) bool {
	return p.client.OnString(id, string(body))
}

func (p *ElementParser) ParseUTF8(id int, body []byte) bool {
	return p.client.OnString(id, string(body))
}

func NewElementParser(client ElementParserClient, typeMap map[int]int) *ElementParser {
	return &ElementParser{
		id: -1, buf: bytes.NewBuffer([]byte{}), client: client, typeMap: typeMap}
}
