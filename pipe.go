package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kurin/blazer/b2"
	"github.com/ttacon/chalk"
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

	fmt.Printf("Uploading STDIN to %sBackBlaze B2%s cloud storage:\n\n", chalk.Red, chalk.Reset)

	fileDisplayString := fmt.Sprintf("%s%s%s", chalk.Blue, fileName, chalk.Reset)
	if err := uploadOneReader(ctx, bucket, r, fileDisplayString, fileName, true); nil != err {
		return err
	}

	dstFull := getDstPathString(bucket.Name(), fileName)
	fmt.Println("EOF received. File has been uploaded to BackBlaze => " + dstFull)
	return nil
}
