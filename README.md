Up ⤴️
====

![](https://github.com/alpn/up/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpn/up)](https://goreportcard.com/report/github.com/alpn/up)

Up is a simple utility for uploading stuff to BackBlaze's B2 cluod storage service.
It's merely a thin wrapper around [Blazer](https://github.com/kurin/blazer) by @kurin.

## Build / Install
```bash
git clone https://github.com/alpn/up.git
cd up
go build

# optionally, move the binary to a PATH directory, e.g
mv up /usr/local/bin
```

## Usage

### Uploading files and directories
```bash
# upload files and *ignore* directories
up file1 file2 ..

# upload files and directories
up -dir file1 file2 dir1 dir2 ..
```
![](https://github.com/alpn/up/raw/master/.media/files.gif)

### Uploading pipes

```bash

curl -s https://example.com/file.dat | up -pipe -bucket BUCKET_NAME
```
![](https://github.com/alpn/up/raw/master/.media/pipe.gif)

## Credits
* [kurin/blazer](https://github.com/kurin/blazer)
*	[vbauerster/mpb](https://github.com/vbauerster/mpb)
* [ttacon/chalk](https://github.com/ttacon/chalk)

## License
MIT

