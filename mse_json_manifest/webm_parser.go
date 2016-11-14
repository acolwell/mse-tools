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
	"github.com/acolwell/mse-tools/ebml"
	"github.com/acolwell/mse-tools/webm"
	"time"
)

type webMClient struct {
	vcodec          string
	acodec          string
	timecodeScale   uint64
	duration        float64
	headerOffset    int64
	headerSize      int64
	clusterOffset   int64
	clusterSize     int64
	clusterTimecode uint64
	manifest        *JSONManifest
}

func (c *webMClient) OnListStart(offset int64, id int) bool {
	//fmt.Printf("OnListStart(%d, %s)\n", offset, webm.IdToName(id))

	if id == ebml.IdHeader {
		if c.headerSize != -1 {
			return false
		}
		c.headerOffset = offset
		c.headerSize = -1
		c.vcodec = ""
		c.acodec = ""
	} else if id == webm.IdCluster {
		if c.headerSize == -1 {
			c.headerSize = offset - c.headerOffset
			c.manifest.Init = &InitSegment{Offset: c.headerOffset, Size: c.headerSize}
		}
		c.clusterOffset = offset
	}
	return true
}

func (c *webMClient) OnListEnd(offset int64, id int) bool {
	//fmt.Printf("OnListEnd(%d, %s)\n", offset, webm.IdToName(id))
	scaleMult := float64(c.timecodeScale) / 1000000000.0

	if id == webm.IdInfo {
		if c.timecodeScale == 0 {
			c.timecodeScale = 1000000
		}
		if c.duration != -1 {
			c.manifest.Duration = c.duration * scaleMult
		}
		return true
	}

	if id == webm.IdTracks {
		contentType := ""
		if c.vcodec != "" && c.acodec != "" {
			contentType = fmt.Sprintf("video/webm; codecs=\"%s, %s\"", c.vcodec, c.acodec)
		} else if c.vcodec != "" && c.acodec == "" {
			contentType = fmt.Sprintf("video/webm; codecs=\"%s\"", c.vcodec)
		} else if c.vcodec == "" && c.acodec != "" {
			contentType = fmt.Sprintf("audio/webm; codecs=\"%s\"", c.acodec)
		}

		c.manifest.Type = contentType
		return true
	}

	if id == webm.IdCluster {
		c.manifest.Media = append(c.manifest.Media, &MediaSegment{
			Offset:   c.clusterOffset,
			Size:     (offset - c.clusterOffset),
			Timecode: (float64(c.clusterTimecode) * scaleMult),
		})
		return true
	}

	if id == webm.IdSegment {
		fmt.Printf(c.manifest.ToJSON())
	}
	return true
}

func (c *webMClient) OnBinary(id int, value []byte) bool {
	return true
}

func (c *webMClient) OnInt(id int, value int64) bool {
	return true
}

func (c *webMClient) OnUint(id int, value uint64) bool {
	if id == webm.IdTimecodeScale {
		c.timecodeScale = value
		return true
	}
	if id == webm.IdTimecode {
		c.clusterTimecode = value
		return true
	}
	if id == webm.IdDateUTC {
		c.manifest.StartDate = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(value))

		return true
	}
	return true
}

func (c *webMClient) OnFloat(id int, value float64) bool {
	if id == webm.IdDuration {
		c.manifest.Duration = value
	}
	return true
}

func (c *webMClient) OnString(id int, value string) bool {
	if id == webm.IdCodecID {
		switch value {
		case "V_VP8":
			c.vcodec = "vp8"
			break
		case "V_VP9":
			c.vcodec = "vp9"
			break
		case "A_VORBIS":
			c.acodec = "vorbis"
			break
		case "A_OPUS":
			c.acodec = "opus"
			break
		}
	}

	return true
}

func newWebMClient() *webMClient {
	return &webMClient{
		vcodec:          "",
		acodec:          "",
		timecodeScale:   0,
		duration:        -1,
		headerOffset:    -1,
		headerSize:      -1,
		clusterOffset:   -1,
		clusterTimecode: 0,
		manifest:        NewJSONManifest(),
	}
}

func NewWebMParser() *ebml.Parser {
	c := newWebMClient()

	return ebml.NewParser(ebml.GetListIDs(webm.IdTypes()), webm.UnknownSizeInfo(),
		ebml.NewElementParser(c, webm.IdTypes()))
}
