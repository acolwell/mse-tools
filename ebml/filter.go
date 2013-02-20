// Copyright 2013 Google Inc. All Rights Reserved.
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

package ebml

import (
	"log"
)

type filterClient struct {
	ids    map[int]bool
	writer *Writer
}

func (c *filterClient) allowId(id int) bool {
	_, hasId := c.ids[id]
	return !hasId
}

func (c *filterClient) OnListStart(offset int64, id int) bool {
	if c.allowId(id) {
		c.writer.WriteListStart(id)
	}

	return true
}

func (c *filterClient) OnListEnd(offset int64, id int) bool {
	if c.allowId(id) {
		c.writer.WriteListEnd(id)
	}

	return true
}

func (c *filterClient) OnBinary(id int, value []byte) bool {
	if c.allowId(id) {
		c.writer.Write(id, value)
	}

	return true
}

func (c *filterClient) OnInt(id int, value int64) bool {
	if c.allowId(id) {
		c.writer.Write(id, value)
	}

	return true
}

func (c *filterClient) OnUint(id int, value uint64) bool {
	if c.allowId(id) {
		c.writer.Write(id, value)
	}

	return true
}

func (c *filterClient) OnFloat(id int, value float64) bool {
	if c.allowId(id) {
		c.writer.Write(id, value)
	}

	return true
}

func (c *filterClient) OnString(id int, value string) bool {
	if c.allowId(id) {
		c.writer.Write(id, value)
	}

	return true
}

func Filter(input []byte, ids []int, typeMap map[int]int, unknownSizeInfo map[int][]int) []byte {
	idMap := map[int]bool{}
	for i := range ids {
		idMap[ids[i]] = true
	}

	writer := NewBufferWriter(len(input))
	parser := NewParser(GetListIDs(typeMap), unknownSizeInfo,
		NewElementParser(&filterClient{ids: idMap, writer: NewWriter(writer)}, typeMap))
	if !parser.Append(input) {
		log.Printf("Filter failed to parse input.\n")
		return nil
	}
	parser.EndOfData()
	return writer.Bytes()
}
