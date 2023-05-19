package types

import "github.com/aws/aws-sdk-go/aws/session"

const (
	Creds   CurrentPage = "creds"
	Buckets             = "buckets"
)

type CurrentPage string

// This is the main model used for the overall UI and for
// pages to pass information back and forth to each other.
type UiModel struct {
	CurrentPage CurrentPage
	Session     *session.Session
}
