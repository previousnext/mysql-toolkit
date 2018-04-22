package awsutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

// NewSession creates a new AWS session.
func NewSession(region string) (*session.Session, error) {
	// Attempt to use default AWS credential chain.
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get credentials")
	}

	return sess, nil
}
