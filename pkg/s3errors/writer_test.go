package s3errors

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIWriter_Write(t *testing.T) {
	input := &S3Error{
		HTTPStatusCode: 404,
		Code:           "NoSuchBucket",
		Message:        "The specified bucket does not exist",
		RequestID:      "5CZT884BVHY7AYJN",
	}

	expectedBody := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		"<Error>" +
		"<Code>NoSuchBucket</Code>" +
		"<Message>The specified bucket does not exist</Message>" +
		"<RequestId>5CZT884BVHY7AYJN</RequestId>" +
		"<HostId>5CZT884BVHY7AYJN</HostId>" +
		"</Error>"

	recorder := httptest.NewRecorder()

	writer := APIWriter{}
	err := writer.Write(input, recorder)
	require.NoError(t, err)

	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, input.HTTPStatusCode, resp.StatusCode)
	require.Equal(t, input.RequestID, resp.Header.Get("X-Amz-Request-Id"))
	require.Equal(t, input.RequestID, resp.Header.Get("X-Amz-Id-2"))
	require.NotNil(t, recorder.Body)
	require.Equal(t, expectedBody, recorder.Body.String())
}
