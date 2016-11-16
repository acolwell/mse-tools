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
	"flag"
	"fmt"
	"github.com/acolwell/mse-tools/ebml"
	"github.com/acolwell/mse-tools/webm"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	EOL            = "\r\n"
	CR             = 0x0d
	LF             = 0x0a
	END_OF_HEADERS = "\r\n\r\n"

	SEEK_HEAD_RESERVE_SIZE = 124
)

type Cue struct {
	timecode int64
	offset   int64
	trackID  uint64
}

type DemuxerClient struct {
	writer                 *ebml.Writer
	minClusterDurationInMS int
	readEBMLHeader         bool
	audioTrackID           uint64
	videoTrackID           uint64
	timecodeScale          uint64
	duration               float64
	segmentOffset          int64
	startTimecode          int64
	clusterTimecode        int64
	tracks                 []webm.Track
	isVorbis               map[uint64]bool
	blocks                 map[uint64][]*Block
	cues                   []Cue

	minClusterDuration    int64
	outputSegmentOffset   int64
	outputInfoOffset      int64
	outputTracksOffset    int64
	outputClusterOffset   int64
	outputCuesOffset      int64
	outputClusterTimecode int64
}

type Block struct {
	id                  uint64
	isSimple            bool
	timecode            int64
	flags               uint8
	data                []byte
	extraBlockGroupData []byte
}

func (b *Block) IsKeyframe() bool {
	return (b.flags & 0x80) != 0
}

func NewBlock(id uint64, isSimple bool, timecode int64, flags uint8, data []byte,
	extraBlockGroupData []byte) *Block {
	block := &Block{id: id, isSimple: isSimple, flags: flags, timecode: timecode, data: make([]byte, len(data)), extraBlockGroupData: make([]byte, len(extraBlockGroupData))}
	copy(block.data, data)
	copy(block.extraBlockGroupData, extraBlockGroupData)
	return block
}

type BlockGroupClient struct {
	parsedBlock bool
	id          uint64
	rawTimecode int64
	flags       uint8
	blockData   []byte
	writer      *ebml.Writer
}

