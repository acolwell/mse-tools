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
 * Parser status enum.
 * @enum {string}
 */
msetools.ParserStatus = {
  OK: 'ok',
  NEED_MORE_DATA: 'need_more_data',
  ERROR: 'error'
};


/**
 * Element reserved ID constant.
 * @const
 * @type {string}
 */
msetools.RESERVED_ID = '_RESERVED_ID_';


/**
 * Unknown element size constant.
 * @const
 * @type {number}
 */
msetools.UNKNOWN_SIZE = -1;


/**
 * Client interface for msetools.ElementListParser.
 * @interface
 */
msetools.ParserClient = function() {};

/**
 * Parses an element header.
 * @param {Uint8Array} buf The buffer to parse.
 * @return {{status: msetools.ParserStatus, bytesUsed: number,
 * id: string, size: number}}
 */
msetools.ParserClient.prototype.parseElementHeader = function(buf) {};

/**
 * Checks to see if an element ID is corresponds to a list element.
 * @param {string} id An element ID.
 * @return {boolean} True if id is a list element.
 */
msetools.ParserClient.prototype.isIdAList = function(id) {};


/**
 * Called when the start of a list element is parsed.
 * @param {string} id The ID for the list element.
 * @param {number} elementPosition The position of list element header.
 * @param {number} bodyPosition The position of list element body.
 * @return {msetools.ParserStatus} True if the element was accepted by
 * the client. False if the client wants the parser to signal a parse error.
 */
msetools.ParserClient.prototype.onListStart =
  function(id, elementPosition, bodyPosition) {};


/**
 * Called when the end of a list element is parsed.
 * @param {string} id The ID for the list element.
 * @param {number} size The size of the list.
 * @return {boolean} True if the element was accepted by the client.
 * False if the client wants the parser to signal a parse error.
 */
msetools.ParserClient.prototype.onListEnd = function(id, size) {};


/**
 * Called when a binary element is parsed.
 * @param {string} id The ID for the element.
 * @param {Uint8Array} value The value in the element.
 * @return {boolean} True if the element was accepted by the client.
 * False if the client wants the parser to signal a parse error.
 */
msetools.ParserClient.prototype.onBinary = function(id, value) {};


/**
 * Parses list elements.
 * @constructor
 * @param {msetools.ParserClient} client The client to notify about parsed
 * elements.
 */
msetools.ElementListParser = function(client) {
  this.client_ = client;
  this.listStack_ = [];
};


/**
 * Parser buffer.
 * @type {?Uint8Array}
 * @private
 */
msetools.ElementListParser.prototype.buffer_ = null;


/**
 * Stack of the list elements that are currently being parsed.
 * @type {Array.<{id: string, size: number, bytes_left: number,
 * startPosition: number }>}
 * @private
 */
msetools.ElementListParser.prototype.listStack_ = null;


/**
 * Client to notify when element are parsed.
 * @type {?msetools.ParserClient}
 * @private
 */
msetools.ElementListParser.prototype.client_ = null;


/**
 * The number of bytes parsed so far.
 * @type {number}
 * @private
 */
msetools.ElementListParser.prototype.bytePosition_ = 0;


/**
 * Reset the parser state.
 * @param {number} position The parser byte position.
 */
msetools.ElementListParser.prototype.reset = function(position) {
  this.buffer_ = null;
  this.listStack_ = [];
  this.bytePosition_ = position;
};


/**
 * Append data to the parser buffer.
 * @param {Uint8Array} newBuffer The buffer to append to the parse buffer.
 * @return {msetools.ParserStatus} The status of the parse.
 */
