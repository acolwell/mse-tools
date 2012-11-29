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
function ByteStreamValidator() {}


/** @typedef {{major: string, minor: string, codecs: Array.<string>}} */
var ByteStreamTypeInfo;

/**
 * Initializes the validator. Must be called before any other method.
 *
 * @param {ByteStreamTypeInfo} typeInfo Information describing the
 * type of bytestream that will be passed to this validator.
 */
ByteStreamValidator.prototype.init = function(typeInfo) {};

/**
 * Parses new data that has been appended to the bytestream.
 *
 * @param {Uint8Array} data The data that was appended.
 * @return {Array.<string>} List of parse errors detected.
 */
ByteStreamValidator.prototype.parse = function(data) {};

/**
 * Resets the parser state.
 */
ByteStreamValidator.prototype.reset = function() {};

/**
 * Called when the end of the bytestream has been reached.
 */
ByteStreamValidator.prototype.endOfStream = function() {};
