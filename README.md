Up
====

![](https://github.com/alpn/up/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpn/up)](https://goreportcard.com/report/github.com/alpn/up)

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
