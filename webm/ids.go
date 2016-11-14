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

import "github.com/acolwell/mse-tools/ebml"

const (
	IdSegment                    = 0x18538067
	IdSeekHead                   = 0x114D9B74
	IdSeek                       = 0x4DBB
	IdSeekID                     = 0x53AB
	IdSeekPosition               = 0x53AC
	IdInfo                       = 0x1549A966
	IdSegmentUID                 = 0x73A4
	IdSegmentFilename            = 0x7384
	IdPrevUID                    = 0x3CB923
	IdPrevFilename               = 0x3C83AB
	IdNextUID                    = 0x3EB923
	IdNextFilename               = 0x3E83BB
	IdSegmentFamily              = 0x4444
	IdChapterTranslate           = 0x6924
	IdChapterTranslateEditionUID = 0x69FC
	IdChapterTranslateCodec      = 0x69BF
	IdChapterTranslateID         = 0x69A5
	IdTimecodeScale              = 0x2AD7B1
	IdDuration                   = 0x4489
	IdDateUTC                    = 0x4461
	IdTitle                      = 0x7BA9
	IdMuxingApp                  = 0x4D80
	IdWritingApp                 = 0x5741
	IdCluster                    = 0x1F43B675
	IdTimecode                   = 0xE7
	IdSilentTracks               = 0x5854
	IdSilentTrackNumber          = 0x58D7
	IdPosition                   = 0xA7
	IdPrevSize                   = 0xAB
	IdSimpleBlock                = 0xA3
	IdBlockGroup                 = 0xA0
	IdBlock                      = 0xA1
	IdBlockAdditions             = 0x75A1
	IdBlockMore                  = 0xA6
	IdBlockAddID                 = 0xEE
	IdBlockAdditional            = 0xA5
	IdBlockDuration              = 0x9B
	IdReferencePriority          = 0xFA
	IdReferenceBlock             = 0xFB
	IdCodecState                 = 0xA4
	IdDiscardPadding             = 0x75A2
	IdSlices                     = 0x8E
	IdTimeSlice                  = 0xE8
	IdLaceNumber                 = 0xCC
	IdTracks                     = 0x1654AE6B
	IdTrackEntry                 = 0xAE
	IdTrackNumber                = 0xD7
	IdTrackUID                   = 0x73C5
	IdTrackType                  = 0x83
	IdFlagEnabled                = 0xB9
	IdFlagDefault                = 0x88
	IdFlagForced                 = 0x55AA
	IdFlagLacing                 = 0x9C
	IdMinCache                   = 0x6DE7
	IdMaxCache                   = 0x6DF8
	IdDefaultDuration            = 0x23E383
	IdTrackTimecodeScale         = 0x23314F
	IdMaxBlockAdditionId         = 0x55EE
	IdName                       = 0x536E
	IdLanguage                   = 0x22B59C
	IdCodecID                    = 0x86
	IdCodecPrivate               = 0x63A2
	IdCodecName                  = 0x258688
	IdAttachmentLink             = 0x7446
	IdCodecDecodeAll             = 0xAA
	IdTrackOverlay               = 0x6FAB
	IdCodecDelay                 = 0x56AA
	IdSeekPreRoll                = 0x56BB
	IdTrackTranslate             = 0x6624
	IdTrackTranslateEditionUID   = 0x66FC
	IdTrackTranslateCodec        = 0x66BF
	IdTrackTranslateTrackID      = 0x66A5
	IdVideo                      = 0xE0
	IdFlagInterlaced             = 0x9A
	IdStereoMode                 = 0x53B8
	IdAlphaMode                  = 0x53C0
	IdPixelWidth                 = 0xB0
	IdPixelHeight                = 0xBA
	IdPixelCropBottom            = 0x54AA
	IdPixelCropTop               = 0x54BB
	IdPixelCropLeft              = 0x54CC
	IdPixelCropRight             = 0x54DD
	IdDisplayWidth               = 0x54B0
	IdDisplayHeight              = 0x54BA
	IdDisplayUnit                = 0x54B2
	IdAspectRatioType            = 0x54B3
	IdColorSpace                 = 0x2EB524
	IdFrameRate                  = 0x2383E3
	IdAudio                      = 0xE1
	IdSamplingFrequency          = 0xB5
	IdOutputSamplingFrequency    = 0x78B5
	IdChannels                   = 0x9F
	IdBitDepth                   = 0x6264
	IdTrackOperation             = 0xE2
	IdTrackCombinePlanes         = 0xE3
	IdTrackPlane                 = 0xE4
	IdTrackPlaneUID              = 0xE5
	IdTrackPlaneType             = 0xE6
	IdJoinBlocks                 = 0xE9
	IdTrackJoinUID               = 0xED
	IdContentEncodings           = 0x6D80
	IdContentEncoding            = 0x6240
	IdContentEncodingOrder       = 0x5031
	IdContentEncodingScope       = 0x5032
	IdContentEncodingType        = 0x5033
	IdContentCompression         = 0x5034
	IdContentCompAlgo            = 0x4254
	IdContentCompSettings        = 0x4255
	IdContentEncryption          = 0x5035
	IdContentEncAlgo             = 0x47E1
	IdContentEncKeyID            = 0x47E2
	IdContentSignature           = 0x47E3
	IdContentSigKeyID            = 0x47E4
	IdContentSigAlgo             = 0x47E5
	IdContentSigHashAlgo         = 0x47E6
	IdCues                       = 0x1C53BB6B
	IdCuePoint                   = 0xBB
	IdCueTime                    = 0xB3
	IdCueTrackPositions          = 0xB7
	IdCueTrack                   = 0xF7
	IdCueClusterPosition         = 0xF1
	IdCueRelativePosition        = 0xF0
	IdCueBlockNumber             = 0x5378
	IdCueCodecState              = 0xEA
	IdCueReference               = 0xDB
	IdCueRefTime                 = 0x96
	IdAttachments                = 0x1941A469
	IdAttachedFile               = 0x61A7
	IdFileDescription            = 0x467E
	IdFileName                   = 0x466E
	IdFileMimeType               = 0x4660
	IdFileData                   = 0x465C
	IdFileUID                    = 0x46AE
	IdChapters                   = 0x1043A770
	IdEditionEntry               = 0x45B9
	IdEditionUID                 = 0x45BC
	IdEditionFlagHidden          = 0x45BD
	IdEditionFlagDefault         = 0x45DB
	IdEditionFlagOrdered         = 0x45DD
	IdChapterAtom                = 0xB6
	IdChapterUID                 = 0x73C4
	IdChapterTimeStart           = 0x91
	IdChapterTimeEnd             = 0x92
	IdChapterFlagHidden          = 0x98
	IdChapterFlagEnabled         = 0x4598
	IdChapterSegmentUID          = 0x6E67
	IdChapterSegmentEditionUID   = 0x6EBC
	IdChapterPhysicalEquiv       = 0x63C3
	IdChapterTrack               = 0x8F
	IdChapterTrackNumber         = 0x89
	IdChapterDisplay             = 0x80
	IdChapString                 = 0x85
	IdChapLanguage               = 0x437C
	IdChapCountry                = 0x437E
	IdChapProcess                = 0x6944
	IdChapProcessCodecID         = 0x6955
	IdChapProcessPrivate         = 0x450D
	IdChapProcessCommand         = 0x6911
	IdChapProcessTime            = 0x6922
	IdChapProcessData            = 0x6933
	IdTags                       = 0x1254C367
	IdTag                        = 0x7373
	IdTargets                    = 0x63C0
	IdTargetTypeValue            = 0x68CA
	IdTargetType                 = 0x63CA
	IdTagTrackUID                = 0x63C5
	IdTagEditionUID              = 0x63C9
	IdTagChapterUID              = 0x63C4
	IdTagAttachmentUID           = 0x63C6
	IdSimpleTag                  = 0x67C8
	IdTagName                    = 0x45A3
	IdTagLanguage                = 0x447A
	IdTagDefault                 = 0x4484
	IdTagString                  = 0x4487
	IdTagBinary                  = 0x4485
)

