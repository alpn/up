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

func run(isPipeMode bool, allowDirectories bool, bucketName string, args []string) error {

	if isPipeMode {
		if allowDirectories {
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
			return handleFiles(bucketName, args, allowDirectories)
		}
	}

	return nil
}

func main() {

	var isPipeMode bool
	var allowDirectories bool
	var bucketName string

	flag.BoolVar(&isPipeMode, "pipe", false,
		"Read and upload data from STDIN until EOF is reached. Does NOT ask for confirmation")
	flag.BoolVar(&allowDirectories, "dir", false,
		"Upload directories (recursively). When this option is not specified, directories are ignored")
	flag.StringVar(&bucketName, "bucket", "", "name of destination bucket")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUp - A simple utility for uploading stuff to BackBlaze's B2\n\n")
		fmt.Fprintf(os.Stderr, "Usage: up [-bucket BUCKET_NAME] [-pipe TARGET_NAME] [-dir] [file1 .. fileN]\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	fmt.Printf("\nUp⤴️\n\n")
	if err := run(isPipeMode, allowDirectories, bucketName, args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\ndone")
}
