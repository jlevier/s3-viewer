package types

import "github.com/aws/aws-sdk-go/aws/session"

const (
	Creds   CurrentPage = "creds"
	Buckets             = "buckets"
	Files               = "files"
)

type CurrentPage string

// Used to change the current page.  Use this msg instead of just setting directly the current page on the
// UiModel because it is necessary for the main Update() method in control to call the Init() function of the new page.
type ChangeCurrentPageMsg struct {
	CurrentPage   CurrentPage
	CurrentBucket string
}

// This is the main model used for the overall UI and for
// pages to pass information back and forth to each other.
type UiModel struct {
	CurrentPage   CurrentPage
	Session       *session.Session
	CurrentBucket string
}