msetools.ElementListParser.prototype.append = function(newBuffer) {
  if (this.buffer_) {
    var oldBuffer = this.buffer_;
    this.buffer_ = new Uint8Array(oldBuffer.length + newBuffer.length);
    this.buffer_.set(oldBuffer, 0);
    this.buffer_.set(newBuffer, oldBuffer.length);
  } else {
    this.buffer_ = newBuffer;
  }

  /** @type {number} */ var i = 0;
  var status = msetools.ParserStatus.OK;
  while (i < this.buffer_.length) {
    var buf = this.buffer_.subarray(i, this.buffer_.length);
    var res = this.client_.parseElementHeader(buf);
    if (res.status == msetools.ParserStatus.ERROR)
      return msetools.ParserStatus.ERROR;

    if (res.status == msetools.ParserStatus.NEED_MORE_DATA) {
      status = msetools.ParserStatus.NEED_MORE_DATA;
      break;
    }

    //console.log('pos ' + this.bytePosition_.toString(16) + ' id ' + res.id);
//    /** @type {{status: msetools.ParserStatus, bytesUsed: number, id: string,
//        size: number}} */(res);
    if (this.client_.isIdAList(res.id)) {
      //console.log('msetools.ElementListParser.append : ' +
      //   ' id ' + res.id +
      //   ' size ' + res.size +
      //   ' list start');
      var wholeElementSize = res.bytesUsed + res.size;

      if (res.size == msetools.UNKNOWN_SIZE)
        wholeElementSize = msetools.UNKNOWN_SIZE;

      var elementPosition = this.bytePosition_;
      var bodyPosition = elementPosition + res.bytesUsed;
      var startRes = this.client_.onListStart(res.id, elementPosition,
                                              bodyPosition);

      if (startRes != msetools.ParserStatus.OK &&
          startRes != msetools.ParserStatus.NEED_MORE_DATA) {
        return msetools.ParserStatus.ERROR;
      }

      this.listStack_.push({
        id: res.id,
        startPosition: this.bytePosition_,
        size: wholeElementSize,
        bytes_left: res.size});
      i += res.bytesUsed;
      this.bytePosition_ += res.bytesUsed;

      if (res.size == 0)
        this.handleListEnd_(0);

      if (startRes == msetools.ParserStatus.NEED_MORE_DATA) {
        status = msetools.ParserStatus.NEED_MORE_DATA;
        break;
      }

      continue;
    }

    if (res.size == msetools.UNKNOWN_SIZE)
      return msetools.ParserStatus.ERROR;

    var wholeElementSize = res.bytesUsed + res.size;
    if (buf.length < wholeElementSize) {
      status = msetools.ParserStatus.NEED_MORE_DATA;
      break;
    }

    var elementBody = buf.subarray(res.bytesUsed,
                                   res.bytesUsed + res.size);

    if (!this.client_.onBinary(res.id, elementBody))
      return msetools.ParserStatus.ERROR;

    i += wholeElementSize;
    this.bytePosition_ += wholeElementSize;
    if (this.handleListEnd_(wholeElementSize) != msetools.ParserStatus.OK) {
      return msetools.ParserStatus.ERROR;
    }

  }
  this.buffer_ = this.buffer_.subarray(i, this.buffer_.length);

  if (status == msetools.ParserStatus.OK &&
      (this.listStack_.length > 0 || this.buffer_.length > 0))
    status = msetools.ParserStatus.NEED_MORE_DATA;
  return status;
};


/**
 * Helper function called after each element is parsed to detect when the end
 *  of a list element is reached.
 * @param {number} bytesUsed The number of bytes parsed in the last element.
 * @return {msetools.ParserStatus} The current parser status.
 * @private
 */
msetools.ElementListParser.prototype.handleListEnd_ = function(bytesUsed) {
  var listBytesUsed = bytesUsed;
  while (this.listStack_.length > 0) {
    var li = this.listStack_[this.listStack_.length - 1];
    if (li.size == msetools.UNKNOWN_SIZE) {
      //break;
      //console.log('List with unknown size not supported : ' + li.id);
      return msetools.ParserStatus.ERROR;
    }

    li.bytes_left -= listBytesUsed;

    if (li.bytes_left > 0)
      break;

    if (li.bytes_left < 0) {
      return msetools.ParserStatus.ERROR;
    }

    var listSize = this.bytePosition_ - li.startPosition;
    if (!this.client_.onListEnd(li.id, listSize)) {
      return msetools.ParserStatus.ERROR;
    }

    //console.log('msetools.ElementListParser.handleListEnd_ : ' +
    //               ' id ' + li.id +
    //       ' list end');
    this.listStack_.pop();
    listBytesUsed = li.size;
  }
  return msetools.ParserStatus.OK;
};
