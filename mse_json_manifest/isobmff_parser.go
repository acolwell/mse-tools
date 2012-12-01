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

package main

import (
	"fmt"
	"github.com/acolwell/mse-tools/isobmff"
)

type isobmffClient struct {
	foundInitSegment   bool
	mediaSegmentOffset int64
	manifest           *JSONManifest
}

func (c *isobmffClient) OnHeader(offset int64, hdr []byte, id string, size int64) bool {
	fmt.Printf("OnHeader(%d, %s, %d)\n", offset, id, size)
	if offset == 0 && id != "ftyp" {
		fmt.Printf("File must start with a 'ftyp' box\n")
		return false
	}

	if id == "moov" {
		if c.foundInitSegment {
			fmt.Printf("Multiple 'moov' boxes not supported\n")
			return false
		}
	} else if id == "moof" {
		if !c.foundInitSegment {
			fmt.Printf("'moof' boxes must come after the 'moov' box.\n")
			return false
		}
		c.mediaSegmentOffset = offset
	} else if id == "mdat" {
		if c.mediaSegmentOffset == -1 {
			fmt.Printf("'mdat' boxes must come after the 'moof' box.\n")
			return false
		}
	}

	return true
}

func (c *isobmffClient) OnBody(offset int64, body []byte) bool {
	//fmt.Printf("OnBody(%d, %d)\n", offset, len(body))
	return true
}

func (c *isobmffClient) OnElementEnd(offset int64, id string) bool {
	fmt.Printf("OnElementEnd(%d, %s)\n", offset, id)

	if id == "moov" {
		c.foundInitSegment = true
		c.manifest.Init = &InitSegment{Offset: 0, Size: offset}
	} else if id == "mdat" {
		c.manifest.Media = append(c.manifest.Media, &MediaSegment{
			Offset:   c.mediaSegmentOffset,
			Size:     (offset - c.mediaSegmentOffset),
			Timecode: float64(-1),
		})

		c.mediaSegmentOffset = -1
	}
	return true
}

func (c *isobmffClient) OnEndOfData(offset int64) {
	fmt.Printf(c.manifest.ToJSON())
}

func newISOBMFFClient() *isobmffClient {
	return &isobmffClient{
		foundInitSegment:   false,
		mediaSegmentOffset: -1,
		manifest:           NewJSONManifest(),
	}
}

func NewISOBMFFParser() *isobmff.Parser {
	return isobmff.NewParser(newISOBMFFClient())
}
