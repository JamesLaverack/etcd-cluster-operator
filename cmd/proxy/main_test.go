package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

var backupUrlTests = []struct {
	blobUrl string
	bucket  string
	path    string
}{
	{"gs://foo/bar", "gs://foo", "/bar"},
	{"s3://foo/bar", "s3://foo", "/bar"},
	{"s3://foo/bar/baz.mkv", "s3://foo", "/bar/baz.mkv"},
	{"s3://mybucket/foo/bar.mkv?endpoint=my.minio.local:8080&disableSSL=true&s3ForcePathStyle=true",
		"s3://mybucket?endpoint=my.minio.local:8080&disableSSL=true&s3ForcePathStyle=true",
		"/foo/bar.mkv"},
}

func TestParseBackupURL(t *testing.T) {
	for _, tt := range backupUrlTests {
		t.Run(tt.blobUrl, func(t *testing.T) {
			actualBucket, actualPath, err := parseBackupURL(tt.blobUrl)
			assert.NoError(t, err)
			assert.Equal(t, tt.bucket, actualBucket)
			assert.Equal(t, tt.path, actualPath)
		})
	}
}

func TestFoo(t *testing.T) {
	u, err := url.Parse("proxy.eco-system.svc:8080")
	fmt.Print(u)
	u.Port()
	require.NoError(t, err)
	require.Equal(t, "proxy.eco-system.svc", u.Host)
}