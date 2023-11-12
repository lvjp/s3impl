package s3errors

import (
	"fmt"
)

type S3Error struct {
	HTTPStatusCode int

	Code      string
	Message   string
	RequestID string
	Resource  string
}

func (e *S3Error) Error() string {
	return fmt.Sprintf("s3error: %s %s", e.Code, e.Message)
}
