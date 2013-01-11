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

package webm

import (
	"github.com/acolwell/mse-tools/ebml"
	"log"
	"math"
)

type InfoElement interface {
	TimecodeScale() uint64
	Duration() float64
	Date() int64
}

type infoParserClient struct {
	timecodeScale uint64
	duration      float64
	date          int64
}

func (p *infoParserClient) TimecodeScale() uint64 {
	return p.timecodeScale
}

func (p *infoParserClient) Duration() float64 {
	return p.duration
}

func (p *infoParserClient) Date() int64 {
	return p.date
}

func (p *infoParserClient) OnListStart(offset int64, id int) bool {
	return false
}

func (p *infoParserClient) OnListEnd(offset int64, id int) bool {
	return false
}

func (p *infoParserClient) OnBinary(id int, value []byte) bool {
	switch id {
	case ebml.IdCRC32,
		ebml.IdVoid,
		IdSegmentUID,
		IdSegmentFilename,
		IdPrevUID,
		IdPrevFilename,
		IdNextUID,
		IdNextFilename,
		IdSegmentFamily,
		IdChapterTranslate,
		IdChapterTranslateEditionUID,
		IdChapterTranslateCodec,
		IdChapterTranslateID,
		IdTitle,
		IdMuxingApp,
		IdWritingApp:
		return true
	}
	return false
}

func (p *infoParserClient) OnInt(id int, value int64) bool {
	if id != IdDateUTC {
		return false
	}
	p.date = value
	return true
}

func (p *infoParserClient) OnUint(id int, value uint64) bool {
	if id == IdTimecodeScale {
		p.timecodeScale = value
		return true
	}

	return false
}

func (p *infoParserClient) OnFloat(id int, value float64) bool {
	if id != IdDuration {
		return false
	}
	p.duration = value
	return true
}

func (p *infoParserClient) OnString(id int, value string) bool {
	return false
}

func ParseInfoElement(buf []byte) InfoElement {
	typeInfo := map[int]int{
		IdTimecodeScale: ebml.TypeUint,
		IdDuration:      ebml.TypeFloat,
		IdDateUTC:       ebml.TypeInt}

	client := &infoParserClient{
		timecodeScale: 1000000,
		duration:      math.Inf(1),
		date:          0}
	parser := ebml.NewParser(ebml.GetListIDs(typeInfo), map[int][]int{},
		ebml.NewElementParser(client, typeInfo))

	if !parser.Append(buf) {
		log.Printf("Failed to parse info.")
		return nil
	}

	if client.TimecodeScale() == 0 ||
		client.Duration() <= 0 {
		return nil
	}

	return client
}
