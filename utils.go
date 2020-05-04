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

func getDstPathString(bucket *b2.Bucket, fileName string) string {
	return "b2://" + bucket.Name() + "/" + fileName
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

func chooseBucket(buckets []*b2.Bucket) (*b2.Bucket, error) {

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

func getBucket(name string) (context.Context, *b2.Bucket, error) {
	ctx, buckets, err := getAllBuckets()
	if nil != err {
		return nil, nil, err
	}

	for _, b := range buckets {
		if b.Name() == name {
			return ctx, b, nil
		}
	}

	bucket, err := chooseBucket(buckets)
	return ctx, bucket, err
}
