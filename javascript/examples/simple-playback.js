(function(window, undefined) {
  var isFirstOpen = true;

  var appendMoreData = null;

  function createAppendFunction(mediaSource, sourceBuffer, file, isPlaying) {
    var readPending = false;

    var readDone = function(status, buf) {
      readPending = false;
      if (status == 'error') {
        mediaSource.endOfStream('network');
        return;
      } else if (status == 'eof') {
        mediaSource.endOfStream(error);
        return;
      }
      sourceBuffer.append(buf);

      if (isPlaying())
        return;

      appendMoreData();
    };

    return function() {
      if (readPending)
        return;

      if (file.isEndOfFile()) {
        if (mediaSource.readyState != 'ended') {
          mediaSource.endOfStream();
        }

        console.log('No more data to append');
        return;
      }

      if (mediaSource.readyState == 'ended') {
        console.log('mediaSource already ended.');
        return;
      }

      readPending = true;
      file.read(128 * 1024, readDone);
    }
  }

  function onSourceOpen(videoTag, e) {
    var mediaSource = e.target;

    if (!isFirstOpen) {
      appendMoreData();
      return;
    }

    isFirstOpen = false;

    var url = document.getElementById('u').value;
    var codecs = document.getElementById('c').value;

    var type = '';
    if (codecs.indexOf('avc1.') != -1 || codecs.indexOf('avc1.') != -1) {
      type = 'video/mp4; codecs="' + codecs + '"';
    } else if (codecs.indexOf('vp8') != -1 || codecs.indexOf('vorbis') != -1) {
      type = 'video/webm; codecs="' + codecs + '"';
    }

    if (type.length == 0) {
      console.log('Couldn\'t determine type from codec string "' +
	  codecs + '"');
      return;
    }

    var info = { url: url, type: type};
    var sourceBuffer = mediaSource.addSourceBuffer(info.type);
    var file = new msetools.RemoteFile(info.url);
    var isPlaying = function() {
      return videoTag.readyState > videoTag.HAVE_FUTURE_DATA;
    };
    appendMoreData = createAppendFunction(mediaSource, sourceBuffer, file,
                                          isPlaying);
    videoTag.addEventListener('progress', onProgress.bind(videoTag,
							  mediaSource));

    appendMoreData();
  }

  function onProgress(mediaSource, e) {
    appendMoreData();
  }

  function onPageLoad() {
    document.getElementById('b').addEventListener('click', loadUrl);

    var loadURL = false;

    // Extract the 'url' parameter from the document URL.
    var urlRegex = new RegExp('[\\?&]url=([^&#]*)');
    var codecsRegex = new RegExp('[\\?&]codecs=([^&#]*)');
    var results = urlRegex.exec(window.location.href);
    if (results != null) {
      var url = results[1];

      // Assign to the input field.
      var u = document.getElementById('u');
      u.value = url;
    }

    results = codecsRegex.exec(window.location.href);
    if (results != null) {
      var codecs = results[1];

      // Assign to the input field.
      var c = document.getElementById('c');
      c.value = codecs;
      loadURL = true;
    }

    if (loadURL) {
      loadUrl();
    }
  }

  function loadUrl() {
    window.MediaSource = window.MediaSource || window.WebKitMediaSource;

    var video = document.getElementById('v');
    var mediaSource = new MediaSource();

    msetools.attachValidator(mediaSource);
    mediaSource.addEventListener('webkitsourceopen',
                                 onSourceOpen.bind(this, video));
    video.src = window.URL.createObjectURL(mediaSource);
  }

  window['onPageLoad'] = onPageLoad;
})(window);
