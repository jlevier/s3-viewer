package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// "Objects" is the name that s3 gives to anything in a bucket.  This includes "directories" and files.
// Directories is quoted here because they do not distinguish between directories and files.  So you might encounter
// an Object whose value is /foo/bar/ as well as another that is file.json.
func GetObjects(session *session.Session, bucket, directory string) (*s3.ListObjectsV2Output, error) {
	client := s3.New(session)
	input := s3.ListObjectsV2Input{Bucket: &bucket, Delimiter: &directory}
	o, err := client.ListObjectsV2(&input)

	if err != nil {
		return nil, err
	}

	return o, nil
}
