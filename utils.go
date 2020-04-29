package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

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

func getBucket() (context.Context, *b2.Bucket, error) {

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
	b2, err := b2.NewClient(ctx, id, key)
	if err != nil {
		return nil, nil, err
	}

	buckets, err := b2.ListBuckets(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(buckets) > 0 {
		bucket := buckets[0]
		return ctx, bucket, nil
	}

	return nil, nil, errors.New("No buckets were found")
}
