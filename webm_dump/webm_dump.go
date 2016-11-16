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
	"golang.org/x/net/websocket"
	"io"
	"net/url"
	"os"
	"strings"
)

type TestClient struct {
	depth           int
	clusterTimecode uint64
}

func (c *TestClient) OnListStart(offset int64, id int) bool {
	fmt.Printf("%s<%s type=\"list\" offset=\"%d\">\n", c.indent(), webm.IdToName(id), offset)
	c.depth++
	return true
}

func (c *TestClient) OnListEnd(offset int64, id int) bool {
	c.depth--
	fmt.Printf("%s</%s>\n", c.indent(), webm.IdToName(id))
	return true
}

func (c *TestClient) OnBinary(id int, value []byte) bool {
	if id != webm.IdSimpleBlock && id != webm.IdBlock {
		fmt.Printf("%s<%s type=\"binary\" size=\"%d\"/>\n", c.indent(), webm.IdToName(id), len(value))
	} else {
		blockInfo := webm.ParseSimpleBlock(value)
		presentationTimecode := int64(c.clusterTimecode) + int64(blockInfo.Timecode)
		if blockInfo != nil {
			fmt.Printf("%s<%s type=\"binary\" size=\"%d\" trackNum=\"%d\" timecode=\"%d\" presentationTimecode=\"%d\" flags=\"%x\"/>\n",
				c.indent(), webm.IdToName(id), len(value), blockInfo.Id, blockInfo.Timecode, presentationTimecode, blockInfo.Flags)
		} else {
			fmt.Printf("%s<%s type=\"binary\" size=\"%d\" invalid=\"true\"/>\n", c.indent(), webm.IdToName(id), len(value))
		}
	}
	return true
}

func (c *TestClient) OnInt(id int, value int64) bool {
	fmt.Printf("%s<%s type=\"int\" value=\"%d\"/>\n", c.indent(), webm.IdToName(id), value)
	return true
}

func (c *TestClient) OnUint(id int, value uint64) bool {
	if id != webm.IdSeekID {
		fmt.Printf("%s<%s type=\"uint\" value=\"%d\"/>\n", c.indent(), webm.IdToName(id), value)
	} else {
		fmt.Printf("%s<%s type=\"uint\" id_name=\"%s\" value=\"%d\"/>\n", c.indent(), webm.IdToName(id), webm.IdToName(int(value)), value)
	}
	if id == webm.IdTimecode {
		c.clusterTimecode = value
	}
	return true
}

func (c *TestClient) OnFloat(id int, value float64) bool {
	fmt.Printf("%s<%s type=\"float\" value=\"%f\"/>\n", c.indent(), webm.IdToName(id), value)
	return true
}

func (c *TestClient) OnString(id int, value string) bool {
	fmt.Printf("%s<%s type=\"string\" value=\"%s\"/>\n", c.indent(), webm.IdToName(id), value)
	return true
}

func (c *TestClient) indent() string {
	return strings.Repeat("  ", c.depth)
}

func NewTestClient() *TestClient {
	return &TestClient{depth: 0}
}

func checkError(str string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s - %s\n", str, err.Error())
		os.Exit(-1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <infile>\n", os.Args[0])
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

	buf := [1024]byte{}

	c := NewTestClient()

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
