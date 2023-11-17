package s3router

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextFactory_Build(t *testing.T) {
	publicHost := "public.example.com"
	privateHost := "private.example.com"

	type testCase struct {
		name string

		url string

		expected Context
	}

	testCases := []testCase{
		{
			"VPath",
			"http://" + publicHost,
			Context{
				Hostname: publicHost,
				BaseHost: publicHost,
				Style:    RequestStyleVPath,
			},
		},
		{
			"VPath+bucket",
			"http://" + publicHost + "/bucket",
			Context{
				Hostname: publicHost,
				BaseHost: publicHost,
				Style:    RequestStyleVPath,
				Bucket:   "bucket",
			},
		},
		{
			"VPath+object",
			"http://" + publicHost + "/bucket/key",
			Context{
				Hostname: publicHost,
				BaseHost: publicHost,
				Style:    RequestStyleVPath,
				Bucket:   "bucket",
				Key:      "key",
			},
		},
		{
			"VPath+object+other host",
			"http://" + privateHost + "/bucket/key",
			Context{
				Hostname: privateHost,
				BaseHost: privateHost,
				Style:    RequestStyleVPath,
				Bucket:   "bucket",
				Key:      "key",
			},
		},

		{
			"VHost+bucket",
			"http://bucket." + publicHost,
			Context{
				Hostname: "bucket." + publicHost,
				BaseHost: publicHost,
				Style:    RequestStyleVHost,
				Bucket:   "bucket",
			},
		},
		{
			"VHost+object",
			"http://bucket." + publicHost + "/key",
			Context{
				Hostname: "bucket." + publicHost,
				BaseHost: publicHost,
				Style:    RequestStyleVHost,
				Bucket:   "bucket",
				Key:      "key",
			},
		},
		{
			"VHost+object+other host",
			"http://bucket." + privateHost + "/key",
			Context{
				Hostname: "bucket." + privateHost,
				BaseHost: privateHost,
				Style:    RequestStyleVHost,
				Bucket:   "bucket",
				Key:      "key",
			},
		},

		{
			"CName",
			"http://my-cname.com",
			Context{
				Hostname: "my-cname.com",
				BaseHost: "my-cname.com",
				Style:    RequestStyleCName,
				Bucket:   "my-cname.com",
			},
		},
		{
			"CName+object",
			"http://my-cname.com/object",
			Context{
				Hostname: "my-cname.com",
				BaseHost: "my-cname.com",
				Style:    RequestStyleCName,
				Bucket:   "my-cname.com",
				Key:      "object",
			},
		},
	}

	factory := ContextFactory{
		Hosts: []string{
			publicHost,
			privateHost,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create input
			r, err := http.NewRequest(http.MethodGet, tc.url, nil)
			require.NoError(t, err)

			// Run and test
			actual := factory.Build(r)
			require.NotNil(t, actual)
			require.Equal(t, &tc.expected, actual)
		})
	}
}
