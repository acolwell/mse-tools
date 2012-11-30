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
 * A random access file abstraction built on top of XMLHttpRequest.
 * @param {string} url URL passed to XMLHttpRequest to fetch data.
 * @constructor
 */
function RemoteFile(url) {
  this.url_ = url;
}


/**
 * The minimum number of bytes to request with a single
 * XMLHttpRequest. This is used to avoid hammering the
 * server with small range requests when the application
 * calls read() for only a few bytes at a time.
 * @private
 * @type {number}
 */
RemoteFile.prototype.MIN_REQUEST_SIZE = 1024;


/**
 * The URL passed to XMLHttpRequest to fetch data.
 * @private
 * @type {string}
 */
RemoteFile.prototype.url_ = '';


/**
 * The current read position.
 * @private
 * @type {number}
 */
RemoteFile.prototype.position_ = 0;


/**
 * The size of the file. -1 indicates the size is unknown.
 * @private
 * @type {number}
 */
RemoteFile.prototype.size_ = -1;


/**
 * Buffer for storing downloaded data.
 * @private
 * @type {?Uint8Array}
 */
RemoteFile.prototype.buffer_ = null;


/**
 * The file position of the first byte in buffer_.
 * @private
 * @type {number}
 */
RemoteFile.prototype.bufferPosition_ = -1;


/** @typedef {?function(string, ?Uint8Array)} */
RemoteFile.ReadCallback;


/**
 * Callback passed to read().
 * @private
 * @type {RemoteFile.ReadCallback}
 */
RemoteFile.prototype.readCallback_ = null;


/**
 * Regular expression used to extract the file size from the
 * Content-Range header in the response.
 * @private
 * @const
 * @type {RegExp}
 */
RemoteFile.ContentRangeRegEx_ = /^bytes \d+-\d+\/(\d+)$/;


/**
 * Gets the current read position.
 * @return {number} The current read position.
 */
RemoteFile.prototype.getPosition = function() {
  return this.position_;
};


/**
 * Gets the size of the file.
 * @return {number} The size of the file if known. Returns -1 if size is not
 *  known.
 */
RemoteFile.prototype.getSize = function() {
  return this.size_;
};


/**
 * Checks if the current position has reached the end of the file.
 * @return {boolean} True if current position is at the end of the file.
 */
RemoteFile.prototype.isEndOfFile = function() {
  return this.size_ == this.position_;
};


/**
 * Seek to a new position.
 * @param {number} newPosition The new position to seek to.
 */
RemoteFile.prototype.seek = function(newPosition) {
  this.position_ = newPosition;
  if (!this.buffer_)
    return;

  var start = this.bufferPosition_;
  var end = start + this.buffer_.length;
  if (this.position_ >= start && this.position_ < end)
    return;
  this.buffer_ = null;
  this.bufferPosition_ = -1;
};


/**
 * Reads data from the current position.
 * @param {number} size The number of bytes to read.
 * @param {RemoteFile.ReadCallback} doneCallback The callback to run
 * when the read has completed.
 */
RemoteFile.prototype.read = function(size, doneCallback) {
  if (this.readCallback_) {
    throw 'Last read() still pending';
  }

  this.readCallback_ = doneCallback;
  this.doRead_(size);
};


/**
 * Attempts to fulfill a read() request.
 * @param {number} size The number of bytes to read.
 * @private
 */
RemoteFile.prototype.doRead_ = function(size) {
  var readEnd = this.position_ + size;

  if (this.size_ != -1) {
    if (this.isEndOfFile()) {
      this.runCallback_('eof', null);
      return;
    }

    if (readEnd > this.size_)
      readEnd = this.size_;
  }

  var downloadPosition = -1;
  if (this.buffer_) {
    var bufferEnd = this.bufferPosition_ + this.buffer_.length;
    if (readEnd > bufferEnd)
      downloadPosition = bufferEnd;
  } else {
    downloadPosition = this.position_;
  }

  var readSize = readEnd - this.position_;

  if (downloadPosition != -1) {
    var downloadSize = readEnd - downloadPosition;
    this.downloadData_(readSize, downloadPosition, downloadSize);
    return;
  }

  var offset = this.position_ - this.bufferPosition_;
  var buffer = this.buffer_.subarray(offset, offset + readSize);
  this.position_ += readSize;
  if ((this.position_ - this.bufferPosition_) >= this.buffer_.length) {
    this.buffer_ = null;
    this.bufferPosition_ = -1;
  }
  this.runCallback_('ok', buffer);
};


