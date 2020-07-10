package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kurin/blazer/b2"
	"github.com/ttacon/chalk"
)

func uploadSingleFile(ctx context.Context, bucket *b2.Bucket, path string) error {

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
	fileString := fmt.Sprintf("%s/%s%s %s(%s)%s", dir, chalk.Yellow, fileName, chalk.Blue, formmatedSize, chalk.Reset)

	if err = uploadSingleReader(ctx, bucket, f, fileString, path, false); nil != err {
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
			if err = uploadSingleFile(ctx, bucket, path); nil != err {
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

func uploadFiles(ctx context.Context, bucket *b2.Bucket, files []string, allowDirectories bool) error {

	for _, f := range files {

		path, err := filepath.Abs(f)
		if err != nil {
			return err
		}

		fs, err := os.Stat(path)
		if nil != err {
			if os.IsNotExist(err) {
				fmt.Printf("%s%s%s does not exist\n",
					chalk.Yellow, f, chalk.Reset)
				continue
			}
			return err
		}

		if fs.IsDir() {
			if allowDirectories {
				return uploadDirectory(ctx, bucket, path)
			}
			fmt.Printf("%s%s%s is a directory, but '-dir' was not provided - skipping\n",
				chalk.Yellow, f, chalk.Reset)

		} else if err = uploadSingleFile(ctx, bucket, path); nil != err {
			return err
		}

	}

	return nil
}

func handleFiles(bucketName string, files []string, allowDirectories bool) error {

	ctx, buckets, err := getAllBuckets()
	if nil != err {
		return err
	}

	bucket, err := pickBucketPrompt(buckets, bucketName)
	if nil != err {
		return err
	}

	printBucket(bucket.Name())
	fmt.Printf("Uploading files to %sBackBlaze B2%s cloud storage:\n\n", chalk.Red, chalk.Reset)

	return uploadFiles(ctx, bucket, files, allowDirectories)

}
