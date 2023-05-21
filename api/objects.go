package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetObjects(session *session.Session, bucket string) (*s3.ListObjectsV2Output, error) {
	client := s3.New(session)
	input := s3.ListObjectsV2Input{Bucket: &bucket}
	o, err := client.ListObjectsV2(&input)

	if err != nil {
		return nil, err
	}

	return o, nil
}
