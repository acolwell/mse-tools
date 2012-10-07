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

  /**
   * @constructor
   * @param {MediaSourceValidator} parent
   * @param {string} id
   * @param {string} type
   * @param {SourceBuffer} sourceBuffer
   */
  function SourceBufferValidator(parent, id, type, sourceBuffer) {
    console.log(id + ': new SourceBuffer(' + type + ')');
    this.parent_ = parent;
    this.id_ = id;
    this.type_ = type;
    this.sourceBuffer_ = sourceBuffer;
    this.bytestreamValidator_ = null;

    var typeInfo = this.parseType_(type);

    if (typeInfo != null) {
      if (typeInfo.minor == 'webm') {
        this.bytestreamValidator_ = new msetools.WebMValidator();
      } else if (typeInfo.minor == 'mp4') {
        this.bytestreamValidator_ = new msetools.ISOBMFFValidator();
      }
    }

    if (this.bytestreamValidator_) {
      this.bytestreamValidator_.init(typeInfo);
    } else {
      console.log(id + ': Unsupported type \'' + type + '\'. ' +
                  ' Validation for this SourceBuffer disabled.');
    }

    var appendFunc = this.sourceBuffer_.append.bind(
      this.sourceBuffer_);

    var abortFunc = this.sourceBuffer_.abort.bind(
      this.sourceBuffer_);

    this.sourceBuffer_.append = this.append.bind(this, appendFunc);
    this.sourceBuffer_.abort = this.abort.bind(this, abortFunc);
  };

  /**
   * @param {function(Uint8Array)} originalMethod
   * @param {Uint8Array} data
   */
  SourceBufferValidator.prototype.append = function(originalMethod, data) {
    //console.log(this.id_ + ': SourceBuffer.append(' + data.length + ')');

    if (this.bytestreamValidator_) {
      var errors = this.parse_(data);
      
      for (var i = 0; i < errors.length; ++i) {
        console.log(errors[i]);
      }
    }

    try {
      originalMethod(data);
    } catch (e) {
      throw e;
    }
  };

  /**
   * @param {function()} originalMethod
   */
  SourceBufferValidator.prototype.abort = function(originalMethod) {
    console.log(this.id_ + ': SourceBuffer.abort()');

    if (this.bytestreamValidator_) {
      this.bytestreamValidator_.abort();
    }

    try {
      originalMethod();
    } catch (e) {
      throw e;
    }
  };

  /**
   * @param {MediaSource.EndOfStreamError=} error
   */
  SourceBufferValidator.prototype.endOfStream = function(error) {
    if (this.bytestreamValidator_) {
      this.bytestreamValidator_.endOfStream();
    }
  }

  SourceBufferValidator.prototype.parseType_ = function(type) {
    var trimmedType = type.replace(' ', '');
    if ((trimmedType.indexOf(';') < 0) ||
        (trimmedType.indexOf(';') == (type.length - 1)))
    {
      console.log('\'' + type + '\' missing codecs.');
      return null;
    }

    if (trimmedType.indexOf(';') < 1) {
      console.log('\'' + type + '\' missing major & minor type.');
      return null;
    }

    var sections = trimmedType.split(';');
    var majorMinor = sections[0].split('/');
    if ((majorMinor.length != 2) ||
        (majorMinor[0].length <= 0) ||
        (majorMinor[1].length <= 0)) {
      console.log('\'' + type + '\' has invalid type & subtype.');
      return null;
    }
    
    var majorType = majorMinor[0].toLowerCase();
    var minorType = majorMinor[1].toLowerCase();
    if ((majorType != 'video') && (majorType != 'audio')) {
      console.log('\'' + type + '\' has unsupported major type \'' +
                  majorType + '\'.');
      return null;
    }

    var codecs = [];
    for (var i = 1; i < sections.length; ++i) {
      var param = sections[i];
      if (param.indexOf('codecs=') != 0) {
        continue;
      }

      param = param.replace('codecs=', '');
      if (param[0] != '"' || param[param.length - 1] != '"') {
        console.log('\'' + type + '\' has codec parameter doesn\'t have \'"\'s.');
        return null;
      }

      codecs = param.substring(1, param.length - 1).split(',');
      for (var i = 0; i < codecs.length; ++i) {
        if (codecs[i].length == 0) {
          console.log('\'' + type + '\' has codec parameter is invalid.');
          return null;
        }
      }
    }

    if (codecs.length == 0) {
      console.log('\'' + type + '\' has no codecs specified.');
      return null;
    }

    return { major: majorType, minor: minorType, codecs: codecs };
  }

  SourceBufferValidator.prototype.parse_ = function(data) {
    return this.bytestreamValidator_.parse(data);
  };

  /**
   * @constructor
   * @param {string} id The unique ID for this validator.
   * @param {MediaSource} mediaSource The MediaSource object to attach the
   * validator to.
   */
  function MediaSourceValidator(id, mediaSource) {
    this.id_ = id;
    this.mediaSource_ = mediaSource;
    this.sourceBuffers_ = [];
    this.nextSourceBufferId_ = 0;

    var addSourceBufferFunc = this.mediaSource_.addSourceBuffer.bind(
      this.mediaSource_);
    var removeSourceBufferFunc = this.mediaSource_.removeSourceBuffer.bind(
      this.mediaSource_);
    var endOfStreamFunc = this.mediaSource_.endOfStream.bind(
      this.mediaSource_);

    this.mediaSource_.addSourceBuffer = this.addSourceBuffer.bind(
      this, addSourceBufferFunc);
    this.mediaSource_.removeSourceBuffer = this.removeSourceBuffer.bind(
      this, removeSourceBufferFunc);
    this.mediaSource_.endOfStream = this.endOfStream.bind(
      this, endOfStreamFunc);
  };

  /**
   * @type {Array.<SourceBufferValidator>}
   */
  MediaSourceValidator.prototype.sourceBuffers_ = null;

  /**
   * @return {MediaSource}
   */
  MediaSourceValidator.prototype.mediaSource = function() {
    return this.mediaSource_;
  }

  /**
   * @param {function(string) : SourceBuffer} originalMethod
   * @param {string} type
   */
  MediaSourceValidator.prototype.addSourceBuffer = function(
    originalMethod, type) {
    console.log(this.id_ + ': MediaSource.addSourceBuffer(' + type + ')');
    var sourceBuffer = null;

    try {
      sourceBuffer = originalMethod(type);
    } catch (e) {
      throw e;
    }

    var id = this.id_ + '-' + this.nextSourceBufferId_;
    this.sourceBuffers_.push(new SourceBufferValidator(
      this, id, type, sourceBuffer));
    this.nextSourceBufferId_++;
    return sourceBuffer;
  };

  /**
   * @param {function(SourceBuffer)} originalMethod
   * @param {SourceBuffer} sourceBuffer
   */
  MediaSourceValidator.prototype.removeSourceBuffer = function(
    originalMethod, sourceBuffer) {
    console.log(this.id_ + ': MediaSource.removeSourceBuffer()');
    var i = this.findSourceBufferIndex(sourceBuffer);

    try {
      originalMethod(sourceBuffer);
    } catch (e) {
      throw e;
    }

    if (i >= 0) {
      // Remove the validator from the list.
      this.sourceBuffers_.splice(i, 1);
    }
  };

  /**
   * @param {function(MediaSource.EndOfStreamError=)} originalMethod
   * @param {MediaSource.EndOfStreamError=} error
   */
  MediaSourceValidator.prototype.endOfStream = function(originalMethod, error) {
    console.log(this.id_ + ': MediaSource.endOfStream(' + 
                ((error != undefined) ? error : '') + ')');

    try {
      if (error == undefined) {
        originalMethod()
      } else {
        originalMethod(error);
      }
    } catch (e) {
      if (e.code ==  DOMException.INVALID_STATE_ERR) {
        console.log(this.id_ + ': MediaSource.endOfStream()' +
                    ' called in unexpected readyState "' + 
                    this.mediaSource_.readyState + '"');
      }
      throw e;
    }

    for (var i = 0; i < this.sourceBuffers_.length; ++i) {
      this.sourceBuffers_[i].endOfStream(error);
    }
  };

  /**
   * @param {SourceBuffer} sourceBuffer
   * @return {number}
   */
  MediaSourceValidator.prototype.findSourceBufferIndex = function(
    sourceBuffer) {
    for (var i = 0; i < this.sourceBuffers_.length; ++i) {
      if (this.sourceBuffers_[i].sourceBuffer_ == sourceBuffer) {
        return i;
      }
    }
    return -1;
  };

  /** @type {number} */
  var nextMediaSourceId = 0;

  /**
   * @param {MediaSource} mediaSource
   */
  msetools.attachValidator = function(mediaSource) {
    var id = nextMediaSourceId.toString();
    nextMediaSourceId++;
    new MediaSourceValidator(id, mediaSource);
  };

})(window, msetools);
