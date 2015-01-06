// Copyright 2015 Google Inc. All Rights Reserved.
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
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"fmt"
	"github.com/acolwell/mse-tools/ebml"
	"github.com/acolwell/mse-tools/webm"
	"io"
	"net/url"
	"os"
	"strings"
)

type TestClient struct {
	clusterTimecode uint64
	currentTrackId  uint64

	videoTrackId uint64
	codec4cc     uint32
	width        uint16
	height       uint16
	frameRate    float64
	timeScale    uint32
	frameCount   uint32
	out          io.WriteSeeker
}

func (c *TestClient) WriteHeader() {
	frameRate := uint32(c.frameRate * float64(c.timeScale))
	fmt.Printf("width %d height %d frameRate %d timeScale %d frameCount %d\n",
		c.width,
		c.height,
		frameRate,
		c.timeScale,
		c.frameCount)

	c.out.Seek(0, os.SEEK_SET)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(0x444b4946)) // 'DKIF' magic
	binary.Write(buf, binary.LittleEndian, uint16(0x0))        // version
	binary.Write(buf, binary.LittleEndian, uint16(32))         // header length
	binary.Write(buf, binary.BigEndian, c.codec4cc)
	binary.Write(buf, binary.LittleEndian, c.width)
	binary.Write(buf, binary.LittleEndian, c.height)
	binary.Write(buf, binary.LittleEndian, frameRate)
	binary.Write(buf, binary.LittleEndian, c.timeScale)
	binary.Write(buf, binary.LittleEndian, c.frameCount)
	binary.Write(buf, binary.LittleEndian, uint32(0x0)) // unused
	c.out.Write(buf.Bytes())
}

func (c *TestClient) OnListStart(offset int64, id int) bool {
	if id == webm.IdSegment {
		c.WriteHeader()
	}
	return true
}

func (c *TestClient) OnListEnd(offset int64, id int) bool {
	if id == webm.IdTrackEntry && c.videoTrackId == 0 {
		c.videoTrackId = c.currentTrackId
	} else if id == webm.IdSegment {
		c.WriteHeader()
	}
	return true
}

func (c *TestClient) OnBinary(id int, value []byte) bool {
	if id == webm.IdSimpleBlock {
		blockInfo := webm.ParseSimpleBlock(value)
		presentationTimecode := int64(c.clusterTimecode) + int64(blockInfo.Timecode)
		if blockInfo != nil {
			frameData := value[blockInfo.HeaderSize:];
			fmt.Printf("frame size %d timestamp %d\n", len(frameData), presentationTimecode)
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, uint32(len(frameData)))
			binary.Write(buf, binary.LittleEndian, uint64(presentationTimecode))
			binary.Write(buf, binary.BigEndian, frameData)
			c.out.Write(buf.Bytes())
			c.frameCount += 1
		} else {
			fmt.Printf("Invalid simple block")
			return false
		}
	}
	return true
}

func (c *TestClient) OnInt(id int, value int64) bool {
	return true
}

func (c *TestClient) OnUint(id int, value uint64) bool {
	if id == webm.IdTimecode {
		c.clusterTimecode = value
	}

	if c.videoTrackId == 0 {
		if id == webm.IdTrackNumber {
			c.currentTrackId = value
		} else if id == webm.IdPixelWidth {
			c.width = uint16(value)
		} else if id == webm.IdPixelHeight {
			c.height = uint16(value)
		} else if id == webm.IdTimecodeScale {
			c.timeScale = uint32(value)
		}
	}

	return true
}

func (c *TestClient) OnFloat(id int, value float64) bool {
	if c.videoTrackId == 0 && id == webm.IdFrameRate {
		c.frameRate = value
	}
	return true
}

func (c *TestClient) OnString(id int, value string) bool {
	if c.videoTrackId == 0 {
		if id == webm.IdCodecID {
			if value == "V_VP9" {
				c.codec4cc = 0x56503930
			} else if value == "V_VP8" {
				c.codec4cc = 0x56503830
			}
		}
	}
	return true
}

func NewTestClient(out io.WriteSeeker) *TestClient {
	return &TestClient{
		videoTrackId: 0,
		codec4cc:     0,
		width:        0,
		height:       0,
		frameRate:    0,
		timeScale:    0,
		frameCount:   0,
		out:          out,
	}
}

func checkError(str string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s - %s\n", str, err.Error())
		os.Exit(-1)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <infile> <outfile>\n", os.Args[0])
		return
	}

	var in io.Reader = nil
	if os.Args[1] == "-" {
		in = os.Stdin
	} else if strings.HasPrefix(os.Args[1], "ws://") {
		url, err := url.Parse(os.Args[1])
		checkError("Output url", err)

		origin := "http://localhost/"
		ws, err := websocket.Dial(url.String(), "", origin)
		checkError("WebSocket Dial", err)
		in = io.Reader(ws)
	} else {
		file, err := os.Open(os.Args[1])
		checkError(fmt.Sprintf("can't open file %s", os.Args[1]), err)
		in = io.Reader(file)
	}

	file, err := os.Create(os.Args[2])
	checkError(fmt.Sprintf("Failed to create file %s", os.Args[2]), err)
	out := io.WriteSeeker(file)

	buf := [4096]byte{}

	c := NewTestClient(out)

	parser := ebml.NewParser(ebml.GetListIDs(webm.IdTypes()), webm.UnknownSizeInfo(), ebml.NewElementParser(c, webm.IdTypes()))

	for {
		bytesRead, err := in.Read(buf[:])
		if err != nil {
			parser.EndOfData()
			break
		}

		parser.Append(buf[0:bytesRead])
	}
}
