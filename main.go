package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/kurin/blazer/b2"
)

// B2Client - a temporary global pointer to an authenticated b2.Client
var B2Client *b2.Client

func run(isPipeMode bool, isDirectoryMode bool, bucketName string, args []string) error {

	if isPipeMode {
		if isDirectoryMode {
			return errors.New("directory mode and pipe mode are mutually exclusive")
		}
		return handlePipeMode(bucketName, args)
	}

	if len(args) > 0 {

		fmt.Println("About to upload the following file[/s]:", args)
		confirmed, err := confirmAction()
		if nil != err {
			return err
		}

		if confirmed {
			return handleFiles(bucketName, args, isDirectoryMode)
		}
	}

	return nil
}

func main() {

	var isPipeMode bool
	var isDirectoryMode bool
	var bucketName string

	flag.BoolVar(&isPipeMode, "pipe", false, "reads and uploads data from STDIN until EOF is reached. Does NOT ask for confirmation")
	flag.BoolVar(&isDirectoryMode, "dir", false, "recursively uploads an entire directory.")
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

	fmt.Println("\ndone")
}