func (c *BlockGroupClient) OnListStart(offset int64, id int) bool {
	//log.Printf("OnListStart(%d, %s)\n", offset, webm.IdToName(id))
	log.Printf("OnListStart() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *BlockGroupClient) OnListEnd(offset int64, id int) bool {
	//log.Printf("OnListEnd(%d, %s)\n", offset, webm.IdToName(id))
	log.Printf("OnListEnd() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *BlockGroupClient) OnBinary(id int, value []byte) bool {
	if id == webm.IdBlock {
		blockInfo := webm.ParseSimpleBlock(value)
		c.id = blockInfo.Id
		c.rawTimecode = int64(blockInfo.Timecode)
		c.flags = blockInfo.Flags & 0x0f
		c.blockData = value[blockInfo.HeaderSize:]
		c.parsedBlock = true
		return true
	} else if id == webm.IdBlockAdditions {
		c.writer.Write(id, value)
		return true
	}
	log.Printf("OnBinary() : Unexpected element %s size %d\n", webm.IdToName(id), len(value))
	return false
}

func (c *BlockGroupClient) OnInt(id int, value int64) bool {
	if id == webm.IdDiscardPadding || id == webm.IdReferenceBlock {
		c.writer.Write(id, value)
		return true
	}
	log.Printf("OnInt() : Unexpected element %s %d\n", webm.IdToName(id), value)
	return false
}

func (c *BlockGroupClient) OnUint(id int, value uint64) bool {
	if id == webm.IdBlockDuration {
		c.writer.Write(id, value)
		return true
	}
	log.Printf("OnUint() : Unexpected element %s %u\n", webm.IdToName(id), value)
	return false
}

func (c *BlockGroupClient) OnFloat(id int, value float64) bool {
	log.Printf("OnFloat() : Unexpected element %s %f\n", webm.IdToName(id), value)
	return false
}

func (c *BlockGroupClient) OnString(id int, value string) bool {
	log.Printf("OnString() : Unexpected element %s %s\n", webm.IdToName(id), value)
	return false
}

func (c *DemuxerClient) OnListStart(offset int64, id int) bool {
	//log.Printf("OnListStart(%d, %s)\n", offset, webm.IdToName(id))

	if !c.readEBMLHeader {
		log.Printf("Unexpected element %s before EBMLHeader\n", webm.IdToName(id))
		return false
	}

	if id == webm.IdSegment {
		c.segmentOffset = offset

		c.writer.WriteListStart(webm.IdSegment)
		c.outputSegmentOffset = c.writer.Offset()
		c.writer.WriteVoid(SEEK_HEAD_RESERVE_SIZE)
		return true
	}

	if id == webm.IdCluster {
		c.clusterTimecode = -1
		if c.outputClusterOffset == -1 {
			c.outputClusterOffset = c.writer.Offset()
		}
		return true
	}

	log.Printf("OnListStart() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *DemuxerClient) OnListEnd(offset int64, id int) bool {
	//log.Printf("OnListEnd(%d, %s)\n", offset, webm.IdToName(id))

	if id == webm.IdSegment {
		if c.outputClusterTimecode != -1 {
			c.writeRemainingBlocks()
			c.writer.WriteListEnd(webm.IdCluster)
		}

		if c.writer.CanSeek() {
			c.writeCues()
		}

		// Rewrite seek head.
		oldOffset := c.writer.Offset()
		if c.writer.SetOffset(c.outputSegmentOffset) {
			c.writeSeekHead()

			if c.writer.Offset() < c.outputInfoOffset {
				c.writer.WriteVoid(int(c.outputInfoOffset - c.writer.Offset()))
			}

			c.writer.SetOffset(oldOffset)
		}

		c.writer.WriteListEnd(webm.IdSegment)
		return true
	}

	if id == webm.IdCluster {
		return true
	}

	log.Printf("OnListEnd() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *DemuxerClient) OnBinary(id int, value []byte) bool {
	if id == ebml.IdHeader {
		if c.readEBMLHeader {
			log.Printf("Already read an EBMLHeader\n")
			return false
		}
		if !c.ParseEBMLHeader(value) {
			return false
		}
		c.readEBMLHeader = true
		webm.WriteHeader(c.writer)
		//c.writer.Write(id, value)
		return true
	}

	if !c.readEBMLHeader {
		log.Printf("Unexpected element %s before EBMLHeader\n", webm.IdToName(id))
		return false
	}

	if id == ebml.IdVoid {
		return true
	}

	if id == webm.IdSeekHead {
		return true
	}

	if id == webm.IdInfo {
		if !c.ParseInfo(value) {
			return false
		}
		c.outputInfoOffset = c.writer.Offset()
		c.writer.Write(id, value)
		return true

	}

	if id == webm.IdTracks {
		if !c.ParseTracks(value) {
			return false
		}
		c.outputTracksOffset = c.writer.Offset()

		// Filter out deprecated values.
		filteredValue := webm.Filter(value, []int{webm.IdFrameRate})

		c.writer.Write(id, filteredValue)
		return true
	}

	if id == webm.IdSimpleBlock {
		return c.ParseSimpleBlock(value)
	}

	if id == webm.IdBlockGroup {
		return c.ParseBlockGroup(value)
	}

	switch id {
	case webm.IdCues,
		webm.IdPrevSize,
		webm.IdPosition:
		return true
	}

	log.Printf("OnBinary() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *DemuxerClient) OnInt(id int, value int64) bool {
	log.Printf("OnInt() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *DemuxerClient) OnUint(id int, value uint64) bool {
	if !c.readEBMLHeader {
		log.Printf("Unexpected element %s before EBMLHeader\n", webm.IdToName(id))
		return false
	}

	if id == webm.IdTimecode {
		c.clusterTimecode = int64(value)
		//log.Printf("Input Cluster timecode %d\n", c.clusterTimecode)
		return true
	}

	log.Printf("OnUint() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *DemuxerClient) OnFloat(id int, value float64) bool {
	if !c.readEBMLHeader {
		log.Printf("Unexpected element %s before EBMLHeader\n", webm.IdToName(id))
		return false
	}

	log.Printf("OnFloat() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *DemuxerClient) OnString(id int, value string) bool {
	if !c.readEBMLHeader {
		log.Printf("Unexpected element %s before EBMLHeader\n", webm.IdToName(id))
		return false
	}

	log.Printf("OnString() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *DemuxerClient) writeSeekHead() {
	c.writer.WriteListStart(webm.IdSeekHead)
	if c.outputInfoOffset > c.outputSegmentOffset {
		c.writeSeek(webm.IdInfo, c.outputInfoOffset)
	}
	if c.outputTracksOffset > c.outputSegmentOffset {
		c.writeSeek(webm.IdTracks, c.outputTracksOffset)
	}
	if c.outputClusterOffset > c.outputSegmentOffset {
		c.writeSeek(webm.IdCluster, c.outputClusterOffset)
	}
	if c.outputCuesOffset > c.outputSegmentOffset {
		c.writeSeek(webm.IdCues, c.outputCuesOffset)
	}
	c.writer.WriteListEnd(webm.IdSeekHead)
}

func (c *DemuxerClient) writeSeek(id int, offset int64) {
	c.writer.WriteListStart(webm.IdSeek)
	c.writer.Write(webm.IdSeekID, uint32(id))
	c.writer.Write(webm.IdSeekPosition, uint64(offset-c.outputSegmentOffset))
	c.writer.WriteListEnd(webm.IdSeek)
}

func (c *DemuxerClient) ParseEBMLHeader(buf []byte) bool {
	header := ebml.ParseHeader(buf)
	if header == nil {
		log.Printf("Failed to parse EBML header\n")
		return false
	}

	if header.DocType() != "webm" {
		log.Printf("EBML header has an unsupported DocType '%s'\n", header.DocType())
		return false
	}

	if header.DocTypeReadVersion() != 2 {
		log.Printf("EBML header has an unsupported DocTypeReadVersion %d\n", header.DocTypeReadVersion())
		return false
	}

	return true
}

func (c *DemuxerClient) ParseInfo(buf []byte) bool {
	info := webm.ParseInfoElement(buf)
	if info == nil {
		log.Printf("Failed to parse Info element\n")
		return false
	}

	scale := float64(1000000000 / info.TimecodeScale())
	c.minClusterDuration = int64(scale * float64(c.minClusterDurationInMS) / 1000.0)

	return true
}

func (c *DemuxerClient) ParseTracks(buf []byte) bool {
	c.tracks = webm.ParseTracksElement(buf)
	for i := range c.tracks {
		id := c.tracks[i].ID()
		c.blocks[id] = []*Block{}
		c.isVorbis[id] = c.tracks[i].CodecID() == "A_VORBIS"
	}

	return c.tracks != nil
}

func (c *DemuxerClient) ParseSimpleBlock(buf []byte) bool {
	if c.clusterTimecode == -1 {
		panic("Got a simple block before the cluster timecode.")
	}

	if len(buf) < 3 {
		log.Printf("Invalid simple block size %d\n", len(buf))
		return false
	}

	mask := byte(0x80)
	idSize := 1
	for ; (buf[0]&mask) == 0 && idSize < 8; idSize++ {
		mask >>= 1
	}

	if len(buf) < idSize+3 {
		log.Printf("Invalid simple block size %d\n", len(buf))
		return false
	}

	id := uint64(buf[0] & (mask - 1))
	for i := 1; i < idSize; i++ {
		id = (id << 8) | uint64(buf[i])
	}

	rawTimecode := int64(buf[idSize])<<8 | int64(buf[idSize+1])
	if (rawTimecode & 0x8000) != 0 {
		rawTimecode |= (-1 << 16)
	}
	timecode := c.clusterTimecode + rawTimecode
	flags := buf[idSize+2]
	//log.Printf("in track %d %d 0x%x %d\n", id, timecode, flags, len(buf)-3-idSize)

	if c.startTimecode == -1 {
		c.startTimecode = timecode
	}

	// Fix any Vorbis blocks that don't have the keyframe flag set. This has been
	// observed in WebM files that specify Flix as the MuxingApp and WritingApp.
	if c.isVorbis[id] && (flags&0x80) != 0x80 {
		flags |= 0x80
	}

	blockList, ok := c.blocks[id]
	if !ok {
		return false
	}

	c.blocks[id] = append(blockList, NewBlock(id, true, timecode, flags, buf[idSize+3:], []byte{}))

	c.tryWritingNextBlock()
	return true
}

func (c *DemuxerClient) ParseBlockGroup(buf []byte) bool {
	if c.clusterTimecode == -1 {
		panic("Got a block group before the cluster timecode.")
	}

	typeInfo := map[int]int{
		webm.IdBlock:          ebml.TypeBinary,
		webm.IdBlockAdditions: ebml.TypeBinary,
		webm.IdBlockDuration:  ebml.TypeUint,
		webm.IdReferenceBlock: ebml.TypeInt,
		webm.IdDiscardPadding: ebml.TypeInt,
	}

	bw := ebml.NewBufferWriter(len(buf))
	bc := &BlockGroupClient{parsedBlock: false, writer: ebml.NewWriter(bw)}
	p := ebml.NewParser(ebml.GetListIDs(typeInfo), webm.UnknownSizeInfo(),
		ebml.NewElementParser(bc, typeInfo))

	if !p.Append(buf) {
		log.Printf("Parser error")
		return false
	}
	p.EndOfData()

	id := bc.id
	rawTimecode := bc.rawTimecode
	flags := bc.flags

	timecode := c.clusterTimecode + rawTimecode

	//log.Printf("in track %d %d 0x%x %d\n", id, timecode, flags, len(buf)-3-idSize)

	if c.startTimecode == -1 {
		c.startTimecode = timecode
	}

	blockList, ok := c.blocks[id]
	if !ok {
		return false
	}

	c.blocks[id] = append(blockList, NewBlock(id, false, timecode, flags, bc.blockData, bw.Bytes()))

	c.tryWritingNextBlock()
	return true
}

func (c *DemuxerClient) tryWritingNextBlock() {
	audioID := uint64(0)
	videoID := uint64(0)
	var audio []*Block = nil
	var video []*Block = nil

	for i := range c.tracks {
		if c.tracks[i].Type() == webm.VIDEO_TRACK {
			videoID = c.tracks[i].ID()
			video = c.blocks[videoID]
		} else if c.tracks[i].Type() == webm.AUDIO_TRACK {
			audioID = c.tracks[i].ID()
			audio = c.blocks[audioID]
		}
	}

	if video == nil {
		c.writeNextSingleStreamBlock(audioID)
		return
	}

	if audio == nil {
		c.writeNextSingleStreamBlock(videoID)
		return
	}

	if len(video) < 1 || len(audio) < 2 {
		return
	}
	videoBlock := video[0]
	audioBlock1 := audio[0]
	audioBlock2 := audio[1]

	if videoBlock.IsKeyframe() &&
		audioBlock1.IsKeyframe() &&
		audioBlock1.timecode <= videoBlock.timecode &&
		audioBlock2.timecode > videoBlock.timecode &&
		(audioBlock1.timecode-c.outputClusterTimecode) >= c.minClusterDuration {
		// This is the situation where a new cluster is allowed.
		c.startNewCluster(audioBlock1.id, audioBlock1.timecode)
	}

	if audioBlock1.timecode <= videoBlock.timecode {
		c.writeBlock(audioBlock1)
		c.blocks[audioID] = audio[1:]
		return
	}
	c.writeBlock(videoBlock)
	c.blocks[videoID] = video[1:]
}

func (c *DemuxerClient) writeNextSingleStreamBlock(trackID uint64) {
	blocks := c.blocks[trackID]

	if len(blocks) < 2 {
		return
	}

	block := blocks[0]
	clusterDuration := block.timecode - c.outputClusterTimecode
	if block.IsKeyframe() &&
		clusterDuration >= c.minClusterDuration {
		c.startNewCluster(block.id, block.timecode)
	}
	c.writeBlock(block)
	c.blocks[trackID] = blocks[1:]
}

func (c *DemuxerClient) startNewCluster(id uint64, timecode int64) {
	//log.Printf("Output Cluster timecode %d\n", timecode)

	if c.outputClusterTimecode != -1 {
		c.writer.WriteListEnd(webm.IdCluster)
	}

	c.cues = append(c.cues, Cue{timecode: timecode, offset: c.writer.Offset(), trackID: id})

	if timecode < 0 {
		panic(fmt.Sprintf("Negative cluster timecode (%d) not allowed!", timecode))
	}
	c.outputClusterTimecode = timecode
	c.writer.WriteListStart(webm.IdCluster)
	c.writer.Write(webm.IdTimecode, c.outputClusterTimecode)

}

func (c *DemuxerClient) writeBlock(block *Block) {
	//log.Printf("out track %d %d 0x%x %d\n", block.id, block.timecode, block.flags, len(block.data))

	if block.id > 0x7f {
		panic("Can't write track numbers yet.")
	}

	if c.outputClusterTimecode == -1 {
		if !block.IsKeyframe() {
			panic("First block is not a keyframe!")
		}
		c.startNewCluster(block.id, block.timecode)
	}

	rawTimecode := block.timecode - c.outputClusterTimecode

	if rawTimecode > 0x7fff {
		panic(fmt.Sprintf("rawTimecode is too big %d (%d)\n", rawTimecode, block.timecode))
	}

	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte(0x80 | byte(block.id))
	buffer.WriteByte(byte(rawTimecode >> 8))
	buffer.WriteByte(byte(rawTimecode & 0xff))
	buffer.WriteByte(block.flags)
	buffer.Write(block.data)
	if block.isSimple {
		c.writer.Write(webm.IdSimpleBlock, buffer.Bytes())
	} else {
		c.writer.WriteListStart(webm.IdBlockGroup)
		c.writer.Write(webm.IdBlock, buffer.Bytes())
		c.writer.WriteToOutput(block.extraBlockGroupData)
		// TODO
		c.writer.WriteListEnd(webm.IdBlockGroup)
	}
}

func (c *DemuxerClient) writeRemainingBlocks() {
	for {
		var minBlock *Block = nil
		for _, blockList := range c.blocks {
			if len(blockList) == 0 {
				continue
			}

			block := blockList[0]
			if minBlock == nil || block.timecode < minBlock.timecode {
				minBlock = block
			}
		}

		if minBlock == nil {
			break
		}

		c.writeBlock(minBlock)
		c.blocks[minBlock.id] = c.blocks[minBlock.id][1:]
	}
}

func (c *DemuxerClient) writeCues() {
	c.outputCuesOffset = c.writer.Offset()
	c.writer.WriteListStart(webm.IdCues)
	for i := range c.cues {
		cue := c.cues[i]
		c.writer.WriteListStart(webm.IdCuePoint)
		c.writer.Write(webm.IdCueTime, cue.timecode)
		c.writer.WriteListStart(webm.IdCueTrackPositions)
		c.writer.Write(webm.IdCueTrack, cue.trackID)
		c.writer.Write(webm.IdCueClusterPosition, cue.offset-c.outputSegmentOffset)
		c.writer.WriteListEnd(webm.IdCueTrackPositions)
		c.writer.WriteListEnd(webm.IdCuePoint)
	}
	c.writer.WriteListEnd(webm.IdCues)
}
func NewDemuxerClient(writer *ebml.Writer, minClusterDurationInMS int) *DemuxerClient {
	return &DemuxerClient{
		writer:                 writer,
		minClusterDurationInMS: minClusterDurationInMS,
		readEBMLHeader:         false,
		audioTrackID:           0,
		videoTrackID:           0,
		timecodeScale:          0,
		duration:               0.0,
		segmentOffset:          -1,
		startTimecode:          -1,
		clusterTimecode:        -1,
		tracks:                 nil,
		isVorbis:               map[uint64]bool{},
		blocks:                 map[uint64][]*Block{},
		cues:                   []Cue{},
		outputSegmentOffset:    -1,
		outputInfoOffset:       -1,
		outputTracksOffset:     -1,
		outputClusterOffset:    -1,
		outputCuesOffset:       -1,
		outputClusterTimecode:  -1,
	}
}

func checkError(str string, err error) {
	if err != nil {
		log.Printf("Error: %s - %s\n", str, err.Error())
		os.Exit(-1)
	}
}

func main() {
	var minClusterDurationInMS int
	flag.IntVar(&minClusterDurationInMS, "cm", 250, "Minimum Cluster Duration (ms)")
	flag.Parse()

	if minClusterDurationInMS < 0 || minClusterDurationInMS > 30000 {
		log.Printf("Invalid minimum cluster duration\n")
		os.Exit(-1)
	}

	if len(flag.Args()) < 2 {
		log.Printf("Usage: %s [-cm <duration>] <infile> <outfile>\n", os.Args[0])
		return
	}

	var in *os.File = nil
	var err error = nil

	inputArg := flag.Arg(0)
	outputArg := flag.Arg(1)

	if inputArg == "-" {
		in = os.Stdin
	} else {
		in, err = os.Open(inputArg)
		checkError("Open input", err)
	}

	var out *ebml.Writer = nil
	if outputArg == "-" {
		out = ebml.NewNonSeekableWriter(io.WriteSeeker(os.Stdout))
	} else {
		if inputArg == outputArg {
			log.Printf("Input and output filenames can't be the same.\n")
			return
		}

		if strings.HasPrefix(outputArg, "ws://") {
			url, err := url.Parse(outputArg)
			checkError("Output url", err)

			origin := "http://localhost/"
			ws, err := websocket.Dial(url.String(), "", origin)
			checkError("WebSocket Dial", err)
			out = ebml.NewNonSeekableWriter(io.Writer(ws))
		} else {
			file, err := os.Create(outputArg)
			if err != nil {
				log.Printf("Failed to create '%s'; err=%s\n", outputArg, err.Error())
				os.Exit(1)
			}
			out = ebml.NewWriter(io.WriteSeeker(file))
		}
	}

	buf := [1024]byte{}
	c := NewDemuxerClient(out, minClusterDurationInMS)

	typeInfo := map[int]int{
		ebml.IdHeader:      ebml.TypeBinary,
		webm.IdSegment:     ebml.TypeList,
		webm.IdInfo:        ebml.TypeBinary,
		webm.IdTracks:      ebml.TypeBinary,
		webm.IdCluster:     ebml.TypeList,
		webm.IdTimecode:    ebml.TypeUint,
		webm.IdSimpleBlock: ebml.TypeBinary,
	}

	parser := ebml.NewParser(ebml.GetListIDs(typeInfo), webm.UnknownSizeInfo(),
		ebml.NewElementParser(c, typeInfo))

	for done := false; !done; {
		bytesRead, err := in.Read(buf[:])
		if err == io.EOF || err == io.ErrClosedPipe {
			parser.EndOfData()
			done = true
			continue
		}

		if !parser.Append(buf[0:bytesRead]) {
			log.Printf("Parser error")
			done = true
			continue
		}
	}
}
