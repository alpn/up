package main

import (
	"context"
	"io"
	"sync"

	"github.com/kurin/blazer/b2"
)

func uploadOneReader(ctx context.Context, bucket *b2.Bucket, src io.Reader, dstName string) error {

	obj := bucket.Object(dstName)
	writer := obj.NewWriter(ctx)

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer func() {
		stop <- true
		wg.Wait()
	}()

	go showProgress(stop, &wg, bucket.Name(), dstName)

	if _, err := io.Copy(writer, src); nil != err {
		writer.Close()
		return err
	}

	return writer.Close()
}
