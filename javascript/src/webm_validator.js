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

/**
 * @constructor
 * @implements msetools.ParserClient
 * @implements ByteStreamValidator
 */
function WebMValidator() {
  this.parser_ = new msetools.ElementListParser(this);
  this.parserError_ = false;
  this.errors_ = [];
}

/**
 * @override
 */
WebMValidator.prototype.init = function(typeInfo) {

};

/**
 * @override
 */
WebMValidator.prototype.parse = function(data) {
  if (this.parserError_) {
    return ['Previously encountered a parser error.'];
  }

  this.errors_ = [];
  if (this.parser_.append(data) == msetools.ParserStatus.ERROR) {
    this.parserError_ = true;
  }

  var errors = this.errors_;
  this.errors_ = [];
  return errors;
};

/**
 * @override
 */
WebMValidator.prototype.reset = function() {
};

/**
 * @override
 */
WebMValidator.prototype.endOfStream = function() {

};


/**
 * WebM ID to name map.
 * @private
 * @const
 * @type {Object.<string, string>}
 */
var ID_TO_NAME_MAP_ = {
  '-1': 'ReservedID',
  'EC': 'Void',
  '1A45DFA3': 'EBMLHeader',
  '4286': 'EBMLVersion',
  '42F7': 'EBMLReadVersion',
  '42F2': 'EBMLMaxIDLength',
  '42F3': 'EBMLMaxSizeLength',
  '4282': 'DocType',
  '4287': 'DocTypeVersion',
  '4285': 'DocTypeReadVersion',
  '18538067': 'Segment',
  '114D9B74': 'SeekHead',
  '4DBB': 'Seek',
  '53AB': 'SeekID',
  '53AC': 'SeekPosition',
  '1549A966': 'Info',
  '2AD7B1': 'TimecodeScale',
  '4489': 'Duration',
  '4461': 'DateUTC',
  '4D80': 'MuxingApp',
  '5741': 'WritingApp',
  '1F43B675': 'Cluster',
  'E7': 'Timecode',
  'AB': 'PrevSize',
  'A3': 'SimpleBlock',
  'A0': 'BlockGroup',
  'A1': 'Block',
  '9B': 'BlockDuration',
  'FB': 'ReferenceBlock',
  '8E': 'Slices',
  'E8': 'TimeSlice',
  'CC': 'LaceNumber',
  '1654AE6B': 'Tracks',
  'AE': 'TrackEntry',
  'D7': 'TrackNumber',
  '73C5': 'TrackUID',
  '83': 'TrackType',
  'B9': 'FlagEnabled',
  '88': 'FlagDefault',
  '55AA': 'FlagForced',
  '9C': 'FlagLacing',
  '23E383': 'DefaultDuration',
  '536E': 'Name',
  '22B59C': 'Language',
  '86': 'CodecID',
  '63A2': 'CodecPrivate',
  '258688': 'CodecName',
  'E0': 'Video',
  '9A': 'FlagInterlaced',
  '53B8': 'StereoMode',
  'B0': 'PixelWidth',
  'BA': 'PixelHeight',
  '54AA': 'PixelCropBottom',
  '54BB': 'PixelCropTop',
  '54CC': 'PixelCropLeft',
  '54DD': 'PixelCropRight',
  '54B0': 'DisplayWidth',
  '54BA': 'DisplayHeight',
  '54B2': 'DisplayUnit',
  '54B3': 'AspectRatioType',
  'E1': 'Audio',
  'B5': 'SamplingFrequency',
  '78B5': 'OutputSamplingFrequency',
  '9F': 'Channels',
  '6264': 'BitDepth',
  '1C53BB6B': 'Cues',
  'BB': 'CuePoint',
  'B3': 'CueTime',
  'B7': 'CueTrackPositions',
  'F7': 'CueTrack',
  'F1': 'CueClusterPosition',
  '5378': 'CueBlockNumber'
};

/**
 * ElementType enum.
 * @enum {number}
 * @const
 */
var ElementType = {
  UNKNOWN: -1,
  BINARY: 1,
  DATE: 2,
  FLOAT: 3,
  INT: 4,
  LIST: 5,
  STRING: 6,
  UINT: 7,
  UTF8: 8
};

/**
 * Maps WebM IDs to ElementTypes.
 * @private
 * @type {Object.<string, ElementType>}
 */
