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

(function(window, undefined) {
  window['msetools'] = {};
  var nameOfThisFile = 'msetools-dev.js';
  var files = [
    'file.js',
    'element_list_parser.js',
    'webm_validator.js',
    'isobmff_validator.js',
    'validator.js'
  ];

  var basePath = null;
  var head = document.querySelector('head');
  var children = head.childNodes;
  for (var i = 0, max = children.length; i < max; ++i) {
    var child = children[i];
    if (child.nodeName == 'SCRIPT' &&
        (child.src.indexOf(nameOfThisFile) >= 0)) {
      basePath = child.src.replace(nameOfThisFile, '');
      break;
    }
  }

  if (basePath != null) {
    for (var i = 0; i < files.length; ++i) {
      var scriptTag = document.createElement('script');
      scriptTag.type = 'text/javascript';
      scriptTag.src = basePath + files[i];
      head.appendChild(scriptTag);
    }
  }
})(window);
