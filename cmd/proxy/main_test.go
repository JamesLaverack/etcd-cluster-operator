package main

import (
	"github.com/stretchr/testify/assert"
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
}
func TestParseBackupUrl(t *testing.T) {
	for _, tt := range backupUrlTests {
		t.Run(tt.blobUrl, func(t *testing.T) {
			actualBucket, actualPath, err := parseBackupUrl(tt.blobUrl)
			assert.NoError(t, err)
			assert.Equal(t, tt.bucket, actualBucket)
			assert.Equal(t, tt.path, actualPath)
		})
	}
}