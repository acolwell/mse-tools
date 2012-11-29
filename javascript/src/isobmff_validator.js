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
function ISOBMFFValidator() {
  this.parser_ = new msetools.ElementListParser(this);
  this.parserError_ = false;
  this.errors_ = [];
}

/**
 * @type {msetools.ElementListParser}
 */
ISOBMFFValidator.prototype.parser_ = null;

/**
 * @type {number}
 */
ISOBMFFValidator.prototype.default_sample_duration_ = -1;

/**
 * @type {number}
 */
ISOBMFFValidator.prototype.default_sample_size_ = -1;

/**
 * @type {number}
 */
ISOBMFFValidator.prototype.default_sample_flags_ = 0;

/**
 * @override
 */
ISOBMFFValidator.prototype.init = function(typeInfo) {

};

/**
 * @override
 */
ISOBMFFValidator.prototype.parse = function(data) {
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
ISOBMFFValidator.prototype.reset = function() {

};

/**
 * @override
 */
ISOBMFFValidator.prototype.endOfStream = function() {

};

/**
 * @override
 */
ISOBMFFValidator.prototype.parseElementHeader = function(buf) {
  var ERROR_STATUS = {status: msetools.ParserStatus.ERROR,
                      bytesUsed: 0, id: '', size: 0 };

  if (buf.length < 8) {
    return {status: msetools.ParserStatus.NEED_MORE_DATA,
            bytesUsed: 0, id: '', size: 0 };
  }

  var size = this.getUint32_(buf);
  var bytesUsed = 4;

  if (size == 0) {
    this.errors_.push('Box size of 0 not allowed!.');
    return ERROR_STATUS;
  }

  if (size == 1) {
    this.errors_.push('64-bit box sizes not supported yet!.');
    return ERROR_STATUS;
  }


  var id = '';
  for (var i = 4; i < 8; ++i) {
    id += String.fromCharCode(buf[i]);
    bytesUsed++;
  }

  if (id == 'uuid') {
    this.errors_.push('uuid boxes not supported yet!.');
    return ERROR_STATUS;
  }

  if (size < bytesUsed) {
    this.errors_.push('Invalid box size ' + size);
    return ERROR_STATUS;
  }
  // Subtract off the header size.
  size -= bytesUsed;

  //console.log('id ' + id + ' size ' + size);
  return {
    status: msetools.ParserStatus.OK,
    bytesUsed: bytesUsed,
    id: id,
    size: size
  };
};

/**
 * Indicates which element IDs are list elements.
 *
 * @type {Object.<string, boolean>}
 */
var ID_IS_LIST_MAP = {
  'moov': true,
  'moof': true,
  'traf': true,
  'trak': true,
  'mdia': true,
  'minf': true,
  'stbl': true
};

/**
 * @override
 */
ISOBMFFValidator.prototype.isIdAList = function(id) {
  return ID_IS_LIST_MAP[id] || false;
};

/**
 * Indicates which element IDs are full box elements.
 *
 * @type {Object.<string, boolean>}
 */
var ID_IS_FULL_BOX_MAP = {
  'trun': true,
  'tfhd': true
};

/**
 * Checks to see if the id is a full box element.
 *
 * @param {string} id The id to check.
 * @return {boolean} True if the id is a full box element.
 * False otherwise.
 */
ISOBMFFValidator.prototype.isIdAFullBox = function(id) {
  return ID_IS_FULL_BOX_MAP[id] || false;
};

/**
 * @override
 */
ISOBMFFValidator.prototype.onListStart = function(id, elementPosition, 
                                                  bodyPosition) {
  console.log('onListStart(' + id +
              ', ' + elementPosition +
              ', ' + bodyPosition + ')');
  return msetools.ParserStatus.OK;
};


/**
 * @override
 */
ISOBMFFValidator.prototype.onListEnd = function(id, size) {
  console.log('onListEnd(' + id + ', ' + size + ')');
  return true;
};


/**
 * @override
 */
ISOBMFFValidator.prototype.onBinary = function(id, value) {
  if (this.isIdAFullBox(id)) {
    if (value.length < 4) {
      console.log('Invalid FullBox \'' + id + '\'');
      return false;
    }
    var tmp = this.getUint32_(value);
    var version = (tmp >> 24) & 0xff;
    var flags = tmp & 0xffffff;
    return this.onFullBox(id, version, flags, value.subarray(4));
  }

  console.log('onBinary(' + id + ', ' + value.length + ')');

  return true;
};

/**
 * Called when a full box element has been received.
 *
 * @param {string} id Element id.
 * @param {number} version The full box version field.
 * @param {number} flags The full box flags field.
 * @param {Uint8Array} value The body of the full box.
 * @return {boolean} True if the element was successfully parsed.
 */
ISOBMFFValidator.prototype.onFullBox = function(id, version, flags, value) {
  console.log('onFullBox(' + id +
              ', ' + version +
              ', 0x' + flags.toString(16) +
              ', ' + value.length + ')');

  if (id == 'trun') {
    return this.parseTrun(version, flags, value);
  } else if (id == 'tfhd') {
    return this.parseTfhd(version, flags, value);
  }

  return true;
};

/**
 * Called when trun box has been received.
 *
 * @param {number} version The full box version field.
 * @param {number} flags The full box flags field.
 * @param {Uint8Array} value The body of the full box.
 * @return {boolean} True if the element was successfully parsed.
 */
ISOBMFFValidator.prototype.parseTrun = function(version, flags, value) {
  var hasDataOffset = (flags & 0x1) != 0;
  var hasFirstSampleFlag = (flags & 0x4) != 0;
  var hasSampleDuration = (flags & 0x100) != 0;
  var hasSampleSize = (flags & 0x200) != 0;
  var hasSampleFlags = (flags & 0x400) != 0;
  var hasSampleCompositionOffsets = (flags & 0x800) != 0;

  var sampleCount = this.getUint32_(value.subarray(0));
  console.log('trun.sample_count ' + sampleCount);
  var i = 4;
  if (hasDataOffset) {
    var offset = this.getUint32_(value.subarray(i));
    console.log('trun.data_offset ' + offset);
    i += 4;
  }

  var firstSampleFlags = -1;
  if (hasFirstSampleFlag) {
    firstSampleFlags = this.getUint32_(value.subarray(i));
    console.log('trun.first_sample_flags ' +
                       this.sampleFlagsToString_(firstSampleFlags));
    i += 4;
  }

  for (var j = 0; j < sampleCount; ++j) {
    var duration = this.default_sample_duration_;
    var size = this.default_sample_size_;
    var sample_flags = this.default_sample_flags_;
    var compositionOffset = -1;

    if (j == 0 && firstSampleFlags != -1) {
      sample_flags = firstSampleFlags;
    }

    if (hasSampleDuration) {
      duration = this.getUint32_(value.subarray(i));
      i += 4;
    }

    if (hasSampleSize) {
      size = this.getUint32_(value.subarray(i));
      i += 4;
    }

    if (hasSampleFlags) {
      sample_flags = this.getUint32_(value.subarray(i));
      i += 4;
    }

    if (hasSampleCompositionOffsets) {
      compositionOffset = this.getUint32_(value.subarray(i));
      i += 4;
    }
    console.log('trun : ' + duration +
                ' ' + size +
                ' ' + this.sampleFlagsToString_(sample_flags) +
                ' ' + compositionOffset);
  }

  return true;
};

/**
 * Called when tfhd box has been received.
 *
 * @param {number} version The full box version field.
 * @param {number} flags The full box flags field.
 * @param {Uint8Array} value The body of the full box.
 * @return {boolean} True if the element was successfully parsed.
 */
ISOBMFFValidator.prototype.parseTfhd = function(version, flags, value) {
  var hasDataOffset = (flags & 0x1) != 0;
  var hasIndex = (flags & 0x2) != 0;
  var hasDuration = (flags & 0x8) != 0;
  var hasSize = (flags & 0x10) != 0;
  var hasFlags = (flags & 0x20) != 0;
  var isDurationEmpty = (flags & 0x10000) != 0;

  var trackId = this.getUint32_(value.subarray(0));
  var i = 4;
  var offset = -1;
  var index = -1;
  this.default_sample_duration_ = -1;
  this.default_sample_size_ = -1;
  this.default_sample_flags_ = 0;

  if (hasDataOffset) {
    offset = this.getUint64_(value.subarray(i));
    i += 8;
  }

  if (hasIndex) {
    index = this.getUint32_(value.subarray(i));
    i += 4;
  }

  if (hasDuration) {
    this.default_sample_duration_ = this.getUint32_(value.subarray(i));
    i += 4;
  }

  if (hasSize) {
    this.default_sample_size_ = this.getUint32_(value.subarray(i));
    i += 4;
  }

  if (hasFlags) {
    this.default_sample_flags_ = this.getUint32_(value.subarray(i));
    i += 4;
  }

  console.log('tfhd :' +
              ' ' + trackId +
              ' ' + offset +
              ' ' + index +
              ' ' + this.default_sample_duration_ +
              ' ' + this.default_sample_size_ +
              ' ' + this.sampleFlagsToString_(this.default_sample_flags_));
  return true;
};

/**
 * Conversts a sample_flags field to a string.
 *
 * @param {number} flags The contents of a sample_flags field.
 * @return {string} A string representation of the flags.
 */
ISOBMFFValidator.prototype.sampleFlagsToString_ = function(flags) {
  var str = '[';

  str += ' DO' + ((flags >> 24) & 0x3);
  str += ' IDO' + ((flags >> 22) & 0x3);
  str += ' HR' + ((flags >> 20) & 0x3);
  str += ' P' + ((flags >> 17) & 0x7);
  str += ' D' + ((flags >> 16) & 0x1);
  str += ' PR' + (flags & 0xffff);
  str += ' ]';
  return str;
};

/**
 * Extracts a 32-bit big endian integer from buf.
 *
 * @param {Uint8Array} buf The buffer to read the integer from.
 * @return {number} The 32-bit integer extracted from buf.
 */
ISOBMFFValidator.prototype.getUint32_ = function(buf) {
  var result = 0;
  for (var i = 0; i < 4; ++i) {
    result *= 256;
    result += buf[i] & 0xff;
  }
  return result;
};

/**
 * Extracts a 64-bit big endian integer from buf.
 *
 * @param {Uint8Array} buf The buffer to read the integer from.
 * @return {number} The 64-bit integer extracted from buf.
 */
ISOBMFFValidator.prototype.getUint64_ = function(buf) {
  var result = 0;
  for (var i = 0; i < 8; ++i) {
    result *= 256;
    result += buf[i] & 0xff;
  }
  return result;
};

msetools.ISOBMFFValidator = ISOBMFFValidator;
