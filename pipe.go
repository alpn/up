package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kurin/blazer/b2"
)

func uploadSTDIN(ctx context.Context, bucket *b2.Bucket, files []string) error {

	// TODO: can't budio.NewReader fail?
	// https://github.com/golang/go/issues/14162
	r := bufio.NewReader(os.Stdin)

	var fileName string

	if len(files) > 0 {
		fileName = files[0]
	} else {
		t := time.Now().UTC()
		fileName = "stdin_" + t.Format("2006-01-02-T15:04:05.9999")
	}

	fmt.Println("Uploading STDIN to BackBlaze..")

	if err := uploadOneReader(ctx, bucket, r, fileName); nil != err {
		return err
	}

	dstFull := getDstPathString(bucket.Name(), fileName)
	fmt.Println("EOF received. File has been uploaded to BackBlaze => " + dstFull)
	return nil
}
