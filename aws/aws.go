package aws

import (
	aw "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWS is our helper type for aws-related operations.
type AWS struct {
	Ses *session.Session
}

// New creates a new AWS instance.
func New() (*AWS, error) {
	ses, err := NewSession(0, "")
	if err != nil {
		return nil, err
	}

	return &AWS{
		Ses: ses,
	}, nil
}

// NewSession creates a new aws session.
func NewSession(maxRetries int, region string) (*session.Session, error) {
	awsCfg := aw.NewConfig().WithMaxRetries(maxRetries) //.WithLogLevel(aws.LogDebug))
	if region != "" {
		awsCfg = awsCfg.WithRegion(region)
	}
	return session.NewSessionWithOptions(session.Options{Config: *awsCfg, SharedConfigState: session.SharedConfigEnable})
}