var webmIdTypes = map[int]int{
	IdSegment:                 ebml.TypeList,
	IdSeekHead:                ebml.TypeList,
	IdSeek:                    ebml.TypeList,
	IdSeekID:                  ebml.TypeUint,
	IdSeekPosition:            ebml.TypeUint,
	IdInfo:                    ebml.TypeList,
	IdTimecodeScale:           ebml.TypeUint,
	IdDuration:                ebml.TypeFloat,
	IdDateUTC:                 ebml.TypeUint,
	IdMuxingApp:               ebml.TypeUTF8,
	IdWritingApp:              ebml.TypeUTF8,
	IdCluster:                 ebml.TypeList,
	IdTimecode:                ebml.TypeUint,
	IdPrevSize:                ebml.TypeUint,
	IdSimpleBlock:             ebml.TypeBinary,
	IdBlockGroup:              ebml.TypeList,
	IdBlock:                   ebml.TypeBinary,
	IdBlockDuration:           ebml.TypeUint,
	IdReferenceBlock:          ebml.TypeInt,
	IdCodecState:              ebml.TypeBinary,
	IdDiscardPadding:          ebml.TypeInt,
	IdSlices:                  ebml.TypeList,
	IdTimeSlice:               ebml.TypeList,
	IdLaceNumber:              ebml.TypeUint,
	IdTracks:                  ebml.TypeList,
	IdTrackEntry:              ebml.TypeList,
	IdTrackNumber:             ebml.TypeUint,
	IdTrackUID:                ebml.TypeUint,
	IdTrackType:               ebml.TypeUint,
	IdFlagEnabled:             ebml.TypeUint,
	IdFlagDefault:             ebml.TypeUint,
	IdFlagForced:              ebml.TypeUint,
	IdFlagLacing:              ebml.TypeUint,
	IdDefaultDuration:         ebml.TypeUint,
	IdName:                    ebml.TypeUTF8,
	IdLanguage:                ebml.TypeString,
	IdCodecID:                 ebml.TypeString,
	IdCodecPrivate:            ebml.TypeBinary,
	IdCodecName:               ebml.TypeUTF8,
	IdCodecDelay:              ebml.TypeUint,
	IdSeekPreRoll:             ebml.TypeUint,
	IdVideo:                   ebml.TypeList,
	IdFlagInterlaced:          ebml.TypeUint,
	IdStereoMode:              ebml.TypeUint,
	IdAlphaMode:               ebml.TypeUint,
	IdPixelWidth:              ebml.TypeUint,
	IdPixelHeight:             ebml.TypeUint,
	IdPixelCropBottom:         ebml.TypeUint,
	IdPixelCropTop:            ebml.TypeUint,
	IdPixelCropLeft:           ebml.TypeUint,
	IdPixelCropRight:          ebml.TypeUint,
	IdDisplayWidth:            ebml.TypeUint,
	IdDisplayHeight:           ebml.TypeUint,
	IdDisplayUnit:             ebml.TypeUint,
	IdAspectRatioType:         ebml.TypeUint,
	IdFrameRate:               ebml.TypeFloat,
	IdAudio:                   ebml.TypeList,
	IdSamplingFrequency:       ebml.TypeFloat,
	IdOutputSamplingFrequency: ebml.TypeFloat,
	IdChannels:                ebml.TypeUint,
	IdBitDepth:                ebml.TypeUint,
	IdCues:                    ebml.TypeList,
	IdCuePoint:                ebml.TypeList,
	IdCueTime:                 ebml.TypeUint,
	IdCueTrackPositions:       ebml.TypeList,
	IdCueTrack:                ebml.TypeUint,
	IdCueClusterPosition:      ebml.TypeUint,
	IdCueRelativePosition:     ebml.TypeUint,
	IdCueBlockNumber:          ebml.TypeUint,
}
var idTypes map[int]int = nil

