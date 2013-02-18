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
	"log"
)

type Header interface {
	Version() uint64
	ReadVersion() uint64
	MaxIDLength() uint64
	MaxSizeLength() uint64
	DocType() string
	DocTypeVersion() uint64
	DocTypeReadVersion() uint64
}

type parserClient struct {
	version            uint64
	readVersion        uint64
	maxIDLength        uint64
	maxSizeLength      uint64
	docType            string
	docTypeVersion     uint64
	docTypeReadVersion uint64
}

func (p *parserClient) Version() uint64 {
	return p.version
}

func (p *parserClient) ReadVersion() uint64 {
	return p.readVersion
}

func (p *parserClient) MaxIDLength() uint64 {
	return p.maxIDLength
}

func (p *parserClient) MaxSizeLength() uint64 {
	return p.maxSizeLength
}

func (p *parserClient) DocType() string {
	return p.docType
}

func (p *parserClient) DocTypeVersion() uint64 {
	return p.docTypeVersion
}

func (p *parserClient) DocTypeReadVersion() uint64 {
	return p.docTypeReadVersion
}

func (p *parserClient) OnListStart(offset int64, id int) bool {
	return false
}

func (p *parserClient) OnListEnd(offset int64, id int) bool {
	return false
}

func (p *parserClient) OnBinary(id int, value []byte) bool {
	return id == IdCRC32 || id == IdVoid
}

func (p *parserClient) OnInt(id int, value int64) bool {
	return false
}

func (p *parserClient) OnUint(id int, value uint64) bool {
	if id == IdVersion {
		p.version = value
		return true
	}
	if id == IdReadVersion {
		p.readVersion = value
		return true
	}
	if id == IdMaxIDLength {
		p.maxIDLength = value
		return true
	}
	if id == IdMaxSizeLength {
		p.maxSizeLength = value
		return true
	}
	if id == IdDocTypeVersion {
		p.docTypeVersion = value
		return true
	}
	if id == IdDocTypeReadVersion {
		p.docTypeReadVersion = value
		return true
	}

	return false
}

func (p *parserClient) OnFloat(id int, value float64) bool {
	return false
}

func (p *parserClient) OnString(id int, value string) bool {
	if id != IdDocType {
		return false
	}

	p.docType = value
	return true
}

func ParseHeader(buf []byte) Header {
	typeInfo := map[int]int{
		IdVersion:            TypeUint,
		IdReadVersion:        TypeUint,
		IdMaxIDLength:        TypeUint,
		IdMaxSizeLength:      TypeUint,
		IdDocType:            TypeString,
		IdDocTypeVersion:     TypeUint,
		IdDocTypeReadVersion: TypeUint}

	client := &parserClient{
		version:            1,
		readVersion:        1,
		maxIDLength:        4,
		maxSizeLength:      8,
		docType:            "",
		docTypeVersion:     1,
		docTypeReadVersion: 1}
	parser := NewParser(GetListIDs(typeInfo), map[int][]int{},
		NewElementParser(client, typeInfo))

	if !parser.Append(buf) {
		log.Printf("Failed to parse header.")
		return nil
	}

	if client.Version() != 1 {
		log.Printf("Unsupported EBML Version %d", client.Version())
		return nil
	}

	if client.ReadVersion() != 1 {
		log.Printf("Unsupported EBML ReadVersion %d", client.ReadVersion())
		return nil
	}

	if client.MaxIDLength() > 4 {
		log.Printf("Unsupported EBML MaxIDLength %d", client.MaxIDLength())
		return nil

	}

	if client.MaxSizeLength() > 8 {
		log.Printf("Unsupported EBML MaxSizeLength %d", client.MaxSizeLength())
		return nil
	}

	if client.DocType() == "" {
		log.Printf("Empty EBML DocType not supported")
		return nil
	}

	if client.DocTypeVersion() < 1 {
		log.Printf("Unsupported EBML DocTypeVersion %d", client.DocTypeVersion())
		return nil
	}

	if client.DocTypeReadVersion() < 1 {
		log.Printf("Unsupported EBML DocTypeReadVersion %d", client.DocTypeReadVersion())
		return nil
	}

	return client
}
