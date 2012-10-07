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
 * @param {Window} window
 * @param {Object} msetools MSE Tools module.
 * @param {Object=} undefined
 */
(function(window, msetools, undefined) {
  var console = window.console;

  function ISOBMFFValidator() {
    this.parser_ = new msetools.ElementListParser(this);
    this.parserError_ = false;
    this.errors_ = [];
  }

  ISOBMFFValidator.prototype.init = function(typeInfo) {

  };

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

  ISOBMFFValidator.prototype.abort = function(data) {

  };

  ISOBMFFValidator.prototype.endOfStream = function() {

  };

  ISOBMFFValidator.prototype.parseElementHeader = function(buf) {
    if (buf.length < 8) {
      return {status: msetools.ParserStatus.NEED_MORE_DATA, 
              bytesUsed: 0, id: 0, size: 0 };
    }

    var bytesUsed = 0;
    var size = 0;
    for (var i = 0; i < 4; ++i) {
      size *= 256;
      size += buf[i];
      bytesUsed++;
    }

    var id = '';
    for (var i = 4; i < 8; ++i) {
      id += String.fromCharCode(buf[i]);
      bytesUsed++;
    }

    if (id == 'uuid') {
      this.errors_.push('uuid boxes not supported yet!.');
      return {status: msetools.ParserStatus.ERROR, 
              bytesUsed: 0, id: 0, size: 0 };
    }

    if (size == 0) {
      this.errors_.push('Box size of 0 supported yet!.');
      return {status: msetools.ParserStatus.ERROR, 
              bytesUsed: 0, id: 0, size: 0 };
    }

    if (size < bytesUsed) {
      this.errors_.push('Invalid box size ' + size);
      return {status: msetools.ParserStatus.ERROR, 
              bytesUsed: 0, id: 0, size: 0 };
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
  }

  var ID_IS_LIST_MAP = {
    'moov': true,
    'moof': true,
    'traf': true,
    'trak': true,
    'mdia': true,
    'minf': true,
    'stbl': true
  };
  ISOBMFFValidator.prototype.isIdAList = function(id) {
    return ID_IS_LIST_MAP[id] || false;
  }

  ISOBMFFValidator.prototype.onListStart = function(id, elementPosition, 
                                                 bodyPosition) {
    console.log('onListStart(' + id + 
                ', ' + elementPosition +
                ', ' + bodyPosition + ')');
    return msetools.ParserStatus.OK;
  };


  /**
   * Called when the end of a list element is parsed.
   * @param {number} id The ID for the list element.
   * @param {number} size The size of the list.
   * @return {boolean} True if the element was accepted by the client.
   * False if the client wants the parser to signal a parse error.
   */
  ISOBMFFValidator.prototype.onListEnd = function(id, size) {
    console.log('onListEnd(' + id + 
                ', ' + size + ')');
    return true;
  };


  /**
   * Called when a binary element is parsed.
   * @param {number} id The ID for the element.
   * @param {Uint8Array} value The value in the element.
   * @return {boolean} True if the element was accepted by the client.
   * False if the client wants the parser to signal a parse error.
   */
  ISOBMFFValidator.prototype.onBinary = function(id, value) {
    console.log('onBinary(' + id + ', ' + value.length + ')');
    return true;
  };

  msetools.ISOBMFFValidator = ISOBMFFValidator;
})(window, msetools);
