package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SessionResponse struct {
	Session *session.Session
	Err     error
}

func GetSession(ch chan<- *SessionResponse) {
	creds := credentials.NewEnvCredentials()
	_, err := creds.Get()

	if err != nil {
		ch <- &SessionResponse{nil, err}
		return
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	}))

	ch <- &SessionResponse{sess, nil}
}
