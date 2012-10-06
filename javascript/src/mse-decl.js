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
 * @interface
 */
function MediaSource() {}

/**
 * @type {Array.<SourceBuffer>}
 */
MediaSource.prototype.sourceBuffers;

/**
 * @type {Array.<SourceBuffer>}
 */
MediaSource.prototype.activeSourceBuffers;

/**
 * @type {number}
 */
MediaSource.prototype.duration;

/**
 * @param {string} type
 * @return {SourceBuffer}
 */
MediaSource.prototype.addSourceBuffer = function(type) {};

/**
 * @param {SourceBuffer} sourceBuffer
 */
MediaSource.prototype.removeSourceBuffer = function(sourceBuffer) {};

/**
 * @enum {string}
 */
MediaSource.State = {
  closed: 'closed',
  open: 'open',
  ended: 'ended'
};

/**
 * @type {MediaSource.State}
 */
MediaSource.prototype.readyState;


/**
 * @enum {string}
 */
MediaSource.EndOfStreamError = {
  network: 'network',
  decode: 'decode'
};

/**
 * @param {MediaSource.EndOfStreamError=} error
 */
MediaSource.prototype.endOfStream = function(error) {};

/**
 * @interface
 */
function SourceBuffer() {}

/**
 * @type {TimeRanges}
 */
SourceBuffer.prototype.buffered;

/**
 * @type {number}
 */
SourceBuffer.prototype.timestampOffset;

/**
 * @param {Uint8Array} data
 */
SourceBuffer.prototype.append = function(data) {};

SourceBuffer.prototype.abort = function() {};
