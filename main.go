package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func run(isPipeMode bool, isDirectoryMode bool, args []string) error {

	if isPipeMode && isDirectoryMode {
		return errors.New("bad params")
	}

	if isPipeMode {
		return uploadSTDIN(args)
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
			return uploadDirectory(rootAbs)
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
			return uploadFiles(args)
		}
	}
	return nil
}

func main() {

	var isPipeMode bool
	var isDirectoryMode bool

	flag.BoolVar(&isPipeMode, "pipe", false, "reads and uploads data from STDIN until EOF is reached. Does NOT ask for confirmation")
	flag.BoolVar(&isDirectoryMode, "dir", false, "recursively uploads an entire directory. if no path is provided, current directory will be assumed")

	// TODO: add support for bucketName
	//var bucketName string
	//flag.StringVar(&bucketName, "bucket", "", "name of destination bucket")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A simple utility for uploading stuff to BackBlaze's B2\n\n")
		fmt.Fprintf(os.Stderr, "usage: up [-pipe dst_name] [-dir path] [file1 .. fileN]\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	if err := run(isPipeMode, isDirectoryMode, args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
