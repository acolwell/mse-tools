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

	buf := [1024]byte{}

	parser := NewWebMParser()

	for done := false; !done; {
		bytesRead, err := in.Read(buf[:])
		if err == io.EOF || err == io.ErrClosedPipe {
			done = true
		} else {
			parser.Append(buf[0:bytesRead])
		}
	}
	parser.EndOfData()
}
