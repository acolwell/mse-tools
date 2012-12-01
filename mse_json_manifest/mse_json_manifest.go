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
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

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

	buf := [4096]byte{}

	var parser Parser = nil
	for done := false; !done; {
		bytesRead, err := in.Read(buf[:])
		if err == io.EOF || err == io.ErrClosedPipe {
			done = true
			continue
		}

		if parser == nil {
			if len(buf) < 8 {
				log.Printf("Not enough bytes to detect file type.\n")
				break
			} else if binary.BigEndian.Uint32(buf[0:4]) == 0x1a45dfa3 {
				parser = NewWebMParser()
			} else if bytes.NewBuffer(buf[4:8]).String() == "ftyp" {
				parser = NewISOBMFFParser()
			}

			if parser == nil {
				log.Printf("Unknown file type.\n")
				break
			}
		}

		if !parser.Append(buf[0:bytesRead]) {
			log.Printf("Parse error\n")
		}
	}

	if parser != nil {
		parser.EndOfData()
	}
}
