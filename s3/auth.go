package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
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

// This method is called via Bubble cmd so there is no need
// to run it in a goroutine (done via Bubbletea already)
func GetSessionFromInput(key, secret string) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(key, secret, "")
	_, err := creds.Get()

	if err != nil {
		return nil, err
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	}))

	client := sts.New(sess)
	_, err = client.GetCallerIdentity(nil)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
