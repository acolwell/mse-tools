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
)

const (
	VIDEO_TRACK int = 1
	AUDIO_TRACK int = 2
)

type Track interface {
	ID() uint64
	Type() int
	CodecID() string
}

type tracksParserClient struct {
	tracks      []Track
	trackNumber uint64
	trackType   int
	codecID     string
}

type track struct {
	id        uint64
	trackType int
	codecID   string
}

func (t *track) ID() uint64 {
	return t.id
}

func (t *track) Type() int {
	return t.trackType
}

func (t *track) CodecID() string {
	return t.codecID
}

func (p *tracksParserClient) Tracks() []Track {
	return p.tracks
}

func (p *tracksParserClient) OnListStart(offset int64, id int) bool {
	if id != IdTrackEntry {
		return false
	}

	p.trackNumber = 0
	p.trackType = 0
	p.codecID = ""

	return true
}

func (p *tracksParserClient) OnListEnd(offset int64, id int) bool {
	if id != IdTrackEntry {
		return false
	}

	p.tracks = append(p.tracks, &track{id: p.trackNumber, trackType: p.trackType, codecID: p.codecID})
	return true
}

func (p *tracksParserClient) OnBinary(id int, value []byte) bool {
	return true
}

func (p *tracksParserClient) OnInt(id int, value int64) bool {
	return false
}

func (p *tracksParserClient) OnUint(id int, value uint64) bool {
	if id == IdTrackNumber {
		p.trackNumber = value
		return true
	}

	if id == IdTrackType {
		p.trackType = int(value)
		return true
	}

	return false
}

func (p *tracksParserClient) OnFloat(id int, value float64) bool {
	return false
}

func (p *tracksParserClient) OnString(id int, value string) bool {
	if id == IdCodecID {
		p.codecID = value
		return true
	}
	return false
}

func ParseTracksElement(buf []byte) []Track {
	typeInfo := map[int]int{
		IdTrackEntry:  ebml.TypeList,
		IdTrackNumber: ebml.TypeUint,
		IdTrackType:   ebml.TypeUint,
		IdCodecID:     ebml.TypeString}

	client := &tracksParserClient{
		tracks:      []Track{},
		trackNumber: 0,
		trackType:   0,
		codecID:     ""}
	parser := ebml.NewParser(ebml.GetListIDs(typeInfo), map[int][]int{},
		ebml.NewElementParser(client, typeInfo))

	if !parser.Append(buf) {
		log.Printf("Failed to parse tracks.")
		return nil
	}

	return client.Tracks()
}
