package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kurin/blazer/b2"
	"github.com/ttacon/chalk"
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

	formmatedSize := ByteCountSI(fi.Size())
	fileName := filepath.Base(path)
	dir := filepath.Dir(path)
	fileString := fmt.Sprintf("%s/%s%s %s(%s)%s", dir, chalk.Blue, fileName, chalk.Yellow, formmatedSize, chalk.Reset)

	if err = uploadOneReader(ctx, bucket, f, fileString, path, false); nil != err {
		return err
	}

	return nil
}

func uploadDirectory(ctx context.Context, bucket *b2.Bucket, rootAbs string) error {

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

	err := filepath.Walk(rootAbs, walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func uploadFiles(ctx context.Context, bucket *b2.Bucket, files []string) error {

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