func IdTypes() map[int]int {
	if idTypes == nil {
		idTypes = ebml.IdTypes()
		for k, v := range webmIdTypes {
			idTypes[k] = v
		}
	}
	return idTypes
}

var idToName = map[int]string{
	IdSegment:                    "Segment",
	IdSeekHead:                   "SeekHead",
	IdSeek:                       "Seek",
	IdSeekID:                     "SeekID",
	IdSeekPosition:               "SeekPosition",
	IdInfo:                       "Info",
	IdSegmentUID:                 "SegmentUID",
	IdSegmentFilename:            "SegmentFilename",
	IdPrevUID:                    "PrevUID",
	IdPrevFilename:               "PrevFilename",
	IdNextUID:                    "NextUID",
	IdNextFilename:               "NextFilename",
	IdSegmentFamily:              "SegmentFamily",
	IdChapterTranslate:           "ChapterTranslate",
	IdChapterTranslateEditionUID: "ChapterTranslateEditionUID",
	IdChapterTranslateCodec:      "ChapterTranslateCodec",
	IdChapterTranslateID:         "ChapterTranslateID",
	IdTimecodeScale:              "TimecodeScale",
	IdDuration:                   "Duration",
	IdDateUTC:                    "DateUTC",
	IdTitle:                      "Title",
	IdMuxingApp:                  "MuxingApp",
	IdWritingApp:                 "WritingApp",
	IdCluster:                    "Cluster",
	IdTimecode:                   "Timecode",
	IdSilentTracks:               "SilentTracks",
	IdSilentTrackNumber:          "SilentTrackNumber",
	IdPosition:                   "Position",
	IdPrevSize:                   "PrevSize",
	IdSimpleBlock:                "SimpleBlock",
	IdBlockGroup:                 "BlockGroup",
	IdBlock:                      "Block",
	IdBlockAdditions:             "BlockAdditions",
	IdBlockMore:                  "BlockMore",
	IdBlockAddID:                 "BlockAddID",
	IdBlockAdditional:            "BlockAdditional",
	IdBlockDuration:              "BlockDuration",
	IdReferencePriority:          "ReferencePriority",
	IdReferenceBlock:             "ReferenceBlock",
	IdCodecState:                 "CodecState",
	IdDiscardPadding:             "DiscardPadding",
	IdSlices:                     "Slices",
	IdTimeSlice:                  "TimeSlice",
	IdLaceNumber:                 "LaceNumber",
	IdTracks:                     "Tracks",
	IdTrackEntry:                 "TrackEntry",
	IdTrackNumber:                "TrackNumber",
	IdTrackUID:                   "TrackUID",
	IdTrackType:                  "TrackType",
	IdFlagEnabled:                "FlagEnabled",
	IdFlagDefault:                "FlagDefault",
	IdFlagForced:                 "FlagForced",
	IdFlagLacing:                 "FlagLacing",
	IdMinCache:                   "MinCache",
	IdMaxCache:                   "MaxCache",
	IdDefaultDuration:            "DefaultDuration",
	IdTrackTimecodeScale:         "TrackTimecodeScale",
	IdMaxBlockAdditionId:         "MaxBlockAdditionId",
	IdName:                       "Name",
	IdLanguage:                   "Language",
	IdCodecID:                    "CodecID",
	IdCodecPrivate:               "CodecPrivate",
	IdCodecName:                  "CodecName",
	IdAttachmentLink:             "AttachmentLink",
	IdCodecDecodeAll:             "CodecDecodeAll",
	IdTrackOverlay:               "TrackOverlay",
	IdCodecDelay:                 "CodecDelay",
	IdSeekPreRoll:                "SeekPreRoll",
	IdTrackTranslate:             "TrackTranslate",
	IdTrackTranslateEditionUID:   "TrackTranslateEditionUID",
	IdTrackTranslateCodec:        "TrackTranslateCodec",
	IdTrackTranslateTrackID:      "TrackTranslateTrackID",
	IdVideo:                      "Video",
	IdFlagInterlaced:             "FlagInterlaced",
	IdStereoMode:                 "StereoMode",
	IdAlphaMode:                  "AlphaMode",
	IdPixelWidth:                 "PixelWidth",
	IdPixelHeight:                "PixelHeight",
	IdPixelCropBottom:            "PixelCropBottom",
	IdPixelCropTop:               "PixelCropTop",
	IdPixelCropLeft:              "PixelCropLeft",
	IdPixelCropRight:             "PixelCropRight",
	IdDisplayWidth:               "DisplayWidth",
	IdDisplayHeight:              "DisplayHeight",
	IdDisplayUnit:                "DisplayUnit",
	IdAspectRatioType:            "AspectRatioType",
	IdColorSpace:                 "ColorSpace",
	IdFrameRate:                  "FrameRate",
	IdAudio:                      "Audio",
	IdSamplingFrequency:          "SamplingFrequency",
	IdOutputSamplingFrequency:    "OutputSamplingFrequency",
	IdChannels:                   "Channels",
	IdBitDepth:                   "BitDepth",
	IdTrackOperation:             "TrackOperation",
	IdTrackCombinePlanes:         "TrackCombinePlanes",
	IdTrackPlane:                 "TrackPlane",
	IdTrackPlaneUID:              "TrackPlaneUID",
	IdTrackPlaneType:             "TrackPlaneType",
	IdJoinBlocks:                 "JoinBlocks",
	IdTrackJoinUID:               "TrackJoinUID",
	IdContentEncodings:           "ContentEncodings",
	IdContentEncoding:            "ContentEncoding",
	IdContentEncodingOrder:       "ContentEncodingOrder",
	IdContentEncodingScope:       "ContentEncodingScope",
	IdContentEncodingType:        "ContentEncodingType",
	IdContentCompression:         "ContentCompression",
	IdContentCompAlgo:            "ContentCompAlgo",
	IdContentCompSettings:        "ContentCompSettings",
	IdContentEncryption:          "ContentEncryption",
	IdContentEncAlgo:             "ContentEncAlgo",
	IdContentEncKeyID:            "ContentEncKeyID",
	IdContentSignature:           "ContentSignature",
	IdContentSigKeyID:            "ContentSigKeyID",
	IdContentSigAlgo:             "ContentSigAlgo",
	IdContentSigHashAlgo:         "ContentSigHashAlgo",
	IdCues:                       "Cues",
	IdCuePoint:                   "CuePoint",
	IdCueTime:                    "CueTime",
	IdCueTrackPositions:          "CueTrackPositions",
	IdCueTrack:                   "CueTrack",
	IdCueClusterPosition:         "CueClusterPosition",
	IdCueRelativePosition:        "CueRelativePosition",
	IdCueBlockNumber:             "CueBlockNumber",
	IdCueCodecState:              "CueCodecState",
	IdCueReference:               "CueReference",
	IdCueRefTime:                 "CueRefTime",
	IdAttachments:                "Attachments",
	IdAttachedFile:               "AttachedFile",
	IdFileDescription:            "FileDescription",
	IdFileName:                   "FileName",
	IdFileMimeType:               "FileMimeType",
	IdFileData:                   "FileData",
	IdFileUID:                    "FileUID",
	IdChapters:                   "Chapters",
	IdEditionEntry:               "EditionEntry",
	IdEditionUID:                 "EditionUID",
	IdEditionFlagHidden:          "EditionFlagHidden",
	IdEditionFlagDefault:         "EditionFlagDefault",
	IdEditionFlagOrdered:         "EditionFlagOrdered",
	IdChapterAtom:                "ChapterAtom",
	IdChapterUID:                 "ChapterUID",
	IdChapterTimeStart:           "ChapterTimeStart",
	IdChapterTimeEnd:             "ChapterTimeEnd",
	IdChapterFlagHidden:          "ChapterFlagHidden",
	IdChapterFlagEnabled:         "ChapterFlagEnabled",
	IdChapterSegmentUID:          "ChapterSegmentUID",
	IdChapterSegmentEditionUID:   "ChapterSegmentEditionUID",
	IdChapterPhysicalEquiv:       "ChapterPhysicalEquiv",
	IdChapterTrack:               "ChapterTrack",
	IdChapterTrackNumber:         "ChapterTrackNumber",
	IdChapterDisplay:             "ChapterDisplay",
	IdChapString:                 "ChapString",
	IdChapLanguage:               "ChapLanguage",
	IdChapCountry:                "ChapCountry",
	IdChapProcess:                "ChapProcess",
	IdChapProcessCodecID:         "ChapProcessCodecID",
	IdChapProcessPrivate:         "ChapProcessPrivate",
	IdChapProcessCommand:         "ChapProcessCommand",
	IdChapProcessTime:            "ChapProcessTime",
	IdChapProcessData:            "ChapProcessData",
	IdTags:                       "Tags",
	IdTag:                        "Tag",
	IdTargets:                    "Targets",
	IdTargetTypeValue:            "TargetTypeValue",
	IdTargetType:                 "TargetType",
	IdTagTrackUID:                "TagTrackUID",
	IdTagEditionUID:              "TagEditionUID",
	IdTagChapterUID:              "TagChapterUID",
	IdTagAttachmentUID:           "TagAttachmentUID",
	IdSimpleTag:                  "SimpleTag",
	IdTagName:                    "TagName",
	IdTagLanguage:                "TagLanguage",
	IdTagDefault:                 "TagDefault",
	IdTagString:                  "TagString",
	IdTagBinary:                  "TagBinary",
}

func IdToName(id int) string {
	if name, ok := idToName[id]; ok {
		return name
	}
	return ebml.IdToName(id)
}

func UnknownSizeInfo() map[int][]int {
	return map[int][]int{
		IdSegment: []int{
			ebml.IdHeader,
			IdSegment},
		IdCluster: []int{
			IdSegment,
			IdSeekHead,
			IdInfo,
			IdCluster,
			IdTracks,
			IdCues,
			IdAttachments,
			IdChapters,
			IdTags}}
}
