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

import (
	"github.com/acolwell/mse-tools/go/ebml"
)

func WriteHeader(writer *ebml.Writer) (n int, err error) {
	bw := ebml.NewBufferWriter(1)
	w := ebml.NewWriter(bw)
	w.Write(ebml.IdVersion, 1)
	w.Write(ebml.IdReadVersion, 1)
	w.Write(ebml.IdMaxIDLength, 4)
	w.Write(ebml.IdMaxSizeLength, 8)
	w.Write(ebml.IdDocType, "webm")
	w.Write(ebml.IdDocTypeVersion, 2)
	w.Write(ebml.IdDocTypeReadVersion, 2)
	return writer.Write(ebml.IdHeader, bw.Bytes())
}