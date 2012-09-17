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

package ebml

import (
	"errors"
	"fmt"
	"os"
)

type BufferWriter struct {
	buffer []byte
	offset int64
	length int64
}

func (bw *BufferWriter) Write(p []byte) (n int, err error) {
	writeLength := int64(len(p))
	bytesLeft := int64(len(bw.buffer)) - bw.length
	if writeLength > bytesLeft {
		newSize := int64(len(bw.buffer))
		for int64(len(p)) > bytesLeft {
			newSize = 2 * newSize
			bytesLeft = newSize - bw.length
		}
		newBuffer := make([]byte, newSize)
		copy(newBuffer, bw.buffer)
		bw.buffer = newBuffer
	}
	copy(bw.buffer[bw.offset:], p)
	bw.offset += writeLength
	if bw.offset > bw.length {
		bw.length = bw.offset
	}

	return len(p), nil
}

func (bw *BufferWriter) Seek(offset int64, whence int) (ret int64, err error) {
	newOffset := offset
	switch whence {
	case os.SEEK_SET:
		newOffset = offset
	case os.SEEK_CUR:
		newOffset = bw.offset + offset
	case os.SEEK_END:
		newOffset = bw.length + offset
	}

	if newOffset < 0 || newOffset > bw.length {
		return bw.offset, errors.New("Invalid offset")
	}

	bw.offset = newOffset
	return bw.offset, nil
}

func (bw *BufferWriter) Bytes() []byte {
	return bw.buffer[:bw.length]
}

func (bw *BufferWriter) Reset() {
	bw.offset = 0
	bw.length = 0
}

func NewBufferWriter(size int) *BufferWriter {
	bufferSize := size
	if bufferSize <= 0 {
		panic(fmt.Sprintf("Invalid buffer size %d", bufferSize))
	}
	return &BufferWriter{buffer: make([]byte, bufferSize), offset: 0, length: 0}
}