/**
 * Runs & resets readCallback_.
 * @param {string} status The read status.
 * @param {?Uint8Array} buffer The bytes read if status is 'ok'. Null otherwise.
 * @private
 */
RemoteFile.prototype.runCallback_ = function(status, buffer) {
  var callback = this.readCallback_;
  this.readCallback_ = null;
  callback(status, buffer);
};


/**
 * Downloads more data using XMLHttpRequest.
 * @param {number} readSize The size passed to read().
 * @param {number} downloadPosition The starting position for the download.
 * @param {number} downloadSize The number of bytes to download.
 * @private
 */
RemoteFile.prototype.downloadData_ = function(readSize, 
                                              downloadPosition,     
                                              downloadSize) {
  var requestSize = this.MIN_REQUEST_SIZE;
  if (requestSize < downloadSize)
    requestSize = downloadSize;

  var start = downloadPosition;
  var end = downloadPosition + requestSize - 1;
  var xhr = new XMLHttpRequest();
  xhr.open('GET', this.url_, true);
  xhr.setRequestHeader('Range', 'bytes=' + start + '-' + end);
  xhr.responseType = 'arraybuffer';

  xhr.onload = this.onLoad_.bind(this, readSize);
  xhr.onerror = this.onError_.bind(this, readSize);
  xhr.send();
};


/**
 * Called when the onload event fires for request issued by downloadData_().
 * @param {number} readSize The size passed to the read() call that triggered
 * the download.
 * @param {ProgressEvent} event The event object passed to the handler.
 * @private
 */
RemoteFile.prototype.onLoad_ = function(readSize, event) {
  var xhr = (/** @type {XMLHttpRequest} */ event.target);
  if (xhr.status != 206 && xhr.status != 0) {
    this.downloadDataDone_(readSize, null);
    return;
  }

  var buffer = new Uint8Array((/** @type {ArrayBuffer} */xhr.response));
  if (this.size_ == -1) {
    var range = xhr.getResponseHeader('Content-Range');
    var m = RemoteFile.ContentRangeRegEx_.exec(range);
    this.size_ = (m && m.length >= 2) ? parseInt(m[1], 10) : -1;
  }
  this.downloadDataDone_(readSize, buffer);
};


/**
 * Called when the onerror event fires for request issued by downloadData_().
 * @param {number} readSize The size passed to the read() call that triggered
 * the download.
 * @param {ProgressEvent} event The event object passed to the handler.
 * @private
 */
RemoteFile.prototype.onError_ = function(readSize, event) {
  this.downloadDataDone_(readSize, null);
};


/**
 * Called when a download has completed.
 * @param {number} readSize The size passed to the read() call that triggered
 * the download.
 * @param {?Uint8Array} buffer The buffer that was downloaded. Null if an
 * error occurred.
 * @private
 */
RemoteFile.prototype.downloadDataDone_ = function(readSize, buffer) {
  if (!buffer) {
    this.runCallback_('error', buffer);
    return;
  }

  if (!this.buffer_) {
    this.buffer_ = buffer;
    this.bufferPosition_ = this.position_;
  } else {
    var oldBuffer = this.buffer_;
    this.buffer_ = new Uint8Array(oldBuffer.length + buffer.length);
    this.buffer_.set(oldBuffer, 0);
    this.buffer_.set(buffer, oldBuffer.length);
  }

  this.doRead_(readSize);
};

msetools.RemoteFile = RemoteFile;
