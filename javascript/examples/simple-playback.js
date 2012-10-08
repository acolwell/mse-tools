(function (window, undefined) {
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

    var info;
    info = { url: 'test6.mp4', type: 'video/mp4; codecs="avc1.42E000,mp4a.40.2"'};
    //info = { url: '/webm/bear-320x240.webm', type: 'video/webm; codecs="vorbis,vp8"'};

    var sourceBuffer = mediaSource.addSourceBuffer(info.type);
    var file = new msetools.RemoteFile(info.url);
    var isPlaying = function() { 
      return videoTag.readyState > videoTag.HAVE_FUTURE_DATA;
    };
    appendMoreData = createAppendFunction(mediaSource, sourceBuffer, file,
                                          isPlaying);
    videoTag.addEventListener('progress', onProgress.bind(videoTag, mediaSource));

    appendMoreData();
  }

  function onProgress(mediaSource, e) {
    appendMoreData();
  }

  function onPageLoad() {
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