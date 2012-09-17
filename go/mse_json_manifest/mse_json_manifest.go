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
	"github.com/acolwell/mse-tools/go/ebml"
	"github.com/acolwell/mse-tools/go/webm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type TestClient struct {
	hasVorbis       bool
	hasVP8          bool
	timecodeScale   uint64
	duration        float64
	startDate       time.Time
	headerOffset    int64
	headerSize      int64
	clusterOffset   int64
	clusterSize     int64
	clusterTimecode uint64
}

func (c *TestClient) OnListStart(offset int64, id int) bool {
	//log.Printf("OnListStart(%d, %s)\n", offset, webm.IdToName(id))

	if id == ebml.IdHeader {
		if c.headerSize != -1 {
			fmt.Printf("},\n")
		}
		fmt.Printf("{\n")
		c.headerOffset = offset
		c.headerSize = -1
		c.hasVorbis = false
		c.hasVP8 = false
	} else if id == webm.IdCluster {
		if c.headerSize == -1 {
			c.headerSize = offset - c.headerOffset
			fmt.Printf("  init: { offset: %d, size: %d},\n",
				c.headerOffset,
				c.headerSize)
			fmt.Printf("  media: [\n")
		}
		c.clusterOffset = offset
	}
	return true
}

func (c *TestClient) OnListEnd(offset int64, id int) bool {
	//log.Printf("OnListEnd(%d, %s)\n", offset, webm.IdToName(id))
	scaleMult := float64(c.timecodeScale) / 1000000000.0

	if id == webm.IdInfo {
		if c.timecodeScale == 0 {
			return false
		}
		if c.duration == -1 {
			fmt.Printf("  live: true, \n")
		} else {
			fmt.Printf("  duration: %f,\n", c.duration*scaleMult)
		}
		if !c.startDate.IsZero() {
			fmt.Printf("  startDate: '%s', \n", c.startDate.Format(time.RFC3339Nano))
		}
		return true
	}

	if id == webm.IdTracks {
		contentType := ""
		if c.hasVP8 {
			contentType = "video/webm; codecs=\"vp8"
			if c.hasVorbis {
				contentType += ", vorbis"
			}
			contentType += "\""
		} else if c.hasVorbis {
			contentType = "audio/webm; codecs=\"vorbis\""
		}

		fmt.Printf("  type: '%s',\n", contentType)
		return true
	}

	if id == webm.IdCluster {
		fmt.Printf("    { offset: %d, size: %d, timecode: %f },\n",
			c.clusterOffset,
			offset-c.clusterOffset,
			float64(c.clusterTimecode)*scaleMult)
		return true
	}
	return true
}

func (c *TestClient) OnBinary(id int, value []byte) bool {
	return true
}

func (c *TestClient) OnInt(id int, value int64) bool {
	return true
}

func (c *TestClient) OnUint(id int, value uint64) bool {
	if id == webm.IdTimecodeScale {
		c.timecodeScale = value
		return true
	}
	if id == webm.IdTimecode {
		c.clusterTimecode = value
		return true
	}
	if id == webm.IdDateUTC {
		c.startDate = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(value))
		return true
	}
	return true
}

func (c *TestClient) OnFloat(id int, value float64) bool {
	if id == webm.IdDuration {
		c.duration = value
	}
	return true
}

func (c *TestClient) OnString(id int, value string) bool {
	if id == webm.IdCodecID {
		switch value {
		case "V_VP8":
			c.hasVP8 = true
			break
		case "A_VORBIS":
			c.hasVorbis = true
			break
		}
	}

	return true
}

func NewTestClient() *TestClient {
	return &TestClient{
		hasVorbis:       false,
		hasVP8:          false,
		timecodeScale:   0,
		duration:        -1,
		headerOffset:    -1,
		headerSize:      -1,
		clusterOffset:   -1,
		clusterTimecode: 0,
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <infile>\n", os.Args[0])
		return
	}

	var in io.Reader = nil
	var err error = nil

	if os.Args[1] == "-" {
		in = os.Stdin
	} else if strings.HasPrefix(os.Args[1], "http://") {
		resp, err := http.Get(os.Args[1])
		if err != nil {
			log.Printf("can't open url; err=%s\n", err.Error())
			os.Exit(1)
		}
		in = resp.Body
	} else {
		in, err = os.Open(os.Args[1])
		if in == nil {
			log.Printf("can't open file; err=%s\n", err.Error())
			os.Exit(1)
		}
	}

	buf := [1024]byte{}

	c := NewTestClient()

	parser := ebml.NewParser(ebml.GetListIDs(webm.IdTypes()), webm.UnknownSizeInfo(),
		ebml.NewElementParser(c, webm.IdTypes()))

	for done := false; !done; {
		bytesRead, err := in.Read(buf[:])
		if err == io.EOF || err == io.ErrClosedPipe {
			done = true

			fmt.Printf("  ]\n")
			fmt.Printf("}\n")
		} else {
			parser.Append(buf[0:bytesRead])
		}
	}
}
