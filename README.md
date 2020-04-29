Up
====

Up is a simple utility for uploading stuff to BackBlaze's B2 cluod storage service.
It's merely a thin wrapper around [Blazer](https://github.com/kurin/blazer) by @kurin.

## Build
```bash
git clone https://github.com/alpn/up.git
cd up
go build
```
## Usage

### Files
```bash
./up file1 file2 ..
```
### Directories

```bash

./up -dir . 

```
### Upload STDIN

```bash

zfs send snapshot_342 | ./up -pipe

```
