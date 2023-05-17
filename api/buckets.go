package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetBuckets(session *session.Session) ([]*s3.Bucket, error) {
	client := s3.New(session)
	b, err := client.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		return nil, err
	}

	return b.Buckets, nil
}
