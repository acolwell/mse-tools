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

const (
	IdReserved           = 0x1FFFFFFF
	IdVoid               = 0xEC
	IdCRC32              = 0xBF
	IdHeader             = 0x1A45DFA3
	IdVersion            = 0x4286
	IdReadVersion        = 0x42F7
	IdMaxIDLength        = 0x42F2
	IdMaxSizeLength      = 0x42F3
	IdDocType            = 0x4282
	IdDocTypeVersion     = 0x4287
	IdDocTypeReadVersion = 0x4285
)

var idTypes = map[int]int{
	IdVoid:               TypeBinary,
	IdCRC32:              TypeUint,
	IdHeader:             TypeList,
	IdVersion:            TypeUint,
	IdReadVersion:        TypeUint,
	IdMaxIDLength:        TypeUint,
	IdMaxSizeLength:      TypeUint,
	IdDocType:            TypeString,
	IdDocTypeVersion:     TypeUint,
	IdDocTypeReadVersion: TypeUint,
}

func IdTypes() map[int]int {
	return idTypes
}

var idToName = map[int]string{
	IdReserved:           "Reserved",
	IdVoid:               "Void",
	IdCRC32:              "CRC32",
	IdHeader:             "EBMLHeader",
	IdVersion:            "EBMLVersion",
	IdReadVersion:        "EBMLReadVersion",
	IdMaxIDLength:        "EBMLMaxIDLength",
	IdMaxSizeLength:      "EBMLMaxSizeLength",
	IdDocType:            "DocType",
	IdDocTypeVersion:     "DocTypeVersion",
	IdDocTypeReadVersion: "DocTypeReadVersion",
}

func IdToName(id int) string {
	return idToName[id]
}
