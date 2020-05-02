package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func run(isPipeMode bool, isDirectoryMode bool, bucketName string, args []string) error {

	if isPipeMode && isDirectoryMode {
		return errors.New("directory mode and pipe mode are mutually exclusive")
	}

	if isPipeMode && len(bucketName) == 0 {
		return errors.New("bucket name must be provided in pipe mode")
	}

	ctx, bucket, e := getBucket(bucketName)
	if e != nil {
		return e
	}

	fmt.Println("Bucket: ", bucket.Name())

	if isPipeMode {
		return uploadSTDIN(ctx, bucket, args)
	}

	if isDirectoryMode {

		var root string
		if len(args) == 0 {
			root = "."
		} else if len(args) == 1 {
			root = args[0]
		} else {
			flag.Usage()
		}

		rootAbs, err := filepath.Abs(root)
		if nil != err {
			return err
		}

		f, err := os.Stat(rootAbs)
		if nil != err {
			return err
		}

		if !f.IsDir() {
			return errors.New(rootAbs + " is not a directory")
		}

		fmt.Println("About to recuresively upload diretory -", rootAbs)
		confirmed, err := confirmAction()
		if nil != err {
			return err
		}
		if confirmed {
			return uploadDirectory(ctx, bucket, rootAbs)
		}
		return nil
	}

	// Files mode
	if len(args) > 0 {

		fmt.Println("About to upload the following file[/s]:", args)

		confirmed, err := confirmAction()
		if nil != err {
			return err
		}

		if confirmed {
			return uploadFiles(ctx, bucket, args)
		}
	}
	return nil
}

func main() {

	var isPipeMode bool
	var isDirectoryMode bool
	var bucketName string

	flag.BoolVar(&isPipeMode, "pipe", false, "reads and uploads data from STDIN until EOF is reached. Does NOT ask for confirmation")
	flag.BoolVar(&isDirectoryMode, "dir", false, "recursively uploads an entire directory. if no path is provided, current directory will be assumed")
	flag.StringVar(&bucketName, "bucket", "", "name of destination bucket")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nA simple utility for uploading stuff to BackBlaze's B2\n\n")
		fmt.Fprintf(os.Stderr, "usage: up [-pipe dst_name] [-dir path] [file1 .. fileN]\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	if err := run(isPipeMode, isDirectoryMode, bucketName, args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