var ID_TO_TYPE_MAP_ = {
  'Void': ElementType.BINARY,
  'EBMLHeader': ElementType.LIST,
  'EBMLVersion': ElementType.UINT,
  'EBMLReadVersion': ElementType.UINT,
  'EBMLMaxIDLength': ElementType.UINT,
  'EBMLMaxSizeLength': ElementType.UINT,
  'DocType': ElementType.STRING,
  'DocTypeVersion': ElementType.UINT,
  'DocTypeReadVersion': ElementType.UINT,
  'Segment': ElementType.LIST,
  'SeekHead': ElementType.LIST,
  'Seek': ElementType.LIST,
  'SeekID': ElementType.UINT,
  'SeekPosition': ElementType.UINT,
  'Info': ElementType.LIST,
  'TimecodeScale': ElementType.UINT,
  'Duration': ElementType.FLOAT,
  'DateUTC': ElementType.DATE,
  'MuxingApp': ElementType.UTF8,
  'WritingApp': ElementType.UTF8,
  'Cluster': ElementType.LIST,
  'Timecode': ElementType.UINT,
  'PrevSize': ElementType.UINT,
  'SimpleBlock': ElementType.BINARY,
  'BlockGroup': ElementType.LIST,
  'Block': ElementType.BINARY,
  'BlockDuration': ElementType.UINT,
  'ReferenceBlock': ElementType.INT,
  'Slices': ElementType.LIST,
  'TimeSlice': ElementType.LIST,
  'LaceNumber': ElementType.UINT,
  'Tracks': ElementType.LIST,
  'TrackEntry': ElementType.LIST,
  'TrackNumber': ElementType.UINT,
  'TrackUID': ElementType.UINT,
  'TrackType': ElementType.UINT,
  'FlagEnabled': ElementType.UINT,
  'FlagDefault': ElementType.UINT,
  'FlagForced': ElementType.UINT,
  'FlagLacing': ElementType.UINT,
  'DefaultDuration': ElementType.UINT,
  'Name': ElementType.UTF8,
  'Language': ElementType.STRING,
  'CodecID': ElementType.STRING,
  'CodecPrivate': ElementType.BINARY,
  'CodecName': ElementType.UTF8,
  'Video': ElementType.LIST,
  'FlagInterlaced': ElementType.UINT,
  'StereoMode': ElementType.UINT,
  'PixelWidth': ElementType.UINT,
  'PixelHeight': ElementType.UINT,
  'PixelCropBottom': ElementType.UINT,
  'PixelCropTop': ElementType.UINT,
  'PixelCropLeft': ElementType.UINT,
  'PixelCropRight': ElementType.UINT,
  'DisplayWidth': ElementType.UINT,
  'DisplayHeight': ElementType.UINT,
  'DisplayUnit': ElementType.UINT,
  'AspectRatioType': ElementType.UINT,
  'Audio': ElementType.LIST,
  'SamplingFrequency': ElementType.FLOAT,
  'OutputSamplingFrequency': ElementType.FLOAT,
  'Channels': ElementType.UINT,
  'BitDepth': ElementType.UINT,
  'Cues': ElementType.LIST,
  'CuePoint': ElementType.LIST,
  'CueTime': ElementType.UINT,
  'CueTrackPositions': ElementType.LIST,
  'CueTrack': ElementType.UINT,
  'CueClusterPosition': ElementType.UINT,
  'CueBlockNumber': ElementType.UINT
};


/**
 * Get the canonical key string for a WebM ID.
 * @param {number} id The WebM ID to convert.
 * @return {string} A capitalized hex string.
 * @private
 */
function getKeyForId_(id) {
  return id.toString(16).toUpperCase();
};


/**
 * Gets the element name for the WebM ID.
 * @param {number} id A WebM ID.
 * @return {string} The name for the specified ID.
 */
function getNameForId_(id) {
  var idStr = getKeyForId_(id);
  return ID_TO_NAME_MAP_[idStr] || ('UNKNOWN_ID(' + idStr + ')');
};


/**
 * @override
 */
WebMValidator.prototype.isIdAList = function(id) {
  var type = ID_TO_TYPE_MAP_[id] || ElementType.UNKNOWN;
  return (type == ElementType.LIST);
};

/**
 * Parses an element header ID field.
 * @param {?Uint8Array} buf The buffer to parse.
 * @return {{status: msetools.ParserStatus}|{status: msetools.ParserStatus,
 * bytesUsed: number, id: string}}
 * @private
 */
