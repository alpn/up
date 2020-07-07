package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/kurin/blazer/b2"
)

func getDstPathString(bucketName string, fileName string) string {
	return "b2://" + bucketName + "/" + fileName
}

func confirmAction() (bool, error) {

	fmt.Println("Is that OK? [yes/no]")

	r := bufio.NewReader(os.Stdin)
	response, err := r.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = response[:len(response)-1]

	if "yes" == response || "y" == response {
		return true, nil
	}
	return false, nil
}

func getAllBuckets() (context.Context, []*b2.Bucket, error) {

	id := os.Getenv("B2_ACCOUNT_ID")
	key := os.Getenv("B2_ACCOUNT_KEY")

	if 0 == len(id) {
		return nil, nil, errors.New("Account id is missing")
	}

	if 0 == len(key) {
		return nil, nil, errors.New("Account key is missing")
	}

	ctx := context.Background()

	// b2_authorize_account
	b2client, err := b2.NewClient(ctx, id, key)
	if err != nil {
		return nil, nil, err
	}

	buckets, err := b2client.ListBuckets(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(buckets) > 0 {
		B2Client = b2client
		return ctx, buckets, nil
	}

	return nil, nil, errors.New("No buckets were found")
}

func promptUserToChooseBucket(buckets []*b2.Bucket) (*b2.Bucket, error) {

	for {
		fmt.Println("The following buckets are available:")

		for i, b := range buckets {
			fmt.Printf("[%d] - %s\n", i, b.Name())
		}

		fmt.Println("Please choose one:")
		r := bufio.NewReader(os.Stdin)
		response, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		response = response[:len(response)-1]
		i, err := strconv.Atoi(response)
		if nil != err {
			continue
		}
		if i < len(buckets) {
			return buckets[i], nil
		}
	}
}

func pickBucket(buckets []*b2.Bucket, name string) *b2.Bucket {

	if len(name) == 0 && len(buckets) == 1 {
		return buckets[0]
	}

	for _, b := range buckets {
		if b.Name() == name {
			return b
		}
	}

	return nil
}

func pickBucketPrompt(buckets []*b2.Bucket, name string) (*b2.Bucket, error) {

	bucket := pickBucket(buckets, name)
	if nil != bucket {
		return bucket, nil
	}

	bucket, err := promptUserToChooseBucket(buckets)
	if nil != err {
		return nil, err
	}
	return bucket, nil
}

// ByteCountSI - credit: https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
