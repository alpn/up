package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func uploadSTDIN(files []string) error {

	ctx, bucket, err := getBucket()
	if nil != err {
		return err
	}

	// TODO: can't budio.NewReader fail?
	// https://github.com/golang/go/issues/14162
	r := bufio.NewReader(os.Stdin)

	var fileName string

	if len(files) > 0 {
		fileName = files[0]
	} else {
		t := time.Now()
		fileName = "stdin_" + t.Format("2020-01-01-T23:59:00.9999")
	}

	fmt.Println("Uploading STDIN to BackBlaze..")

	if err = uploadOneReader(ctx, bucket, r, fileName); nil != err {
		return err
	}

	dstFull := getDstPathString(bucket, fileName)
	fmt.Println("EOF received. File has been uploaded to BackBlaze => " + dstFull)
	return nil
}
