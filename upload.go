package main

import (
	"context"
	"fmt"
	"io"

	"github.com/kurin/blazer/b2"
)

func uploadOneReader(ctx context.Context, bucket *b2.Bucket, src io.Reader, dstName string) error {

	obj := bucket.Object(dstName)
	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, src); nil != err {
		writer.Close()
		return err
	}
	return writer.Close()
}
