mse-tools
=========================================

Go tools that simplify building applications for the 
[Media Source Extensions](http://dvcs.w3.org/hg/html-media/raw-file/tip/media-source/media-source.html).

## Go Command-line Tools
### Tools
* mse\_webm\_remuxer - Remuxes a WebM file so it conforms to [WebM Byte Stream](http://dvcs.w3.org/hg/html-media/raw-file/tip/media-source/media-source.html#webm) requirements. 
* mse\_json\_manifest - Generates a simple JSON manifest that contains information about the initialization segment and media segments in a WebM file.
* webm\_dump - Simple debugging tool that dumps the element information in a WebM file.

### Requirements
* [Go](http://golang.org/)
* [Mercurial](http://mercurial.selenic.com/)

### Build
- Setup your $GOPATH as described in [How to Write Go Code](http://golang.org/doc/code.html)
- Run the following commands:
    <pre>`cd $GOPATH
    go get github.com/acolwell/mse-tools/...
    `</pre>
- All the command-line tools will appear in $GOPATH/bin
