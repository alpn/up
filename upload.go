package main

import (
	"context"
	"io"
	"sync"

	"github.com/kurin/blazer/b2"
)

type nopSeekWriter struct {
	// TODO: can this be *io.Writer instead?
	writer *b2.Writer
}

func (uw *nopSeekWriter) Write(p []byte) (int, error) {
	return (*uw.writer).Write(p)
}

func uploadSingleReader(ctx context.Context, bucket *b2.Bucket,
	src io.Reader, fileDisplayString string, dstName string, isPipe bool) error {

	obj := bucket.Object(dstName)
	writer := obj.NewWriter(ctx)
	dst := &nopSeekWriter{writer: writer}

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer func() {
		stop <- true
		wg.Wait()
	}()

	go showProgress(stop, &wg, bucket.Name(), fileDisplayString, dstName, isPipe)

	if _, err := io.Copy(dst, src); nil != err {
		writer.Close()
		return err
	}

	return writer.Close()
}
