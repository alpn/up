package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kurin/blazer/b2"
)

func uploadFile(ctx context.Context, bucket *b2.Bucket, path string) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	fmt.Println(path, fi.Size())

	dstName := strings.Trim(path, "/") //filepath.Base(path)

	if err = uploadOneReader(ctx, bucket, f, dstName); nil != err {
		return err
	}

	fmt.Println("=> ", getDstPathString(bucket, dstName))

	return nil
}

func uploadDirectory(rootAbs string) error {

	ctx, bucket, err := getBucket()
	if nil != err {
		return err
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if err = uploadFile(ctx, bucket, path); nil != err {
				return err
			}
		}
		return nil
	}

	err = filepath.Walk(rootAbs, walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func uploadFiles(files []string) error {

	ctx, bucket, err := getBucket()
	if nil != err {
		return err
	}

	for _, f := range files {

		path, err := filepath.Abs(f)
		if err != nil {
			return err
		}

		if err = uploadFile(ctx, bucket, path); nil != err {
			return err
		}

	}

	return nil
}
