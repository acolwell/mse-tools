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

type BlockInfo struct {
	Id         uint64
	Timecode   int
	Flags      uint8
	HeaderSize int
}

func ParseSimpleBlock(buf []byte) *BlockInfo {
	if len(buf) < 4 {
		return nil
	}

	// Compute trackID size.
	mask := byte(0x80)
	idSize := 1
	for ; (buf[0]&mask) == 0 && idSize < 8; idSize++ {
		mask >>= 1
	}

	// Check for invalid ID size.
	if idSize == 8 {
		return nil
	}

	headerSize := idSize + 3
	if len(buf) < headerSize {
		return nil
	}

	id := uint64(buf[0] & (mask - 1))
	for i := 1; i < idSize; i++ {
		id = (id << 8) | uint64(buf[i])
	}

	timecode := int(buf[idSize])<<8 | int(buf[idSize+1])
	if (timecode & 0x8000) != 0 {
		timecode |= (-1 << 16)
	}
	flags := buf[idSize+2]

	return &BlockInfo{Id: id, Timecode: timecode, Flags: flags, HeaderSize: headerSize}
}
