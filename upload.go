package main

import (
	"context"
	"fmt"
	"io"

	"github.com/kurin/blazer/b2"
)

type upWriter struct {
	// TODO: can this be *io.Writer instead?
	writer *b2.Writer
	total  int64
}

func (uw *upWriter) Write(p []byte) (int, error) {

	n, err := (*uw.writer).Write(p)
	uw.total += int64(n)

	if nil == err {
		fmt.Println(uw.total, "bytes written")
	}

	return n, err
}

func uploadOneReader(ctx context.Context, bucket *b2.Bucket, src io.Reader, dstName string) error {

	obj := bucket.Object(dstName)
	writer := obj.NewWriter(ctx)
	dst := &upWriter{writer: writer}

	if _, err := io.Copy(dst, src); nil != err {
		writer.Close()
		return err
	}
	return writer.Close()
}