function parseWebMId_(buf) {
  if (!buf)
    return { status: msetools.ParserStatus.ERROR };
  
  if (buf.length < 1)
    return { status: msetools.ParserStatus.NEED_MORE_DATA };
  
  var bytesNeeded = 0;
  var mask = 0x80;
  var allOnesMask = 0x7f;
  for (var i = 1; i <= 4; i++) {
    if ((buf[0] & mask) == mask) {
      bytesNeeded = i;
      break;
    }
    mask >>= 1;
    allOnesMask >>= 1;
  }
  
  if (bytesNeeded == 0)
    return { status: msetools.ParserStatus.ERROR };
  
  if (buf.length < bytesNeeded)
    return { status: msetools.ParserStatus.NEED_MORE_DATA };


  /** @type {number} */ var raw_id = buf[0];
  var allOnes = (raw_id & allOnesMask) == allOnesMask;
  for (var i = 1; i < bytesNeeded; i++) {
    /** @type {number}*/ var ch = buf[i];
    raw_id = (raw_id * 256) + ch;
    allOnes = allOnes && (ch == 0xff);
  }
  
  var id = msetools.RESERVED_ID;
  if (!allOnes)
    id = getNameForId_(raw_id);
  
  return {
    status: msetools.ParserStatus.OK,
    bytesUsed: bytesNeeded,
    id: id
  };
};


/**
 * Parses an element header size field.
 * @param {?Uint8Array} buf The buffer to parse.
 * @return {{status: msetools.ParserStatus}|{status: msetools.ParserStatus,
 * bytesUsed: number, size: number}}
 * @private
 */
function parseWebMSize_(buf) {
  if (!buf)
    return { status: msetools.ParserStatus.ERROR };

  if (buf.length < 1)
    return { status: msetools.ParserStatus.NEED_MORE_DATA };

  var bytesNeeded = 0;
  var mask = 0x80;
  var sizeMask = 0x7f;
  for (var i = 1; i <= 8; i++) {
    if ((buf[0] & mask) == mask) {
      bytesNeeded = i;
      break;
    }
    mask >>= 1;
    sizeMask >>= 1;
  }

  if (bytesNeeded == 0)
    return { status: msetools.ParserStatus.ERROR, bytesUsed: 0, size: 0 };

  if (buf.length < bytesNeeded)
    return {
      status: msetools.ParserStatus.NEED_MORE_DATA, bytesUsed: 0, size: 0 };

  /** @type {number} */ var size = buf[0] & sizeMask;
  var allOnes = (size == sizeMask);
  for (var i = 1; i < bytesNeeded; i++) {
    /** @type {number}*/ var ch = buf[i];
    size = (size * 256) + ch;
    allOnes = allOnes && (ch == 0xff);
  }

  if (allOnes)
    size = msetools.RESERVED_SIZE;

  return {
    status: msetools.ParserStatus.OK,
    bytesUsed: bytesNeeded,
    size: size
  };
};

/**
 * @override
 */
WebMValidator.prototype.parseElementHeader = function(buf) {
  var currentBuffer = buf;
  var result = parseWebMId_(currentBuffer);

  if (result.status != msetools.ParserStatus.OK)
    return {status: result.status, bytesUsed: 0, id: '', size: 0 };
  var id = result.id;
  var bytesUsed = result.bytesUsed;

  currentBuffer = currentBuffer.subarray(result.bytesUsed);
  result = parseWebMSize_(currentBuffer);

  if (result.status != msetools.ParserStatus.OK)
    return {status: result.status, bytesUsed: 0, id: '', size: 0 };

  bytesUsed += result.bytesUsed;

  //window.console.log('msetools.parseElementHeader : id ' + msetools.getNameForId(id) +
  //     ' size ' + result.size);
  return {
    status: msetools.ParserStatus.OK,
    bytesUsed: bytesUsed,
    id: id,
    size: result.size
  };
};

WebMValidator.prototype.onListStart = function(id, elementPosition, 
                                               bodyPosition) {
  window.console.log('onListStart(' + id + 
                     ', ' + elementPosition +
                     ', ' + bodyPosition + ')');
  return msetools.ParserStatus.OK;
};


/**
 * Called when the end of a list element is parsed.
 * @param {string} id The ID for the list element.
 * @param {number} size The size of the list.
 * @return {boolean} True if the element was accepted by the client.
 * False if the client wants the parser to signal a parse error.
 */
WebMValidator.prototype.onListEnd = function(id, size) {
  window.console.log('onListEnd(' + id + 
                     ', ' + size + ')');
  return true;
};


/**
 * Called when a binary element is parsed.
 * @param {string} id The ID for the element.
 * @param {Uint8Array} value The value in the element.
 * @return {boolean} True if the element was accepted by the client.
 * False if the client wants the parser to signal a parse error.
 */
WebMValidator.prototype.onBinary = function(id, value) {
  //window.console.log('onBinary(' + id + ', ' + value.length + ')');
  return true;
};

msetools.WebMValidator = WebMValidator;
